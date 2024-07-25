package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func NewAuthenticate(authService services.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header["Authorization"]
			if authHeader == nil {
				utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
				return
			}

			bearerToken := strings.Fields(authHeader[0])
			if len(bearerToken) < 2 || bearerToken[0] != "Bearer" {
				utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
				return
			}

			token := bearerToken[1]
			claims, err := authService.ValidateToken(token)
			if err != nil {
				utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Token is invalid"})
				return
			}

			ctx := context.WithValue(r.Context(), requestctx.UserContextKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
