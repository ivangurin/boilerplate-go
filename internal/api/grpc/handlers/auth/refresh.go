package auth

import (
	"context"
	"fmt"

	"boilerplate/internal/pkg/grpc"
	"boilerplate/internal/services/auth"
	"boilerplate/pkg/pb"
)

func (h *handler) Refresh(ctx context.Context, req *pb.AuthRefreshRequest) (*pb.AuthRefreshResponse, error) {
	refreshToken := req.GetRefreshToken()
	if refreshToken == "" {
		refreshToken, _ = grpc.GetRefreshToken(ctx)
	}

	resp, err := h.authService.Refresh(ctx, &auth.AuthRefreshRequest{
		RefreshToken: refreshToken,
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

	return &pb.AuthRefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
