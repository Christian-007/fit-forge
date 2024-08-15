package appcontext

import (
	"log/slog"

	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	AppContextOptions
}

type AppContextOptions struct {
	Logger      *slog.Logger
	Pool        *pgxpool.Pool
	RedisClient *cache.RedisCache
}

func NewAppContext(options AppContextOptions) AppContext {
	return AppContext{
		options,
	}
}
