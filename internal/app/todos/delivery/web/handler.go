package web

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
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

func (t TodoHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
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

	todoResponse, err := t.TodoService.GetAllByUserId(userIdInt)
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

func (t TodoHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Todo ID is invalid"})
		return
	}

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

	todoResponse, err := t.TodoService.GetOneByUserId(userIdInt, todoId)
	if err != nil {
		t.Logger.Error(err.Error())

		if err == apperrors.ErrTodoNotFound {
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Todo not found"})
			return
		}

		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	utils.SendResponse(w, http.StatusOK, todoResponse)
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

	if err = createTodoRequest.Validate(); err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: err.Error()})
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

func (t TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	todoId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusBadRequest, apphttp.ErrorResponse{Message: "Todo ID is invalid"})
		return
	}

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

	err = t.TodoService.Delete(todoId, userIdInt)
	if err != nil {
		if err == apperrors.ErrUserOrTodoNotFound {
			t.Logger.Error(err.Error())
			utils.SendResponse(w, http.StatusNotFound, apphttp.ErrorResponse{Message: "Record not found"})
			return
		}

		t.Logger.Error(err.Error())
		utils.SendResponse(w, http.StatusInternalServerError, apphttp.ErrorResponse{Message: "Internal Server Error"})
		return
	}

	w.WriteHeader(http.StatusOK)
}
