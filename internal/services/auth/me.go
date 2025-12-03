package auth

import (
	"context"

	errors_pkg "boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/metadata"
	users_service "boilerplate/internal/services/users"
)

func (s *service) Me(ctx context.Context) (*users_service.User, error) {
	userID, exists := metadata.GetUserID(ctx)
	if !exists {
		return nil, errors_pkg.NewUnauthorizedError("Не авторизованы")
	}

	user, err := s.usersService.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
