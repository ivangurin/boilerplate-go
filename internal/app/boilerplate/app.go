package boilerplate

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"boilerplate/internal/api/handlers"
	"boilerplate/internal/model"
	"boilerplate/internal/pkg/clients/db"
	closer_pkg "boilerplate/internal/pkg/closer"
	"boilerplate/internal/pkg/http_server"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/repository"
	"boilerplate/internal/service_provider"
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
	dbClient, err := db.New(ctx, a.config.DB.GetDSN())
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create db client: %w", err)
	}

	err = migrations.Migrate(dbClient)
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
	httpServer := http_server.NewServer(a.config.API.Port, handlers.NewHandler(logger, sp))

	go func() {
		logger.Infof(ctx, "http server started on port %s", a.config.API.Port)
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

	logger.Info(ctx, "App is running...")

	return nil
}
