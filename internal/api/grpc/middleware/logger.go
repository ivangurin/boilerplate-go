package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		m.logger.Errorf(ctx, "response: method: %v, err: %v (not a gRPC error)", info.FullMethod, err)
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.Internal {
				return nil, status.Error(codes.Internal, "internal server error. Please try again later")
			}
			return resp, err
		} else {
			return resp, err
		}
	}

	if m.logger.IsWithDebug() {
		raw, err = protojson.Marshal((resp).(proto.Message)) //nolint:errcheck
		if err != nil {
			raw = []byte("<failed to marshal response>")
		}
		m.logger.Debugf(ctx, "response: method: %s, resp: %s", info.FullMethod, string(raw))
	}

	return resp, nil
}
