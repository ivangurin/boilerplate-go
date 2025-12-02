package service_provider

import "boilerplate/internal/services/users"

type services struct {
	users users.IService
}

func (p *Provider) GetUsersService() users.IService {
	if p.services.users == nil {
		p.services.users = users.NewService(p.repo)
	}
	return p.services.users
}
