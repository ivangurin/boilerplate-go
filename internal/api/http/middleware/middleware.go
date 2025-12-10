package middleware

import (
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Logger() gin.HandlerFunc
	Auth() gin.HandlerFunc
}

type middleware struct {
	logger      logger_pkg.Logger
	authService auth.Service
}

func NewMiddleware(
	logger logger_pkg.Logger,
	authService auth.Service,
) Middleware {
	return &middleware{
		authService: authService,
		logger:      logger,
	}
}
