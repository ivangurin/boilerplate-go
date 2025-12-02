package service_provider

import (
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/repository"
)

type Provider struct {
	logger   logger_pkg.Logger
	repo     repository.Repo
	clients  clients // nolint: unused
	services services
}

func NewProvider(
	logger logger_pkg.Logger,
	repo repository.Repo,
) *Provider {
	return &Provider{
		logger: logger,
		repo:   repo,
	}
}

func GetRepo(p *Provider) repository.Repo {
	return p.repo
}
