package requestctx

type contextKey string

const (
	UserContextKey               = contextKey("userId")
	AccessTokenUuidContextKey    = contextKey("accessTokenUuid")
	UserRoleContextKey           = contextKey("role")
	CorrelationIdContextKey      = contextKey("correlationId")
	SubscriptionStatusContextKey = contextKey("subscriptionStatus")
)
