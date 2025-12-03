package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errors_pkg "boilerplate/internal/pkg/errors"
)

func (s *service) Get(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.Users().Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors_pkg.NewNotFoundError(fmt.Sprintf("Пользователь %d не найден", id))
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	return toUser(user), nil
}
