package main

import (
	authweb "github.com/Christian-007/fit-forge/internal/app/auth/delivery/web"
	todosweb "github.com/Christian-007/fit-forge/internal/app/todos/delivery/web"
	usersweb "github.com/Christian-007/fit-forge/internal/app/users/delivery/web"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()

	logRequest := middlewares.NewLogRequest(appCtx.Logger)

	r.Use(logRequest)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	r.Mount("/users", usersweb.Routes(appCtx))
	r.Mount("/todos", todosweb.Routes(appCtx))
	r.Mount("/auth", authweb.Routes(appCtx))

	return r
}
