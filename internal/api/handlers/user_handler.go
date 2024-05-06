package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/repositories"
	"github.com/Christian-007/fit-forge/internal/utils"
	"github.com/jackc/pgx/v5"
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
	users, err := u.UserRepository.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := domains.CollectionRes[domains.User]{Results: users}
	utils.SendResponse(w, http.StatusOK, res)
}

func (u UserHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	user, err := u.UserRepository.GetOne(userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w,r)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	utils.SendResponse(w, http.StatusOK, user)
}
