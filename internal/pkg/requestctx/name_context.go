package requestctx

import (
	"context"
)

func WithName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, NameContextKey, name)
}

func Name(ctx context.Context) (string, bool) {
	result, ok := ctx.Value(NameContextKey).(string)
	return result, ok
}
