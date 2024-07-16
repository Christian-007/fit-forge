package main

import (
	"net/http"

	authweb "github.com/Christian-007/fit-forge/internal/app/auth/delivery/web"
	todosweb "github.com/Christian-007/fit-forge/internal/app/todos/delivery/web"
	usersweb "github.com/Christian-007/fit-forge/internal/app/users/delivery/web"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/justinas/alice"
)

func Routes(appCtx appcontext.AppContext) http.Handler {
	mux := http.NewServeMux()

	logRequest := middlewares.NewLogRequest(appCtx.Logger)
	
	standard := alice.New(logRequest)

	usersweb.Routes(mux, appCtx)
	todosweb.Routes(mux, appCtx)
	authweb.Routes(mux, appCtx)

	return standard.Then(mux)
}
