package repositories

import (
	"context"
	"fmt"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/pkg/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PointsRepositoryPg struct {
	db *pgxpool.Pool
}

func NewPointsRepositoryPg(pool *pgxpool.Pool) PointsRepositoryPg {
	return PointsRepositoryPg{
		db: pool,
	}
}

func (p PointsRepositoryPg) GetOne(tx pgx.Tx, userId int) (domains.PointModel, error) {
	ctx := context.Background()
	query := "SELECT user_id, total_points, created_at, updated_at FROM points WHERE user_id=$1 FOR UPDATE"

	var result domains.PointModel
	var row pgx.Row

	if tx != nil {
		row = tx.QueryRow(ctx, query, userId)
	} else {
		row = p.db.QueryRow(ctx, query, userId)
	}

	err := row.Scan(
		&result.UserId,
		&result.TotalPoints,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return domains.PointModel{}, err
	}

	return result, nil
}

func (p PointsRepositoryPg) Create(tx pgx.Tx, point domains.PointModel) (domains.PointModel, error) {
	ctx := context.Background()
	query := "INSERT INTO points(user_id, total_points) VALUES ($1, $2) RETURNING user_id, total_points, created_at, updated_at"

	var insertedPoint domains.PointModel
	var row pgx.Row

	if tx != nil {
		row = tx.QueryRow(ctx, query, point.UserId, point.TotalPoints)
	} else {
		row = p.db.QueryRow(ctx, query, point.UserId, point.TotalPoints)
	}

	err := row.Scan(
		&insertedPoint.UserId,
		&insertedPoint.TotalPoints,
		&insertedPoint.CreatedAt,
		&insertedPoint.UpdatedAt,
	)
	if err != nil {
		return domains.PointModel{}, err
	}

	return insertedPoint, nil
}

func (p PointsRepositoryPg) UpdateWithTransactionHistory(ctx context.Context, userId int, pointTransaction domains.PointTransactionsModel) (domains.PointModel, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return domains.PointModel{}, err
	}
	defer tx.Rollback(ctx)

	var row pgx.Row

	// Step 1: Add points by updating the existing point
	var updatedPoint domains.PointModel
	query := `
		UPDATE points 
		SET
			total_points = total_points + $1
		WHERE user_id = $2
		RETURNING user_id, total_points
	`

	row = tx.QueryRow(ctx, query, pointTransaction.Points, userId)
	err = row.Scan(&updatedPoint.UserId, &updatedPoint.TotalPoints)
	if err != nil {
		return domains.PointModel{}, err
	}

	// Step 2: Log to the point transaction
	pointTransaction.ID = uuid.New()
	pointTransaction.UserID = userId
	query = "INSERT INTO point_transactions(id, transaction_type, points, reason, user_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, pointTransaction.ID, pointTransaction.TransactionType, pointTransaction.Points, pointTransaction.Reason, pointTransaction.UserID)
	if err != nil {
		return domains.PointModel{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return domains.PointModel{}, err
	}

	return updatedPoint, nil
}

func (p PointsRepositoryPg) FindUsersForSubscriptionDeduction(ctx context.Context, dueDate string) (domains.UsersDueForSubscription, error) {
	dueDateQuery := fmt.Sprintf("'%s'", dueDate)
	query := `
	SELECT
		u.id,
		u.email,
		p.total_points
	FROM
		users u
		INNER JOIN points p ON u.id = p.user_id
	WHERE
		COALESCE(
			(
				SELECT
					MAX(pt.created_at)
				FROM
					point_transactions pt
				WHERE
					pt.transaction_type = 'SUBSCRIPTION_DEDUCTION'
			),
			u.created_at::date
		) + INTERVAL '1 month' = ` + dueDateQuery + ";"

	usersDueForSubscription := domains.UsersDueForSubscription{}

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return usersDueForSubscription, err
	}

	defer rows.Close()

	users := []model.UserWithPoints{}
	for rows.Next() {
		user := model.UserWithPoints{}
		err := rows.Scan(&user.Id, &user.Email, &user.TotalPoints)
		if err != nil {
			return usersDueForSubscription, fmt.Errorf("unable to scan row: %w", err)
		}
		users = append(users, user)
	}

	for _, user := range users {
		if user.TotalPoints >= domains.SubscriptionDeductionAmount {
			usersDueForSubscription.EligibleForDeduction = append(usersDueForSubscription.EligibleForDeduction, user)
		} else {
			usersDueForSubscription.InsufficientPoints = append(usersDueForSubscription.InsufficientPoints, user)
		}
	}

	return usersDueForSubscription, nil
}
