package middleware

import (
	"boilerplate/internal/pkg/metadata"
	"boilerplate/internal/pkg/utils"
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func (m *middleware) Tracer(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = metadata.SetRequestID(ctx, utils.UniqueID())
	ctx = metadata.SetIP(ctx, getClientIP(ctx))

	return handler(ctx, req)
}

func getClientIP(ctx context.Context) string {
	p, exists := peer.FromContext(ctx)
	if !exists {
		return ""
	}

	addr := p.Addr.String()

	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}

	return addr
}
