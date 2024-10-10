package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	userRepositoryPg := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})

	authService := authservices.NewAuthService(authservices.AuthServiceOptions{
		UserService:        userService,
		Cache:              appCtx.RedisClient,
		EnvVariableService: appCtx.EnvVariableService,
	})
	jwtAuthMiddleware := middlewares.JwtAuth(authService)
	strictSessionMiddleware := middlewares.StrictSession(authService)

	authHandler := NewAuthHandler(AuthHandlerOptions{
		AuthService: authService,
		Logger:      appCtx.Logger,
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(jwtAuthMiddleware)
		r.Use(strictSessionMiddleware)
		
		r.Post("/logout", authHandler.Logout)
	})

	return r
}
