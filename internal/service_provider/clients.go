package service_provider

import (
	"boilerplate/internal/model"
	model_mocks "boilerplate/internal/model/mocks"
	"boilerplate/internal/pkg/clients/chrome"
	"boilerplate/internal/pkg/clients/mail"
	"boilerplate/internal/pkg/clients/s3"
)

type clients struct {
	s3Client     s3.Client
	chromeClient chrome.Client
	brokerClient model.BrokerClient
	mailClient   mail.Client
}

func (p *Provider) GetS3Client() s3.Client {
	return p.clients.s3Client
}

func (p *Provider) GetChromeClient() chrome.Client {
	return p.clients.chromeClient
}

func (p *Provider) GetBrokerClient() model.BrokerClient {
	if p.clients.brokerClient == nil {
		p.clients.brokerClient = &model_mocks.BrokerClient{}
	}
	return p.clients.brokerClient
}

func (p *Provider) GetMailClient() mail.Client {
	if p.clients.mailClient == nil {
		p.clients.mailClient = mail.NewClient(
			p.GetLogger(),
			p.config.Mail)
	}
	return p.clients.mailClient
}
