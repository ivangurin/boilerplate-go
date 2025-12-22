package users

import (
	"context"

	"boilerplate/internal/pkg/convert"
	"boilerplate/internal/pkg/grpc"
	"boilerplate/internal/services/users"
	"boilerplate/pkg/pb"
)

func (h *handler) Update(ctx context.Context, req *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	resp, err := h.usersService.Update(ctx, &users.UserUpdateRequest{
		ID:       convert.ToInt(req.GetUserId()),
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		return nil, grpc.Error(err)
	}

	return &pb.UserUpdateResponse{
		User: ToUser(resp),
	}, nil
}
