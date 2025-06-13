package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()

	profileHandler := NewProfileHandler(ProfileHandlerOptions{
		Logger: appCtx.Logger,
	})

	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		Cache: appCtx.RedisClient,
	})
	strictSessionMiddleware := middlewares.StrictSession(authService, appCtx.SecretManagerClient)

	r.Use(strictSessionMiddleware)
	r.Group(func(r chi.Router) {
		r.Get("/", profileHandler.Get)
	})

	return r
}
