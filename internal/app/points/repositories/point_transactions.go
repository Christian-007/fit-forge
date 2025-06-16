package repositories

import (
	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/jackc/pgx/v5"
)

type PointTransactionsRepostiory interface {
	Create(tx pgx.Tx, transaction domains.CreatePointTransactions) error
}
