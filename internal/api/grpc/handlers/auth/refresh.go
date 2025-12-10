package auth

import (
	"context"

	"boilerplate/internal/pkg/grpc"
	"boilerplate/internal/services/auth"
	"boilerplate/pkg/pb"
)

func (h *handler) Refresh(ctx context.Context, req *pb.AuthRefreshRequest) (*pb.AuthRefreshResponse, error) {
	resp, err := h.authService.Refresh(ctx, &auth.AuthRefreshRequest{
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		return nil, grpc.Error(err)
	}

	return &pb.AuthRefreshResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
