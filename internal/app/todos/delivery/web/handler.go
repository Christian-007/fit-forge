package web

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/todos/delivery/domains"
	"github.com/Christian-007/fit-forge/internal/utils"
)

type TodoHandler struct{}

func NewTodoHandler() TodoHandler {
	return TodoHandler{}
}

func (t TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(w, http.StatusOK, domains.TodoModel{
		Id: 1,
		Title: "Test",
		IsCompleted: true,
		UserId: 1,
	})
}
