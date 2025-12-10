package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"boilerplate/internal/pkg/convert"
	"boilerplate/internal/pkg/grpc"
	"boilerplate/pkg/pb"
)

func (h *handler) Me(ctx context.Context, _ *emptypb.Empty) (*pb.User, error) {
	user, err := h.authService.Me(ctx)
	if err != nil {
		return nil, grpc.Error(err)
	}

	return &pb.User{
		Id:        convert.ToInt64(user.ID),
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		IsAdmin:   user.IsAdmin,
		Deleted:   user.Deleted,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		DeletedAt: func() *timestamppb.Timestamp {
			if user.DeletedAt == nil {
				return nil
			}
			return timestamppb.New(*user.DeletedAt)
		}(),
	}, nil
}
