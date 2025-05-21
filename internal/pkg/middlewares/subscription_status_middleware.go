package middlewares

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func SubscriptionStatus() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			subscriptionStatus, ok := requestctx.SubscriptionStatus(ctx)
			if !ok || subscriptionStatus != domains.ActiveSubscriptionStatus {
				utils.SendResponse(w, http.StatusForbidden, apphttp.ErrorResponse{Message: "Forbidden access"})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
