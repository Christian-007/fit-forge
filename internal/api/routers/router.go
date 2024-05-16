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
	userRepository := repositories.NewUserRepository(appCtx.Pool)
	userHandler := handlers.NewUserHandler(handlers.UserHandlerOptions{
		UserService: services.NewUserService(services.UserServiceOptions{
			UserRepository: userRepository,
		}),
		Logger: appCtx.Logger,
	})
	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("GET /users/{id}", userHandler.GetOne)
	mux.HandleFunc("POST /users", userHandler.Create)

	return logRequest(mux)
}
