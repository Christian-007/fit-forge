package appcontext

import (
	"github.com/Christian-007/fit-forge/internal/pkg/applog"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/envvariable"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	AppContextOptions
}

type AppContextOptions struct {
	Logger             applog.Logger
	Pool               *pgxpool.Pool
	RedisClient        *cache.RedisCache
	EnvVariableService envvariable.EnvVariableService
	Publisher          message.Publisher
}

func NewAppContext(options AppContextOptions) AppContext {
	return AppContext{
		options,
	}
}
