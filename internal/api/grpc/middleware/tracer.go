package middleware

import (
	"context"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	metadata_pkg "boilerplate/internal/pkg/metadata"
	"boilerplate/internal/pkg/utils"
)

func (m *middleware) Tracer(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	ctx = metadata_pkg.WithRequestID(ctx, utils.UniqueID())
	if ip, exists := getClientIP(ctx); exists {
		ctx = metadata_pkg.WithIP(ctx, ip)
	}

	return handler(ctx, req)
}

func getClientIP(ctx context.Context) (string, bool) {
	// Сначала пытаемся получить IP из заголовков, установленных nginx
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// Сначала проверяем X-Real-IP
		if xri := md.Get("x-real-ip"); len(xri) > 0 && xri[0] != "" {
			return xri[0], true
		}

		// Затем проверяем X-Forwarded-For (может содержать несколько IP через запятую)
		if xff := md.Get("x-forwarded-for"); len(xff) > 0 {
			// Берем первый IP из списка (реальный клиент)
			ips := strings.Split(xff[0], ",")
			if len(ips) > 0 {
				ip := strings.TrimSpace(ips[0])
				if ip != "" {
					return ip, true
				}
			}
		}
	}

	// Пытаемся получить IP из peer (если не проксируется)
	p, exists := peer.FromContext(ctx)
	if !exists {
		return "", false
	}

	addr := p.Addr.String()

	if host, _, err := net.SplitHostPort(addr); err == nil {
		addr = host
	}

	return addr, true
}
