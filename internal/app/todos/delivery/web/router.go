package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	todorepositories "github.com/Christian-007/fit-forge/internal/app/todos/repositories"
	todoservices "github.com/Christian-007/fit-forge/internal/app/todos/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	todoRepository := todorepositories.NewTodoRepositoryPg(appCtx.Pool)
	todoHandler := NewTodoHandler(TodoHandlerOptions{
		TodoService: todoservices.NewTodoService(todoservices.TodoServiceOptions{
			TodoRepository: todoRepository,
		}),
		Logger: appCtx.Logger,
	})

	userRepositoryPg := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})
	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		UserService:        userService,
		Cache:              appCtx.RedisClient,
		EnvVariableService: appCtx.EnvVariableService,
	})

	strictSessionMiddleware := middlewares.StrictSession(authService)

	r.Use(strictSessionMiddleware)

	r.Get("/all", todoHandler.GetAll)
	r.Get("/", todoHandler.GetAllByUserId)
	r.Get("/{id}", todoHandler.GetOne)
	r.Post("/", todoHandler.Create)
	r.Delete("/{id}", todoHandler.Delete)

	return r
}
