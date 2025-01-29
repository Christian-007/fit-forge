package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	emailservices "github.com/Christian-007/fit-forge/internal/app/email/services"
	"github.com/Christian-007/fit-forge/internal/app/users/repositories"
	userservices "github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
	"github.com/go-chi/chi/v5"
)

func Routes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()
	userRepositoryPg := repositories.NewUserRepositoryPg(appCtx.Pool)
	userService := userservices.NewUserService(userservices.UserServiceOptions{
		UserRepository: userRepositoryPg,
	})
	tokenService := security.NewTokenService(security.TokenServiceOptions{
		SecretKey: appCtx.EnvVariableService.Get("AUTH_SECRET_KEY"),
	})
	emailService := emailservices.NewEmailService(emailservices.EmailServiceOptions{
		Host:         "http://localhost:4000",
		Cache:        appCtx.RedisClient,
		TokenService: tokenService,
	})
	mailtrapSender := emailservices.NewMailtrapEmailService(emailservices.MailtrapSenderOptions{
		Host:   appCtx.EnvVariableService.Get("EMAIL_HOST"),
		ApiKey: appCtx.EnvVariableService.Get("MAILTRAP_API_KEY"),
	})
	userHandler := NewUserHandler(UserHandlerOptions{
		UserService:    userService,
		Logger:         appCtx.Logger,
		EmailService:   emailService,
		MailtrapSender: mailtrapSender,
		Publisher:      appCtx.Publisher,
	})
	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		UserService:        userService,
		Cache:              appCtx.RedisClient,
		EnvVariableService: appCtx.EnvVariableService,
	})

	strictSessionMiddleware := middlewares.StrictSession(authService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/", userHandler.Create)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(strictSessionMiddleware)

		r.Get("/", userHandler.GetAll)
		r.Get("/{id}", userHandler.GetOne)
		r.Delete("/{id}", userHandler.Delete)
		r.Patch("/{id}", userHandler.UpdateOne)
	})

	return r
}
