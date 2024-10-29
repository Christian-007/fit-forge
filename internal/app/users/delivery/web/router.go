package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	userRepositoryPg := repositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})
	userHandler := NewUserHandler(UserHandlerOptions{
		UserService: userService,
		Logger:      appCtx.Logger,
	})
	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		UserService:        userService,
		Cache:              appCtx.RedisClient,
		EnvVariableService: appCtx.EnvVariableService,
	})

	strictSessionMiddleware := middlewares.StrictSession(authService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/", userHandler.Create)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(strictSessionMiddleware)

		r.Get("/", userHandler.GetAll)
		r.Get("/{id}", userHandler.GetOne)
		r.Delete("/{id}", userHandler.Delete)
		r.Patch("/{id}", userHandler.UpdateOne)
	})

	return r
}
