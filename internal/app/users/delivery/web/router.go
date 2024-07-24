package web

import (
	"github.com/Christian-007/fit-forge/internal/app/users/repositories"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux{
	r := chi.NewRouter()
	userRepositoryPg := repositories.NewUserRepositoryPg(appCtx.Pool)
	userHandler := NewUserHandler(UserHandlerOptions{
		UserService: services.NewUserService(services.UserServiceOptions{
			UserRepository: userRepositoryPg,
		}),
		Logger: appCtx.Logger,
	})

	r.Get("/", userHandler.GetAll)
	r.Get("/{id}", userHandler.GetOne)
	r.Post("/", userHandler.Create)
	r.Delete("/{id}", userHandler.Delete)
	r.Patch("/{id}", userHandler.UpdateOne)

	return r
}
