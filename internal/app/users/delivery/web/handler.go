package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type UserHandler struct {
	UserHandlerOptions
}

type UserHandlerOptions struct {
	UserService services.UserService
	Logger      applog.Logger
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
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	res := apphttp.CollectionRes[dto.UserResponse]{Results: userResponses}
	utils.SendResponse(w, http.StatusOK, res)
}

func (u UserHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
		return
	}

	user, err := u.UserService.GetOne(userId)
	if err != nil {
		if err == apperrors.ErrUserNotFound {
			u.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, user)
}

func (u UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var createUserRequest dto.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	if err = createUserRequest.Validate(); err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: err.Error()})
		return
	}

	userResponse, err := u.UserService.Create(createUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, userResponse)
}

func (u UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Invalid ID"})
		return
	}

	err = u.UserService.Delete(userId)
	if err != nil {
		if err == apperrors.ErrUserNotFound {
			u.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u UserHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Invalid ID"})
		return
	}

	var updateUserRequest dto.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&updateUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	if err = updateUserRequest.Validate(); err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: err.Error()})
		return
	}

	user, err := u.UserService.UpdateOne(userId, updateUserRequest)
	if err != nil {
		if err == apperrors.ErrUserNotFound {
			u.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, user)
}
