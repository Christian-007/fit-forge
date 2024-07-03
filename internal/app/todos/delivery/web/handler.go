package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apphttp"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type TodoHandler struct {
	TodoHandlerOptions
}

type TodoHandlerOptions struct {
	TodoService services.TodoService
	Logger      *slog.Logger
}

func NewTodoHandler(options TodoHandlerOptions) TodoHandler {
	return TodoHandler{
		options,
	}
}

func (t TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	todoResponse, err := t.TodoService.GetAll()
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	res := apphttp.CollectionRes[dto.TodoResponse]{
		Results: todoResponse,
	}
	utils.SendResponse(w, http.StatusOK, res)
}

func (t TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("userId")
	if userId == "" {
		t.Logger.Error("User ID is empty")
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "User ID is required"})
		return
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		t.Logger.Error("User ID " + userId + " is invalid")
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "User ID is invalid"})
		return
	}

	var createTodoRequest dto.CreateTodoRequest
	err = json.NewDecoder(r.Body).Decode(&createTodoRequest)
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	todoResponse, err := t.TodoService.Create(userIdInt, createTodoRequest)
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, todoResponse)
}
