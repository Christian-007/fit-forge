package routers

import (
	"net/http"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/handlers"
	"github.com/Christian-007/fit-forge/internal/api/middlewares"
	"github.com/Christian-007/fit-forge/internal/api/repositories"
)

func Routes(appCtx domains.AppContext) http.Handler {
	mux := http.NewServeMux()

	logRequest := middlewares.NewLogRequest(appCtx.Logger)
	userHandler := handlers.NewUserHandler(handlers.UserHandlerOptions{
		UserRepository: repositories.NewUserRepository(appCtx.Pool),
		Logger: appCtx.Logger,
	})
	mux.HandleFunc("GET /users", userHandler.GetAll)
	mux.HandleFunc("GET /users/{id}", userHandler.GetOne)

	return logRequest(mux)
}
