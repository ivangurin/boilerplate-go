package users

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"boilerplate/internal/pkg/convert"
	"boilerplate/internal/services/users"
	"boilerplate/pkg/pb"
)

func ToUser(user *users.User) *pb.User {
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
	}
}
