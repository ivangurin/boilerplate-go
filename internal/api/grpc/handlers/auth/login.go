package auth

import (
	"context"

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

	return &pb.AuthLoginResponse{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}
