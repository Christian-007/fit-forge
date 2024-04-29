package domains

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AppContext struct {
	Logger *slog.Logger
	Db *pgxpool.Pool
}
