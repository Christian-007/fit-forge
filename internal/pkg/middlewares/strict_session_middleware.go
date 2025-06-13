package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/Christian-007/fit-forge/internal/utils"
)

func StrictSession(authService services.AuthService, secretManagerProvider security.SecretManageProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var uuid string

			// Method 1: GCP API Gateway
			apiGatewayUserInfoHeader := r.Header.Get("X-Apigateway-Api-Userinfo")
			if apiGatewayUserInfoHeader != "" {
				decodedHeader, err := base64.RawURLEncoding.DecodeString(apiGatewayUserInfoHeader)
				if err != nil {
					utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request"})
					return
				}

				var payload map[string]any
				err = json.Unmarshal(decodedHeader, &payload)
				if err != nil {
					utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request: Error Unmarshal"})
					return
				}

				uuid = payload["Uuid"].(string)
			} else {
				// Method 2: Manual JWT validation
				authHeader, ok := r.Header["Authorization"]
				if !ok {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
					return
				}

				bearerToken := strings.Fields(authHeader[0])
				if len(bearerToken) < 2 || bearerToken[0] != "Bearer" {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
					return
				}

				privateKey, err := secretManagerProvider.GetPrivateKey(r.Context(), "")
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
					return
				}

				token := bearerToken[1]
				claims, err := authService.ValidateToken(privateKey, token)
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

				uuid = claims.Uuid
			}

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(uuid)
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
			ctx = requestctx.WithSubscriptionStatus(ctx, authData.SubscriptionStatus)
			ctx = requestctx.WithName(ctx, authData.Name)
			ctx = requestctx.WithEmail(ctx, authData.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
