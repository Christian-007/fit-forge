package requestctx

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
)

func WithSubscriptionStatus(ctx context.Context, status domains.SubscriptionStatus) context.Context {
	return context.WithValue(ctx, SubscriptionStatusContextKey, status)
}

func SubscriptionStatus(ctx context.Context) (domains.SubscriptionStatus, bool) {
	result, ok := ctx.Value(SubscriptionStatusContextKey).(domains.SubscriptionStatus)
	return result, ok
}
