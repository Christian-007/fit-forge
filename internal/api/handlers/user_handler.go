package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Christian-007/fit-forge/internal/api/repositories"
)

type UserHandler struct {
	UserHandlerOptions
}

type UserHandlerOptions struct {
	UserRepository repositories.UserRepository
}

func NewUserHandler(options UserHandlerOptions) UserHandler {
	return UserHandler{
		options,
	}
}

func (u UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users := u.UserRepository.GetAll()
	jsonRes, err := json.Marshal(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}
