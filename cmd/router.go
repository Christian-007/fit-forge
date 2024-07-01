package main

import (
	"net/http"

	todosweb "github.com/Christian-007/fit-forge/internal/app/todos/delivery/web"
	usersweb "github.com/Christian-007/fit-forge/internal/app/users/delivery/web"
	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
)

func Routes(appCtx domains.AppContext) http.Handler {
	mux := http.NewServeMux()

	logRequest := middlewares.NewLogRequest(appCtx.Logger)

	usersweb.Routes(mux, appCtx)
	todosweb.Routes(mux)

	return logRequest(mux)
}
