package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"boilerplate/internal/api/grpc/handlers/users"
	"boilerplate/internal/pkg/grpc"
	"boilerplate/pkg/pb"
)

func (h *handler) Me(ctx context.Context, req *emptypb.Empty) (*pb.AuthMeResponse, error) {
	resp, err := h.authService.Me(ctx)
	if err != nil {
		return nil, grpc.Error(err)
	}

	return &pb.AuthMeResponse{
		User: users.ToUser(resp),
	}, nil
}
