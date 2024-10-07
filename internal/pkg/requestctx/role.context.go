package requestctx

import "context"

func WithRole(ctx context.Context, role int) context.Context {
	return context.WithValue(ctx, UserRoleContextKey, role)
}

func Role(ctx context.Context) (int, bool) {
	result, ok := ctx.Value(UserRoleContextKey).(int)
	return result, ok
}
