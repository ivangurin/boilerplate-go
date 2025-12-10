package handlers

import (
	"boilerplate/internal/api/grpc/handlers/auth"
	"boilerplate/internal/api/grpc/handlers/users"
	"boilerplate/internal/model"
	"boilerplate/internal/service_provider"
)

func NewHandlers(
	sp *service_provider.Provider,
) []model.GRPCHandler {
	return []model.GRPCHandler{
		auth.NewHandler(
			sp.GetAuthService(),
		),
		users.NewHandler(
			sp.GetUsersService(),
		),
	}
}
