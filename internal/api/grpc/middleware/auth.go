package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errUnauthenticated
	}

	// Получаем токен из заголовка authorization
	var token string
	values := md.Get("authorization")
	if len(values) > 0 {
		token = strings.TrimPrefix(values[0], "Bearer ")
		token = strings.TrimSpace(token)
	}

	if token == "" {
		return nil, errUnauthenticated
	}

	// Валидируем токен
	authResp, err := m.authService.Validate(ctx, &auth.AuthValidateRequest{
		AccessToken: &token,
	})
	if err != nil {
		return nil, errUnauthenticated
	}

	// Добавляем данные пользователя в контекст
	if authResp.UserID != nil {
		ctx = metadata_pkg.SetUserID(ctx, *authResp.UserID)
	}

	if authResp.UserName != nil {
		ctx = metadata_pkg.SetUserName(ctx, *authResp.UserName)
	}

	return handler(ctx, req)
}
