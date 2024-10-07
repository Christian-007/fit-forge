package requestctx

import "context"

func WithUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, UserContextKey, userId)
}

func UserId(ctx context.Context) (int, bool) {
	result, ok := ctx.Value(UserContextKey).(int)
	return result, ok
}
