package grpc

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	cookieKey    = "cookie"
	accessToken  = "access_token"
	refreshToken = "refresh_token"
)

// getCookie извлекает значение cookie из gRPC metadata
func getCookie(ctx context.Context, name string) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	key := cookieKey + "-" + strings.ToLower(name)
	values := md.Get(key)
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}

// setCookie устанавливает cookie в gRPC metadata для отправки в HTTP ответ
func setCookie(ctx context.Context, cookie *http.Cookie) error {
	md := metadata.Pairs(cookieKey, cookie.String())

	if err := grpc.SetHeader(ctx, md); err != nil {
		return fmt.Errorf("set cookie header: %w", err)
	}

	return nil
}

func SetAccessToken(ctx context.Context, token string, ttl int) error {
	return setCookie(ctx, &http.Cookie{
		Name:     accessToken,
		Value:    token,
		Path:     "/",
		MaxAge:   ttl,
		Secure:   false,
		HttpOnly: true,
	})
}

func GetAccessToken(ctx context.Context) (string, bool) {
	var token string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get("authorization")
		if len(values) > 0 {
			token = strings.TrimPrefix(values[0], "Bearer ")
			token = strings.TrimSpace(token)
		}
	}
	if token != "" {
		return token, true
	}

	return getCookie(ctx, accessToken)
}

func SetRefreshToken(ctx context.Context, token string, ttl int) error {
	return setCookie(ctx, &http.Cookie{
		Name:     refreshToken,
		Value:    token,
		Path:     "/",
		MaxAge:   ttl,
		Secure:   false,
		HttpOnly: true,
	})
}

func GetRefreshToken(ctx context.Context) (string, bool) {
	return getCookie(ctx, refreshToken)
}
