package auth

import (
	"context"
	"errors"
	"fmt"

	errors_pkg "boilerplate/internal/pkg/errors"
	jwt_pkg "boilerplate/internal/pkg/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func (s *service) Refresh(ctx context.Context, req *AuthRefreshRequest) (*AuthRefreshResponse, error) {
	if len(req.RefreshToken) == 0 {
		return nil, errors_pkg.NewBadRequestError("не указан токен обновления")
	}

	token, claims, err := jwt_pkg.ParseToken(req.RefreshToken, s.config)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors_pkg.NewUnauthorizedError("неправильная подпись токена")
		}
		return nil, errors_pkg.NewUnauthorizedError(fmt.Sprintf("разбор токена: %s", err.Error()))
	}

	if !token.Valid {
		return nil, errors_pkg.NewUnauthorizedError("недействительный токен")
	}

	userID, exists := jwt_pkg.GetUserID(claims)
	if !exists {
		return nil, errors_pkg.NewUnauthorizedError("не указан код пользователя")
	}

	user, err := s.usersService.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user.Deleted {
		return nil, errors_pkg.NewForbiddenError("пользователь удален")
	}

	newAccessToken, err := jwt_pkg.GenerateAccessToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("генерация токена доступа: %w", err)
	}

	newRefreshToken, err := jwt_pkg.GenerateRefreshToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("генерация токена обновления: %w", err)
	}

	res := &AuthRefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	return res, nil
}
