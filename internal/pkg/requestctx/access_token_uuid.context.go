package requestctx

import "context"

func WithAccessTokenUuid(ctx context.Context, accessTokenUuid string) context.Context {
	return context.WithValue(ctx, AccessTokenUuidContextKey, accessTokenUuid)
}

func AccessTokenUuid(ctx context.Context) (string, bool) {
	result, ok := ctx.Value(AccessTokenUuidContextKey).(string)
	return result, ok
}
