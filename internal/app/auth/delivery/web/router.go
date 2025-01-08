package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
	userrepositories "github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	userRepositoryPg := userrepositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})

	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		UserService:        userService,
		Cache:              appCtx.RedisClient,
		EnvVariableService: appCtx.EnvVariableService,
	})
	logoutSessionMiddleware := middlewares.LogoutSession(authService)

	// TODO: move secret key to .env
	tokenService := security.NewTokenService(security.TokenServiceOptions{
		SecretKey: "haha",
	})
	emailService := emailservices.NewEmailService(emailservices.EmailServiceOptions{
		Host:         "http://localhost:4000",
		Cache:        appCtx.RedisClient,
		TokenService: tokenService,
	})

	authHandler := NewAuthHandler(AuthHandlerOptions{
		AuthService:  authService,
		Logger:       appCtx.Logger,
		UserService:  userService,
		EmailService: emailService,
		Cache: appCtx.RedisClient,
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/verify/{token}", authHandler.Verify)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(logoutSessionMiddleware)

		r.Post("/logout", authHandler.Logout)
	})

	return r
}
