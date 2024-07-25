package web

import (
	"github.com/Christian-007/fit-forge/internal/app/todos/repositories"
	"github.com/Christian-007/fit-forge/internal/app/todos/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	todoRepository := repositories.NewTodoRepositoryPg(appCtx.Pool)
	todoHandler := NewTodoHandler(TodoHandlerOptions{
		TodoService: services.NewTodoService(services.TodoServiceOptions{
			TodoRepository: todoRepository,
		}),
		Logger: appCtx.Logger,
	})

	r.Get("/all", todoHandler.GetAll)
	r.Get("/", todoHandler.GetAllByUserId)
	r.Get("/{id}", todoHandler.GetOne)
	r.Post("/", todoHandler.Create)
	r.Delete("/{id}", todoHandler.Delete)

	return r
}
