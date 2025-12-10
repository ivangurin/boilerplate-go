package auth

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"boilerplate/internal/model"
	"boilerplate/internal/services/auth"
	"boilerplate/pkg/pb"
)

type handler struct {
	pb.UnimplementedAuthAPIServer
	authService auth.Service
}

func NewHandler(
	authService auth.Service,
) model.GRPCHandler {
	return &handler{
		authService: authService,
	}
}

func (h *handler) RegisterGRPCServer(server *grpc.Server) {
	pb.RegisterAuthAPIServer(server, h)
}

func (h *handler) RegisterHTTPHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return pb.RegisterAuthAPIHandler(ctx, mux, conn)
}
