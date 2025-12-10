package middleware

import (
	"context"

	"google.golang.org/grpc"

	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/services/auth"
)

type Middleware interface {
	Panic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	Tracer(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	Validate(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
	Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error)
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
