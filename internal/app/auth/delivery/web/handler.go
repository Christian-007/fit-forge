package web

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/auth/dto"
	"github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type AuthHandler struct {
	AuthHandlerOptions
}

type AuthHandlerOptions struct {
	AuthService services.AuthService
	Logger      *slog.Logger
}

func NewAuthHandler(options AuthHandlerOptions) AuthHandler {
	return AuthHandler{
		options,
	}
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	userResponse, err := a.AuthService.Authenticate(loginRequest)
	if err != nil {
		a.Logger.Error(err.Error())

		if err == apperrors.ErrUserNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		if err == apperrors.ErrInvalidCredentials {
			utils.SendResponse(w, http.StatusUnauthorized, apphttp.ErrorResponse{Message: "Invalid username or password"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	token, err := a.AuthService.CreateToken(userResponse.Id)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	err = a.AuthService.SaveToken(userResponse.Id, token)
	if err != nil {
		a.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, dto.LoginResponse{AccessToken: token.AccessToken})
}
