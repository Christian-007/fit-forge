package web

import (
	authservice "github.com/Christian-007/fit-forge/internal/app/auth/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	userRepositoryPg := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})

	authService := authservice.NewAuthService(authservice.AuthServiceOptions{
		UserService: userService,
	})

	authHandler := NewAuthHandler(AuthHandlerOptions{
		AuthService: authService,
		Logger:      appCtx.Logger,
	})

	r.Post("/login", authHandler.Login)

	return r
}
