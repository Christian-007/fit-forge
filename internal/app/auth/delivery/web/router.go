package web

import (
	"net/http"

	authservice "github.com/Christian-007/fit-forge/internal/app/auth/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
)

func Routes(mux *http.ServeMux, appCtx appcontext.AppContext) {
	userRepositoryPg := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})

	authService := authservice.NewAuthService(authservice.AuthServiceOptions{
		UserService: userService,
	})

	authHandler := NewAuthHandler(AuthHandlerOptions{
		AuthService: authService,
		Logger: appCtx.Logger,
	})

	mux.HandleFunc("POST /auth/login", authHandler.Login)
}