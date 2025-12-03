package service_provider

import (
	"boilerplate/internal/model"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/repository"
)

type Provider struct {
	config   *model.Config
	logger   logger_pkg.Logger
	repo     repository.Repo
	clients  clients // nolint: unused
	services services
}

func NewProvider(
	config *model.Config,
	logger logger_pkg.Logger,
	repo repository.Repo,
) *Provider {
	return &Provider{
		config: config,
		logger: logger,
		repo:   repo,
	}
}

func GetRepo(p *Provider) repository.Repo {
	return p.repo
}
