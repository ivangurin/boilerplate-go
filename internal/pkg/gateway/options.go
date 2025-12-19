package gateway

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

const (
	cookieKey = "cookie"
)

// WithMetadata прокидывает данных из HTTP запроса в gRPC metadata
func WithMetadata(_ context.Context, req *http.Request) metadata.MD {
	md := metadata.MD{}

	cookies := req.Cookies()
	for _, cookie := range cookies {
		md.Set(cookieKey+"-"+strings.ToLower(cookie.Name), cookie.Value)
	}

	return md
}

// WithForwardResponseOption прокидывает данные из gRPC metadata в HTTP ответ
func WithForwardResponseOption(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	cookies := md.HeaderMD.Get(cookieKey)
	for _, cookie := range cookies {
		w.Header().Add("Set-Cookie", cookie)
	}

	return nil
}
