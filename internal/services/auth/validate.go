package auth

import (
	"context"
	"fmt"

	errors_pkg "boilerplate/internal/pkg/errors"
	jwt_pkg "boilerplate/internal/pkg/jwt"
)

var errUnauthorized = errors_pkg.NewUnauthorizedError("требуется аутентификация")

func (s *service) Validate(ctx context.Context, req *AuthValidateRequest) (*AuthValidateResponse, error) {
	resp := &AuthValidateResponse{}

	if req.AccessToken != nil {
		claims, err := jwt_pkg.ValidateToken(*req.AccessToken, s.config)
		if err != nil {
			return nil, errUnauthorized
		}

		userID, exists := jwt_pkg.GetUserID(claims)
		if !exists {
			return nil, errUnauthorized
		}

		resp.UserID = &userID

		userName, exists := jwt_pkg.GetUserName(claims)
		if exists {
			resp.UserName = &userName
		}

		user, err := s.usersService.Get(ctx, userID)
		if err != nil {
			return nil, errUnauthorized
		}

		if user.Deleted {
			return nil, errUnauthorized
		}

		return resp, nil
	}

	if req.RefreshToken == nil {
		return nil, errUnauthorized
	}

	claims, err := jwt_pkg.ValidateToken(*req.RefreshToken, s.config)
	if err != nil {
		return nil, errUnauthorized
	}

	userID, exists := jwt_pkg.GetUserID(claims)
	if !exists {
		return nil, errUnauthorized
	}

	resp.UserID = &userID

	userName, exists := jwt_pkg.GetUserName(claims)
	if exists {
		resp.UserName = &userName
	}

	user, err := s.usersService.Get(ctx, userID)
	if err != nil {
		return nil, errUnauthorized
	}

	if user.Deleted {
		return nil, errUnauthorized
	}

	newAccessToken, err := jwt_pkg.GenerateAccessToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("create access token: %w", err)
	}

	resp.AccessToken = &newAccessToken

	newRefreshToken, err := jwt_pkg.GenerateRefreshToken(user.ID, user.Name, s.config)
	if err != nil {
		return nil, fmt.Errorf("create refresh token: %w", err)
	}

	resp.RefreshToken = &newRefreshToken

	return resp, nil
}
