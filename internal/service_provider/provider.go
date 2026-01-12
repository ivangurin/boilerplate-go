package service_provider

import (
	"boilerplate/internal/model"
	"boilerplate/internal/pkg/clients/chrome"
	"boilerplate/internal/pkg/clients/s3"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/repository"
)

type Provider struct {
	config   *model.Config
	logger   logger_pkg.Logger
	repo     repository.Repo
	clients  clients
	services services
}

func NewProvider(
	config *model.Config,
	logger logger_pkg.Logger,
	repo repository.Repo,
	s3Client s3.Client,
	chromeClient chrome.Client,
	brokerClient model.BrokerClient,
) *Provider {
	return &Provider{
		config: config,
		logger: logger,
		repo:   repo,
		clients: clients{
			s3Client:     s3Client,
			chromeClient: chromeClient,
			brokerClient: brokerClient,
		},
	}
}

func (p *Provider) GetLogger() logger_pkg.Logger {
	return p.logger
}

func GetRepo(p *Provider) repository.Repo {
	return p.repo
}
