package suite_provider

import (
	"boilerplate/internal/services/auth"
	"boilerplate/internal/services/users"
)

type services struct {
	auth  auth.Service
	users users.Service
}

func (sp *Provider) GetAuthService() auth.Service {
	if sp.services.auth == nil {
		sp.services.auth = auth.NewService(
			&sp.GetConfig().API,
			sp.GetUserService(),
		)
	}
	return sp.services.auth
}

func (sp *Provider) GetUserService() users.Service {
	if sp.services.users == nil {
		sp.services.users = users.NewService(sp.GetRepo())
	}
	return sp.services.users
}
