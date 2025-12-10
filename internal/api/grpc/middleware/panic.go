package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (m *middleware) Panic(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if e := recover(); e != nil {
			m.logger.Errorf(ctx, "panic: %v", e)
			err = status.Errorf(codes.Internal, "panic: %v", e)
		}
	}()

	return handler(ctx, req)
}
