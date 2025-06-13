package requestctx

import (
	"context"
)

func WithEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, EmailContextKey, email)
}

func Email(ctx context.Context) (string, bool) {
	result, ok := ctx.Value(EmailContextKey).(string)
	return result, ok
}
