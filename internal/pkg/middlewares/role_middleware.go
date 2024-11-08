package middlewares

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func Role(role int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userRole, ok := requestctx.Role(ctx)
			if !ok {
				utils.SendResponse(w, http.StatusForbidden, apphttp.ErrorResponse{Message: "Forbidden access"})
				return
			}

			hasCorrectRole := userRole == role
			if !hasCorrectRole {
				utils.SendResponse(w, http.StatusForbidden, apphttp.ErrorResponse{Message: "Forbidden access"})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}