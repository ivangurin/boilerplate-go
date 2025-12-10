package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func (m *middleware) Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	raw, err := protojson.Marshal((req).(proto.Message)) //nolint:errcheck
	if err != nil {
		raw = []byte("<failed to marshal request>")
	}
	m.logger.Debugf(ctx, "request: method: %s, req: %s", info.FullMethod, string(raw))

	resp, err := handler(ctx, req)
	if err != nil {
		m.logger.Errorf(ctx, "response: method: %v, err: %v", info.FullMethod, err)
		return resp, err
	}

	raw, err = protojson.Marshal((resp).(proto.Message)) //nolint:errcheck
	if err != nil {
		raw = []byte("<failed to marshal response>")
	}
	m.logger.Debugf(ctx, "response: method: %s, resp: %s", info.FullMethod, string(raw))

	return resp, nil
}
