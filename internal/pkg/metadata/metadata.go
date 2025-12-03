package metadata

import "context"

const (
	KeyUserID   = "user_id"
	KeyUserName = "user_name"
)

func GetUserID(ctx context.Context) (int, bool) {
	if userID, ok := ctx.Value(KeyUserID).(int); ok {
		return userID, true
	}
	return 0, false
}
