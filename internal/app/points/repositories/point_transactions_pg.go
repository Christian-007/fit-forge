package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PointTransactionsRepositoryPg struct {
	db *pgxpool.Pool
}

func NewPointTransactionsRepositoryPg(pool *pgxpool.Pool) PointTransactionsRepositoryPg {
	return PointTransactionsRepositoryPg{db: pool}
}

func (p PointTransactionsRepositoryPg) Create(tx pgx.Tx, transaction domains.CreatePointTransactions) error {
	query := "INSERT INTO point_transactions(id, transaction_type, points, reason, user_id) VALUES ($1, $2, $3, $4, $5), RETURNING id, transaction_type, points, reason, user_id, created_at"
	_, err := tx.Exec(context.Background(), query)

	return err
}
