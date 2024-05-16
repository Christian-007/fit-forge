package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/api/apperrors"
	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/dto"
	"github.com/Christian-007/fit-forge/internal/api/services"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type UserHandler struct {
	UserHandlerOptions
}

type UserHandlerOptions struct {
	UserService    services.UserService
	Logger         *slog.Logger
}

func NewUserHandler(options UserHandlerOptions) UserHandler {
	return UserHandler{
		options,
	}
}

func (u UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userResponses, err := u.UserService.GetAll()
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, domains.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	res := domains.CollectionRes[dto.UserResponse]{Results: userResponses}
	utils.SendResponse(w, http.StatusOK, res)
}

func (u UserHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusNotFound, domains.ErrorResponse{Message: "Record not found"})
		return
	}

	user, err := u.UserService.GetOne(userId)
	if err != nil {
		if err == apperrors.ErrUserNotFound {
			u.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, domains.ErrorResponse{Message: "Record not found"})
			return
		}

		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, domains.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, user)
}

func (u UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, domains.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	if err = createUserRequest.Validate(); err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, domains.ErrorResponse{Message: err.Error()})
		return
	}

	userResponse, err := u.UserService.Create(createUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, domains.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, userResponse)
}
