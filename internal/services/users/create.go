package users

import (
	"context"
	"fmt"

	"boilerplate/internal/repository"
)

func (s *service) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	user := &repository.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	id, err := s.repo.Users().Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %s", err.Error())
	}

	user, err = s.repo.Users().Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created user: %s", err.Error())
	}

	return toUser(user), nil
}
