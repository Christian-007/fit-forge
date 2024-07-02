package web

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/todos/repositories"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
)

func Routes(mux *http.ServeMux, appCtx appcontext.AppContext) {
	todoRepository := repositories.NewTodoRepositoryPg(appCtx.Pool)
	todoHandler := NewTodoHandler(TodoHandlerOptions{
		TodoService: services.NewTodoService(services.TodoServiceOptions{
			TodoRepository: todoRepository,
		}),
	})

	mux.HandleFunc("GET /todos", todoHandler.GetAll)
}
