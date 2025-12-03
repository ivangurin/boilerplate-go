package suite_provider

import "boilerplate/internal/services/users"

type services struct {
	users users.Service
}

func (sp *Provider) GetUserService() users.Service {
	if sp.services.users == nil {
		sp.services.users = users.NewService(sp.GetRepo())
	}
	return sp.services.users
}
