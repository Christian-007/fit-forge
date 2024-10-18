package middlewares

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func StrictSession(authService services.AuthServiceImpl) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			accessTokenUuid, ok := requestctx.AccessTokenUuid(r.Context())
			if !ok {
				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(accessTokenUuid)
			if err != nil {
				if err == apperrors.ErrRedisKeyNotFound || err == apperrors.ErrRedisValueNotInHash {
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
