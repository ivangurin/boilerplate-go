package suite_provider

import (
	"context"

	"boilerplate/internal/model"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/repository"
)

type Provider struct {
	ctx      context.Context
	config   *model.Config
	logger   logger_pkg.Logger
	repo     repository.Repo
	clients  clients // nolint: unused
	services services
	handlers handlers // nolint: unused
	cleanups []func() error
}

func NewProvider() (*Provider, func()) {
	logger, err := logger_pkg.New()
	if err != nil {
		panic(err)
	}

	sp := &Provider{
		ctx:    context.Background(),
		config: InitConfig(),
		logger: logger,
	}

	sp.cleanups = append(sp.cleanups,
		func() error {
			sp.ctx.Done()
			return nil
		},
	)

	return sp, sp.Cleaner
}

func (sp *Provider) Cleaner() {
	for _, cleanup := range sp.cleanups {
		err := cleanup()
		if err != nil {
			panic(err)
		}
	}
}

func (sp *Provider) Context() context.Context {
	return sp.ctx
}

func (sp *Provider) ContextWithValue(key, val any) context.Context {
	return context.WithValue(sp.Context(), key, val)
}

func (sp *Provider) GetLogger() logger_pkg.Logger {
	return sp.logger
}
