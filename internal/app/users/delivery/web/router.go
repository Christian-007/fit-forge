package web

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/app/users/repositories"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
)

func Routes(mux *http.ServeMux, appCtx appcontext.AppContext) {
	userRepositoryPg := repositories.NewUserRepositoryPg(appCtx.Pool)
	userHandler := NewUserHandler(UserHandlerOptions{
		UserService: services.NewUserService(services.UserServiceOptions{
			UserRepository: userRepositoryPg,
		}),
		Logger: appCtx.Logger,
	})

	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("GET /users/{id}", userHandler.GetOne)
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("DELETE /users/{id}", userHandler.Delete)
	mux.HandleFunc("PATCH /users/{id}", userHandler.UpdateOne)
}
