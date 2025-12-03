package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	errors_pkg "boilerplate/internal/pkg/errors"
)

func (s *service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Users().Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors_pkg.NewNotFoundError(fmt.Sprintf("Пользователь %d не найден", id))
		}
		return fmt.Errorf("get user: %w", err)
	}

	err = s.repo.Users().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
