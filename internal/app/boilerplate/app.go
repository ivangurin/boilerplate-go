package boilerplate

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"google.golang.org/grpc"

	grpc_handlers "boilerplate/internal/api/grpc/handlers"
	grpc_middleware "boilerplate/internal/api/grpc/middleware"
	http_handlers "boilerplate/internal/api/http/handlers"
	consumers_pkg "boilerplate/internal/consumers"
	"boilerplate/internal/model"
	"boilerplate/internal/pkg/clients/chrome"
	"boilerplate/internal/pkg/clients/db"
	nats_client "boilerplate/internal/pkg/clients/nats"
	"boilerplate/internal/pkg/clients/s3"
	closer_pkg "boilerplate/internal/pkg/closer"
	"boilerplate/internal/pkg/gateway"
	logger_pkg "boilerplate/internal/pkg/logger"
	grpc_server "boilerplate/internal/pkg/servers/grpc"
	http_server "boilerplate/internal/pkg/servers/http"
	nats_server "boilerplate/internal/pkg/servers/nats"
	"boilerplate/internal/repository"
	"boilerplate/internal/service_provider"
	"boilerplate/internal/topics"
	"boilerplate/migrations"
)

type App struct {
	config *model.Config
}

func New(
	config *model.Config,
) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	debug := a.config.LogLevel == logger_pkg.LevelDebug

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Logger
	logger, err := logger_pkg.New(logger_pkg.WithLevel(a.config.LogLevel), logger_pkg.WithOutputStdout())
	if err != nil {
		panic(err)
	}

	// Closer
	closer := closer_pkg.New(ctx, logger, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer closer.Wait()

	closer.Add(logger.Close)

	// DB Client
	dbClient, err := db.New(ctx, logger, a.config.DB.GetDSN())
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create db client: %w", err)
	}

	err = migrations.Migrate(ctx, dbClient)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("migrate db: %w", err)
	}

	closer.Add(dbClient.Close)

	// Repository
	repo := repository.NewRepo(dbClient)

	// S3 Client
	s3Client, err := s3.NewClient(ctx, a.config.S3.Host, a.config.S3.Port, a.config.S3.AccessKey, a.config.S3.SecretKey, a.config.S3.Bucket, s3.WithLogger(logger))
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create s3 client: %w", err)
	}
	logger.Info(ctx, "s3 client created")

	bucketCreated, err := s3Client.CreateBucket(ctx, a.config.S3.Bucket)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create bucket %s: %w", a.config.S3.Bucket, err)
	}
	if bucketCreated {
		logger.Infof(ctx, "s3 bucket %s created", a.config.S3.Bucket)
	}

	// Chrome Client
	chromeClient := chrome.NewClient(a.config.Chrome.Host, a.config.Chrome.Port, a.config.Chrome.Timeout)

	// Broker Server
	brokerServer, err := nats_server.NewServer(&a.config.Nats, nats_server.WithName("boilerplate-nats-server"), nats_server.WithJetStream(a.config.Nats.Domain), nats_server.WithDebug(debug))
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create nats server: %w", err)
	}

	logger.Infof(ctx, "nats server starting...")
	if err = brokerServer.Start(); err != nil {
		closer.CloseAll()
		return fmt.Errorf("start nats server: %w", err)
	}
	logger.Infof(ctx, "nats server started")

	closer.Add(func() error {
		err := brokerServer.Stop()
		if err != nil {
			return fmt.Errorf("stop nats server: %w", err)
		}
		logger.Info(ctx, "nats server stopped")
		return nil
	})

	// Broker client
	brokerConn, err := brokerServer.GetConn()
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("get nats server conn: %w", err)
	}

	logger.Info(ctx, "nats server conn obtained")
	brokerClient, err := nats_client.NewClient(logger, nats_client.WithName("boilerplate-nats-client"), nats_client.WithConn(brokerConn))
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create nats client: %w", err)
	}

	logger.Info(ctx, "nats client created")

	closer.Add(func() error {
		err := brokerClient.Close()
		if err != nil {
			return fmt.Errorf("close nats client: %w", err)
		}
		logger.Info(ctx, "nats client closed")
		return nil
	})

	// Service Provider
	sp := service_provider.NewProvider(a.config, logger, repo, s3Client, chromeClient, brokerClient)

	// Create or update topics
	err = topics.CreateOrUpdateTopics(ctx, brokerClient)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create or update topics: %w", err)
	}

	// Start Consumers
	consumers := consumers_pkg.NewConsumers(logger, brokerClient, sp)

	err = consumers.Start(ctx)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("start consumers: %w", err)
	}
	logger.Info(ctx, "consumers started")

	// HTTP Server
	if 1 == 2 {
		httpServer := http_server.NewServer(a.config.API.Host, a.config.API.HTTPPort, http_handlers.NewHandler(logger, sp))

		go func() {
			logger.Infof(ctx, "http server started on port %s", a.config.API.HTTPPort)
			err := httpServer.Start()
			if err != nil {
				closer.CloseAll()
				logger.Errorf(ctx, "start http server: %s", err.Error())
			}
		}()

		closer.Add(func() error {
			err := httpServer.Stop(ctx)
			if err != nil {
				return fmt.Errorf("stop http server: %w", err)
			}
			logger.Info(ctx, "http server stopped")
			return nil
		})
	}

	// GRPC Server
	grpcMiddleware := grpc_middleware.NewMiddleware(logger, sp.GetAuthService())
	grpcHandlers := grpc_handlers.NewHandlers(sp)

	grpcServer := grpc_server.NewServer(
		a.config.API.Host,
		a.config.API.GRPCPort,
		[]grpc.UnaryServerInterceptor{
			grpcMiddleware.Panic,
			grpcMiddleware.Tracer,
			grpcMiddleware.Logger,
			grpcMiddleware.Validate,
			grpcMiddleware.Auth,
		},
		grpcHandlers,
	)

	go func() {
		logger.Infof(ctx, "grpc server started on port %s", a.config.API.GRPCPort)
		if err := grpcServer.Start(); err != nil {
			closer.CloseAll()
			logger.Errorf(ctx, "start grpc server: %s", err.Error())
		}
	}()

	closer.Add(func() error {
		grpcServer.Stop()
		logger.Info(ctx, "grpc server stopped")
		return nil
	})

	httpRouter, err := gateway.Setup(ctx, a.config.API, grpcHandlers)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("setup http gateway: %w", err)
	}

	httpServer := http_server.NewServer(a.config.API.Host, a.config.API.HTTPPort, httpRouter)

	go func() {
		logger.Infof(ctx, "http server started on port %s", a.config.API.HTTPPort)
		if err := httpServer.Start(); err != nil {
			closer.CloseAll()
			logger.Errorf(ctx, "start http server: %s", err.Error())
		}
	}()

	closer.Add(func() error {
		if err := httpServer.Stop(ctx); err != nil {
			return fmt.Errorf("stop http server: %w", err)
		}
		logger.Info(ctx, "http server stopped")
		return nil
	})

	logger.Info(ctx, "App is running...")

	return nil
}
