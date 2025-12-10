package users

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"boilerplate/internal/model"
	"boilerplate/internal/services/users"
	"boilerplate/pkg/pb"
)

type handler struct {
	pb.UnimplementedUsersAPIServer
	usersService users.Service
}

func NewHandler(
	usersService users.Service,
) model.GRPCHandler {
	return &handler{
		usersService: usersService,
	}
}

func (h *handler) RegisterGRPCServer(server *grpc.Server) {
	pb.RegisterUsersAPIServer(server, h)
}

func (h *handler) RegisterHTTPHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return pb.RegisterUsersAPIHandler(ctx, mux, conn)
}
