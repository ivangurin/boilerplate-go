package auth

import (
	"context"
	"errors"
	"fmt"

	errors_pkg "boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/jwt"
	"boilerplate/internal/pkg/pwd"
	"boilerplate/internal/pkg/utils"
	users_service "boilerplate/internal/services/users"
)

func (s *service) Login(ctx context.Context, req *AuthLoginRequest) (*AuthLoginResponse, error) {
	users, err := s.usersService.Search(ctx, &users_service.UserSearchRequest{
		Filter: users_service.UserSearchRequestFilter{
			Email:       []string{req.Email},
			WithDeleted: utils.Ptr(true),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("поиск пользователя: %w", err)
	}
	if len(users.Result) == 0 {
		return nil, errors_pkg.NewNotFoundError("пользователь не найден")
	}
	if len(users.Result) > 1 {
		return nil, errors.New("найдено несколько пользователей с одинаковым логином")
	}

	user := users.Result[0]

	if !pwd.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors_pkg.NewUnauthorizedError("неверные учетные данные")
	}

	if user.Deleted {
		return nil, errors_pkg.NewForbiddenError("пользователь удален")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("генерация токена доступа: %w", err)
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("генерация токена обновления: %w", err)
	}

	res := &AuthLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	return res, nil
}
