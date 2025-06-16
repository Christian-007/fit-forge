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

func (p PointTransactionsRepositoryPg) GetAllWithPagination(ctx context.Context, userId int, limit int, offset int) ([]domains.PointTransaction, int, error) {
	query := `
		SELECT id, transaction_type, points, reason, created_at
		FROM point_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, _ := p.db.Query(ctx, query, userId, limit, offset)
	pointTransactions, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.PointTransaction])
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	countQuery := `SELECT COUNT(*) FROM point_transactions WHERE user_id = $1`
	var total int
	err = p.db.QueryRow(ctx, countQuery, userId).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return pointTransactions, total, nil
}

