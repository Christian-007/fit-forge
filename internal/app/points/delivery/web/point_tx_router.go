package web

import (
	authservices "github.com/Christian-007/fit-forge/internal/app/auth/services"
	"github.com/Christian-007/fit-forge/internal/app/points/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/appcontext"
	"github.com/Christian-007/fit-forge/internal/pkg/middlewares"
	"github.com/go-chi/chi/v5"
)

func PointTxRoutes(appCtx appcontext.AppContext) *chi.Mux {
	r := chi.NewRouter()

	pointTransactionsRepository := repositories.NewPointTransactionsRepositoryPg(appCtx.Pool)
	pointTransactionsHandler := NewPointTransactionsHandler(PointTransactionsOptions{
		Logger: appCtx.Logger,
		PointTransactionsRepository: pointTransactionsRepository,
	})

	authService := authservices.NewAuthServiceImpl(authservices.AuthServiceOptions{
		Cache: appCtx.RedisClient,
	})
	strictSessionMiddleware := middlewares.StrictSession(authService, appCtx.SecretManagerClient)

	r.Use(strictSessionMiddleware)
	r.Get("/", pointTransactionsHandler.GetAllWithPagination)

	return r
}
