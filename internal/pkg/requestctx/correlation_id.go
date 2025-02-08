package requestctx

import (
	"context"

	"github.com/lithammer/shortuuid/v4"
)

func WithCorrelationId(ctx context.Context, correlationId string) context.Context {
	return context.WithValue(ctx, CorrelationIdContextKey, correlationId)
}

func CorrelationId(ctx context.Context) (string, bool) {
	result, ok := ctx.Value(CorrelationIdContextKey).(string)
	if ok {
		return result, ok
	}

	return "gen_" + shortuuid.New(), ok
}
