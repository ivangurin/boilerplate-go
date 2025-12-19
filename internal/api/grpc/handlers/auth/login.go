package auth

import (
	"context"
	"fmt"

	"boilerplate/internal/pkg/grpc"
	"boilerplate/internal/services/auth"
	"boilerplate/pkg/pb"
)

func (h *handler) Login(ctx context.Context, req *pb.AuthLoginRequest) (*pb.AuthLoginResponse, error) {
	resp, err := h.authService.Login(ctx, &auth.AuthLoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, grpc.Error(err)
	}

	if err := grpc.SetAccessToken(ctx, resp.AccessToken, h.authService.GetConfig().AccessTokenTTL); err != nil {
		return nil, fmt.Errorf("set access token:  %w", err)
	}
	if err := grpc.SetRefreshToken(ctx, resp.RefreshToken, h.authService.GetConfig().RefreshTokenTTL); err != nil {
		return nil, fmt.Errorf("set refresh token: %w", err)
	}

	return &pb.AuthLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
