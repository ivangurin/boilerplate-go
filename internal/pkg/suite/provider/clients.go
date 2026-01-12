package suite_provider

import (
	"boilerplate/internal/model"
	model_mocks "boilerplate/internal/model/mocks"
	"boilerplate/internal/pkg/clients/chrome"
	"boilerplate/internal/pkg/clients/mail"
	"boilerplate/internal/pkg/clients/mail/mocks"
	"boilerplate/internal/pkg/clients/s3"
)

type clients struct {
	s3Client     s3.Client
	chromeClient chrome.Client
	brokerClient model.BrokerClient
	mailClient   mail.Client
}

func (p *Provider) GetS3Client() s3.Client {
	if p.clients.s3Client == nil {
		var err error
		p.clients.s3Client, err = s3.NewClient(p.Context(), p.config.S3.Host, p.config.S3.Port, p.config.S3.AccessKey, p.config.S3.SecretKey, p.config.S3.Bucket, s3.WithLogger(p.GetLogger()))
		if err != nil {
			panic(err)
		}

		_, err = p.clients.s3Client.CreateBucket(p.Context(), p.config.S3.Bucket)
		if err != nil {
			panic(err)
		}
	}
	return p.clients.s3Client
}

func (p *Provider) GetChromeClient() chrome.Client {
	if p.clients.chromeClient == nil {
		p.clients.chromeClient = chrome.NewClient(p.config.Chrome.Host, p.config.Chrome.Port, p.config.Chrome.Timeout)
	}
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
		p.clients.mailClient = &mocks.Client{}
	}
	return p.clients.mailClient
}
