package service_provider

import (
	"boilerplate/internal/services/auth"
	"boilerplate/internal/services/users"
)

type services struct {
	auth  auth.Service
	users users.Service
}

func (p *Provider) GetAuthService() auth.Service {
	if p.services.auth == nil {
		p.services.auth = auth.NewService(
			&p.config.API,
			p.GetUsersService(),
		)
	}
	return p.services.auth
}

func (p *Provider) GetUsersService() users.Service {
	if p.services.users == nil {
		p.services.users = users.NewService(p.repo)
	}
	return p.services.users
}
