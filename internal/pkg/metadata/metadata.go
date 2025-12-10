package metadata

import "context"

const (
	KeyUserID    = "user_id"
	KeyUserName  = "user_name"
	KeyRequestID = "request_id"
	KeyIP        = "ip"
)

func GetUserID(ctx context.Context) (int, bool) {
	if res, ok := ctx.Value(KeyUserID).(int); ok {
		return res, true
	}
	return 0, false
}

func GetRequestID(ctx context.Context) (string, bool) {
	if res, ok := ctx.Value(KeyRequestID).(string); ok {
		return res, true
	}
	return "", false
}

func GetIP(ctx context.Context) (string, bool) {
	ip := ctx.Value(KeyIP)
	if res, ok := ip.(string); ok {
		return res, true
	}
	return "", false
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, KeyUserID, userID) //nolint:revive,staticcheck
}

func SetUserName(ctx context.Context, userName string) context.Context {
	return context.WithValue(ctx, KeyUserName, userName) //nolint:revive,staticcheck
}

func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, KeyRequestID, requestID) //nolint:revive,staticcheck
}

func SetIP(ctx context.Context, ip string) context.Context {
	if ip == "" {
		return ctx
	}
	return context.WithValue(ctx, KeyIP, ip) //nolint:revive,staticcheck
}
