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

func LogoutSession(authService services.AuthService, secretManagerProvider security.SecretManageProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var uuid string
			ctx := r.Context()

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

				privateKey, err := secretManagerProvider.GetPrivateKey(ctx, "")
				if err != nil {
					utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
					return
				}

				token := bearerToken[1]
				claims, err := authService.ValidateToken(privateKey, token)
				if err != nil {
					if err == apperrors.ErrExpiredToken {
						utils.SendResponse(w, http.StatusOK, apphttp.ErrorResponse{Message: "Logout successful"})
						return
					}

					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Token is invalid"})
					return
				}

				uuid = claims.Uuid
			}

			// It's important to use `userId` from the cache just in case the JWT has been tampered
			authData, err := authService.GetHashAuthDataFromCache(uuid)
			if err != nil {
				if err == apperrors.ErrRedisValueNotInHash {
					// Handle auth using old data structure
					userId, err := authService.GetAuthDataFromCache(uuid)
					if err != nil {
						if err == apperrors.ErrRedisKeyNotFound {
							utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
							return
						}

						utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
						return
					}

					ctx = requestctx.WithAccessTokenUuid(ctx, uuid)
					ctx = requestctx.WithUserId(ctx, userId)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				if err == apperrors.ErrRedisKeyNotFound {
					utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Unauthorized"})
					return
				}

				utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
				return
			}

			ctx = requestctx.WithAccessTokenUuid(ctx, uuid)
			ctx = requestctx.WithUserId(ctx, authData.UserId)
			ctx = requestctx.WithRole(ctx, authData.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
