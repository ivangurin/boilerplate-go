package boilerplate

import (
	"context"
	"fmt"
	"os"
	"syscall"

	"boilerplate/internal/model"
	"boilerplate/internal/pkg/clients/db"
	closer_pkg "boilerplate/internal/pkg/closer"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/migrations"
)

type App struct {
	config model.Config
}

func New(
	config model.Config,
) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := logger_pkg.New(logger_pkg.WithDebugLevel(), logger_pkg.WithOutputStdout())
	if err != nil {
		panic(err)
	}

	closer := closer_pkg.New(ctx, logger, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer closer.Wait()

	closer.Add(logger.Close)

	dbClient, err := db.New(ctx, a.config.DbDsn)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("create db client: %s", err.Error())
	}

	err = migrations.Migrate(dbClient)
	if err != nil {
		closer.CloseAll()
		return fmt.Errorf("migrate db: %s", err.Error())
	}

	closer.Add(dbClient.Close)

	logger.Info(ctx, "App is running...")

	return nil
}
