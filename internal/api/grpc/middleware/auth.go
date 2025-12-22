package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_pkg "boilerplate/internal/pkg/grpc"
	metadata_pkg "boilerplate/internal/pkg/metadata"
	"boilerplate/internal/services/auth"
)

var errUnauthenticated = status.Error(codes.Unauthenticated, "требуется аутентификация")

// Список публичных методов, не требующих авторизации
var publicMethods = map[string]bool{
	"/auth.AuthAPI/Login":    true,
	"/auth.AuthAPI/Refresh":  true,
	"/users.UsersAPI/Create": true,
}

// nolint:revive
func (m *middleware) Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	// Проверяем, является ли метод публичным
	if publicMethods[info.FullMethod] {
		return handler(ctx, req)
	}

	// Получаем токены
	accessToken, _ := grpc_pkg.GetAccessToken(ctx)
	refreshToken, _ := grpc_pkg.GetRefreshToken(ctx)
	if accessToken == "" && refreshToken == "" {
		return nil, errUnauthenticated
	}

	authReq := &auth.AuthValidateRequest{}
	if accessToken != "" {
		authReq.AccessToken = &accessToken
	}
	if refreshToken != "" {
		authReq.RefreshToken = &refreshToken
	}

	// Валидируем токен
	authResp, err := m.authService.Validate(ctx, authReq)
	if err != nil {
		return nil, errUnauthenticated
	}

	// Добавляем данные пользователя в контекст
	if authResp.UserID != nil {
		ctx = metadata_pkg.WithUserID(ctx, *authResp.UserID)
	}

	if authResp.AccessToken != nil {
		if err := grpc_pkg.SetAccessToken(ctx, *authResp.AccessToken, m.authService.GetConfig().AccessTokenTTL); err != nil {
			return nil, err
		}
	}

	if authResp.RefreshToken != nil {
		if err := grpc_pkg.SetRefreshToken(ctx, *authResp.RefreshToken, m.authService.GetConfig().RefreshTokenTTL); err != nil {
			return nil, err
		}
	}

	return handler(ctx, req)
}
