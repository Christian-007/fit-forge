package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	todorepositories "github.com/Christian-007/fit-forge/internal/app/todos/repositories"
	todoservices "github.com/Christian-007/fit-forge/internal/app/todos/services"
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
		Logger:    appCtx.Logger,
		Publisher: appCtx.Publisher,
	})

	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		Cache:       appCtx.RedisClient,
	})

	strictSessionMiddleware := middlewares.StrictSession(authService, appCtx.SecretManagerClient)
	subscriptionStatusMiddleware := middlewares.SubscriptionStatus()

	// All routes require auth session check
	r.Use(strictSessionMiddleware)

	// Routes that can be accessed by all roles
	r.Group(func(r chi.Router) {
		// Check if the user's subscription status is 'ACTIVE'
		r.Use(subscriptionStatusMiddleware)

		r.Get("/", todoHandler.GetAllByUserId)
		r.Get("/{id}", todoHandler.GetOne)
		r.Post("/", todoHandler.Create)
		r.Delete("/{id}", todoHandler.Delete)
		r.Patch("/{id}", todoHandler.Patch)
	})

	// Routes that can only be accessed by Admin role
	adminRoleEnum := 1
	adminRoleMiddleware := middlewares.Role(adminRoleEnum)
	r.Group(func(r chi.Router) {
		r.Use(adminRoleMiddleware)
		r.Get("/all", todoHandler.GetAll)
	})

	return r
}
