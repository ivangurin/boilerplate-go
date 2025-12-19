package auth

import (
	"boilerplate/internal/pkg/grpc"
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) Logout(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := grpc.SetAccessToken(ctx, "", -1); err != nil {
		return nil, fmt.Errorf("set access token:  %w", err)
	}
	if err := grpc.SetRefreshToken(ctx, "", -1); err != nil {
		return nil, fmt.Errorf("set refresh token: %w", err)
	}

	return nil, nil
}
