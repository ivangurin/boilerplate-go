package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errors_pkg "boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/pwd"
)

func (s *service) Update(ctx context.Context, req *UserUpdateRequest) (*User, error) {
	user, err := s.repo.Users().Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors_pkg.NewNotFoundError(fmt.Sprintf("Пользователь %d не найден", req.ID))
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Password != nil {
		var err error
		user.Password, err = pwd.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("hash password: %w", err)
		}
	}

	err = s.repo.Users().Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	return toUser(user), nil
}
