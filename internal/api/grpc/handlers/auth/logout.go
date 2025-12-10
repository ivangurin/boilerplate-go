package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *handler) Logout(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}
