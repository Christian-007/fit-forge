package web

import (
	"log/slog"
	"net/http"

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
