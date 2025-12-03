package auth

import (
	"context"

	"github.com/gin-gonic/gin"

	"boilerplate/internal/model"
	"boilerplate/internal/services/users"
)

type Service interface {
	GetConfig() *model.ConfigAPI
	Login(ctx context.Context, req *AuthLoginRequest) (*AuthLoginResponse, error)
	Refresh(ctx context.Context, req *AuthRefreshRequest) (*AuthRefreshResponse, error)
	Me(ctx context.Context) (*users.User, error)
	Validate() gin.HandlerFunc
}

type service struct {
	config       *model.ConfigAPI
	usersService users.Service
}

func NewService(
	config *model.ConfigAPI,
	usersService users.Service,
) Service {
	return &service{
		config:       config,
		usersService: usersService,
	}
}
