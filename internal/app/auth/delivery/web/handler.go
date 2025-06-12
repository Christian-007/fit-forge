package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	authdto "github.com/Christian-007/fit-forge/internal/app/auth/dto"
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
	usersdomain "github.com/Christian-007/fit-forge/internal/app/users/domains"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type AuthHandler struct {
	AuthHandlerOptions
}

type AuthHandlerOptions struct {
	AuthService         authservices.AuthServiceImpl
	Logger              applog.Logger
	EmailService        emailservices.EmailService
	UserService         userservices.UserService
	Cache               cache.Cache
	SecretManagerClient security.SecretManageProvider
}

func NewAuthHandler(options AuthHandlerOptions) AuthHandler {
	return AuthHandler{
		options,
	}
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest authdto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}
	userResponse, err := a.UserService.GetOneByEmail(loginRequest.Username)
	if err != nil {
		a.Logger.Error(err.Error())

		if err == apperrors.ErrUserNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	err = a.AuthService.Authenticate(loginRequest.Password, userResponse.Password)
	if err != nil {
		a.Logger.Error(err.Error())

		if err == apperrors.ErrInvalidCredentials {
			utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Invalid username or password"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	privateKey, err := a.SecretManagerClient.GetPrivateKey(r.Context(), os.Getenv("GCP_SECRET_DIR"))
	if err != nil {
		a.Logger.Error("failed to get private key", slog.String("error", err.Error()))
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	token, err := a.AuthService.CreateToken(privateKey, userResponse.Id)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	dto := userdto.UserResponse{
		Id:                 userResponse.Id,
		Name:               userResponse.Name,
		Email:              userResponse.Email,
		Role:               userResponse.Role,
		SubscriptionStatus: userResponse.SubscriptionStatus,
	}
	err = a.AuthService.SaveToken(dto, token)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, authdto.LoginResponse{AccessToken: token.AccessToken})
}

func (a AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	accessTokenUuid, ok := requestctx.AccessTokenUuid(r.Context())
	if !ok {
		a.Logger.Error(apperrors.ErrTypeAssertion.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	err := a.AuthService.InvalidateToken(accessTokenUuid)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	if token == "" {
		a.Logger.Error("No Token Found on Email Verify")
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request"})
		return
	}

	email, hashedToken, err := a.EmailService.Verify(token)
	if err != nil {
		a.Logger.Error(err.Error())

		if err == apperrors.ErrRedisKeyNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Token not found"})
			return
		}

		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Bad Request"})
		return
	}

	currentTime := time.Now()
	subscriptionStatus := usersdomain.ActiveSubscriptionStatus
	updateUserRequest := userdto.UpdateUserRequest{
		EmailVerifiedAt:    &currentTime,
		SubscriptionStatus: &subscriptionStatus,
	}

	userResponse, err := a.UserService.UpdateOneByEmail(email, updateUserRequest)
	if err != nil {
		a.Logger.Error(err.Error())

		if err == apperrors.ErrUserNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	err = a.Cache.Delete(hashedToken)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	a.Logger.Info("Email verified successfully", "email", email)
	utils.SendResponse(w, http.StatusOK, userResponse)
}
