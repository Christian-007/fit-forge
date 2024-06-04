package routers

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/handlers"
	"github.com/Christian-007/fit-forge/internal/api/middlewares"
	"github.com/Christian-007/fit-forge/internal/api/repositories"
	"github.com/Christian-007/fit-forge/internal/api/services"
)

func Routes(appCtx domains.AppContext) http.Handler {
	mux := http.NewServeMux()

	logRequest := middlewares.NewLogRequest(appCtx.Logger)
	userRepositoryPg := repositories.NewUserRepositoryPg(appCtx.Pool)
	userHandler := handlers.NewUserHandler(handlers.UserHandlerOptions{
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

	return logRequest(mux)
}
