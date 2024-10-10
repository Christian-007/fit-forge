package middlewares

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func LogoutSession(authService services.AuthService) func (http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {	
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			accessTokenUuid, ok := requestctx.AccessTokenUuid(ctx)
			if !ok {
				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(accessTokenUuid)
			if err != nil {
				if err == apperrors.ErrRedisValueNotInHash {
					// Handle auth using old data structure
					userId, err := authService.GetAuthDataFromCache(accessTokenUuid)
					if err != nil {
						utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
						return
					}
					
					ctx = requestctx.WithUserId(ctx, userId)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				if err == apperrors.ErrRedisKeyNotFound {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
					return
				}
	
				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
			}

			ctx = requestctx.WithUserId(ctx, authData.UserId)
			ctx = requestctx.WithRole(ctx, authData.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}