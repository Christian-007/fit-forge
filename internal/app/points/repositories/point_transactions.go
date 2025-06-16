package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/jackc/pgx/v5"
)

type PointTransactionsRepostiory interface {
	Create(tx pgx.Tx, transaction domains.CreatePointTransactions) error
	GetAllWithPagination(ctx context.Context, userId int, limit int, offset int) ([]domains.PointTransaction, int, error)
}
