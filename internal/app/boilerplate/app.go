package boilerplate

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"syscall"

	grpc_handlers "boilerplate/internal/api/grpc/handlers"
	grpc_middleware "boilerplate/internal/api/grpc/middleware"
	http_handlers "boilerplate/internal/api/http/handlers"
	"boilerplate/internal/model"
	"boilerplate/internal/pkg/clients/db"
	closer_pkg "boilerplate/internal/pkg/closer"
	"boilerplate/internal/pkg/grpc_server"
	"boilerplate/internal/pkg/http_server"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/pkg/swagger"
	"boilerplate/internal/repository"
	"boilerplate/internal/service_provider"
	"boilerplate/migrations"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	// Service Provider
	sp := service_provider.NewProvider(a.config, logger, repo)

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

	// HTTP Gateway Server
	httpRouter, err := a.setupHTTPGateway(ctx, grpcHandlers)
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

func (a *App) setupHTTPGateway(ctx context.Context, grpcHandlers []model.GRPCHandler) (*http.ServeMux, error) {
	// Create gRPC client connection
	grpcConn, err := grpc.NewClient(
		net.JoinHostPort(a.config.API.Host, a.config.API.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create grpc client: %w", err)
	}

	// Setup main router
	router := http.NewServeMux()

	// Register swagger UI
	swagger.Register(router)

	// Setup gRPC gateway
	gwRouter := runtime.NewServeMux()
	for _, handler := range grpcHandlers {
		if err := handler.RegisterHTTPHandler(ctx, gwRouter, grpcConn); err != nil {
			return nil, fmt.Errorf("register grpc handler: %w", err)
		}
	}

	// Mount gateway to main router
	router.Handle("/", gwRouter)

	return router, nil
}
