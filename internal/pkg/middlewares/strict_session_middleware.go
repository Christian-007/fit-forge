package middlewares

import (
	"net/http"
	"strings"

	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func StrictSession(authService services.AuthService) func(http.Handler) http.Handler {
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
				if err == apperrors.ErrExpiredToken {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Token is expired"})
					return
				}

				if err == apperrors.ErrInvalidSignature {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Token is invalid"})
					return
				}

				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(claims.Uuid)
			if err != nil {
				if err == apperrors.ErrRedisValueNotInHash {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
					return
				}

				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			ctx := requestctx.WithUserId(r.Context(), authData.UserId)
			ctx = requestctx.WithRole(ctx, authData.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
