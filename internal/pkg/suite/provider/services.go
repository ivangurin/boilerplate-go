package suite_provider

import "boilerplate/internal/services/users"

type services struct {
	users users.IService
}

func (sp *Provider) GetUserService() users.IService {
	if sp.services.users == nil {
		sp.services.users = users.NewService(sp.GetRepo())
	}
	return sp.services.users
}
