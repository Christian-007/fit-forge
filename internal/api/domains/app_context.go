package domains

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	AppContextOptions
}

type AppContextOptions struct {
	Logger *slog.Logger
	Pool   *pgxpool.Pool
}

func NewAppContext(options AppContextOptions) AppContext {
	return AppContext{
		options,
	}
}
