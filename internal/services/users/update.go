package users

import (
	"context"
	"fmt"
)

func (s *Service) Update(ctx context.Context, req *UserUpdateRequest) (*User, error) {
	user, err := s.repo.Users().Get(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("get user: %s", err.Error())
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.Password != nil {
		user.Password = *req.Password
	}

	err = s.repo.Users().Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update user: %s", err.Error())
	}

	return toUser(user), nil
}
