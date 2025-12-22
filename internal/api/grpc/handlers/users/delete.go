package users

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"boilerplate/internal/pkg/convert"
	"boilerplate/internal/pkg/grpc"
	"boilerplate/pkg/pb"
)

func (h *handler) Delete(ctx context.Context, req *pb.UserDeleteRequest) (*emptypb.Empty, error) {
	err := h.usersService.Delete(ctx, convert.ToInt(req.GetUserId()))
	if err != nil {
		return nil, grpc.Error(err)
	}

	return &emptypb.Empty{}, nil
}
