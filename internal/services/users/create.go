package users

import (
	"context"
	"fmt"

	"boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/pwd"
	"boilerplate/internal/pkg/utils"
	"boilerplate/internal/repository"
)

func (s *service) Create(ctx context.Context, req *UserCreateRequest) (*User, error) {
	if req.Name == "" {
		return nil, errors.NewBadRequestError("Не указано имя пользователя")
	}
	if req.Email == "" {
		return nil, errors.NewBadRequestError("Не указан email пользователя")
	}
	if req.Password == "" {
		return nil, errors.NewBadRequestError("Не указан пароль пользователя")
	}

	users, err := s.repo.Users().Search(ctx, &repository.UserFilter{
		Emails:      []string{req.Email},
		WithDeleted: utils.Ptr(true),
	})
	if err != nil {
		return nil, fmt.Errorf("search existing users: %w", err)
	}
	if len(users.Result) > 0 {
		return nil, errors.NewBadRequestError("Пользователь с таким email уже существует")
	}

	user := &repository.User{
		Name:  req.Name,
		Email: req.Email,
	}

	if len(req.Password) > 0 {
		var err error
		user.Password, err = pwd.HashPassword(req.Password)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
	}

	id, err := s.repo.Users().Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	user, err = s.repo.Users().Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get created user: %w", err)
	}

	return toUser(user), nil
}
