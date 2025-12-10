package users

import (
	"context"

	"boilerplate/internal/pkg/grpc"
	"boilerplate/internal/services/users"
	"boilerplate/pkg/pb"
)

func (h *handler) Create(ctx context.Context, req *pb.UserCreateRequest) (*pb.User, error) {
	resp, err := h.usersService.Create(ctx, &users.UserCreateRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, grpc.Error(err)
	}

	return ToUser(resp), nil
}
