package users

import (
	"context"

	"boilerplate/internal/pkg/convert"
	"boilerplate/internal/pkg/grpc"
	"boilerplate/pkg/pb"
)

func (h *handler) Get(ctx context.Context, req *pb.UserGetRequest) (*pb.User, error) {
	resp, err := h.usersService.Get(ctx, convert.ToInt(req.Id))
	if err != nil {
		return nil, grpc.Error(err)
	}

	return ToUser(resp), nil
}
