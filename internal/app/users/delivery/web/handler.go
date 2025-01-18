package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	emaildomain "github.com/Christian-007/fit-forge/internal/app/email/domains"
	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
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
	UserService    services.UserService
	Logger         applog.Logger
	EmailService   emailservices.EmailService
	MailtrapSender emailservices.MailtrapSender
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

	// TODO: need to findUserByEmail() first, to check if the user
	// has registered before, but didn't verify the email
	// If yes, then proceed to overwrite the old with the new data
	// and proceed to send the email verification again.

	userResponse, err := u.UserService.Create(createUserRequest)
	if err != nil {
		u.Logger.Error(err.Error())

		if err == apperrors.ErrEmailDuplicate {
			utils.SendResponse(w, http.StatusConflict, apphttp.ErrorResponse{Message: "Email has already been registered"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}
	u.Logger.Info("Succefully created a user", slog.String("email", userResponse.Email))

	verificationLink, err := u.EmailService.CreateVerificationLink(userResponse.Email)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Cannot generate an email verification link"})
		return
	}

	emailRequest := emaildomain.EmailWithTemplateRequest{
		From: emaildomain.EmailAddressOptions{
			Email: "hello@demomailtrap.com",
			Name:  "No Reply at Fit Forge",
		},
		To: []emaildomain.EmailAddressOptions{
			{
				Email: userResponse.Email,
				Name:  userResponse.Name,
			},
		},
		TemplateUuid: "fdbefad8-2410-45d2-bded-9d1b647ac416",
		TemplateVariables: map[string]any{
			"user_name":         userResponse.Name,
			"verification_link": verificationLink,
		},
	}
	err = u.MailtrapSender.SendWithTemplate(emailRequest)
	if err != nil {
		u.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Cannot send an email verification"})
		return
	}
	u.Logger.Info("Succefully send a verification email", slog.String("email", userResponse.Email))

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
