package ctxdata

import "context"

type contextKey string

const contextKeyUserID = "user_id"

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(contextKeyUserID).(string)
	return userId, ok
}
