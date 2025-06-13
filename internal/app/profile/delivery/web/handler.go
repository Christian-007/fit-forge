package web

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/profile/dto"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/requestctx"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type ProfileHandler struct {
	ProfileHandlerOptions
}

type ProfileHandlerOptions struct {
	Logger applog.Logger
}

func NewProfileHandler(options ProfileHandlerOptions) ProfileHandler {
	return ProfileHandler{
		options,
	}
}

func (p ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	profileData, err := getProfileDataFromContext(r)
	if err != nil {
		p.Logger.Error("failed to get profile data from context", slog.String("error", err.Error()))
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
	}
	
	utils.SendResponse(w, http.StatusOK, profileData)
}

func getProfileDataFromContext(r *http.Request) (dto.ProfileResponse, error) {
	userId, ok := requestctx.UserId(r.Context())
	if !ok {
		return dto.ProfileResponse{}, fmt.Errorf("empty 'userId' value")
	}

	role, ok := requestctx.Role(r.Context())
	if !ok {
		return dto.ProfileResponse{}, fmt.Errorf("empty 'role' value")
	}

	subscriptionStatus, ok := requestctx.SubscriptionStatus(r.Context())
	if !ok {
		return dto.ProfileResponse{}, fmt.Errorf("empty 'subscriptionStatus' value")
	}

	email, ok := requestctx.Email(r.Context())
	if !ok {
		return dto.ProfileResponse{}, fmt.Errorf("empty 'email' value")
	}

	name, ok := requestctx.Name(r.Context())
	if !ok {
		return dto.ProfileResponse{}, fmt.Errorf("empty 'name' value")
	}

	return dto.ProfileResponse{
		UserId: userId,
		Role: role,
		SubscriptionStatus: subscriptionStatus,
		Email: email,
		Name: name,
	}, nil

}
