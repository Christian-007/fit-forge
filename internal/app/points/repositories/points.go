package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/jackc/pgx/v5"
)

type PointsRepository interface {
	GetOne(tx pgx.Tx, userId int) (domains.PointModel, error)
	Create(tx pgx.Tx, point domains.PointModel) (domains.PointModel, error)
	UpdateWithTransactionHistory(ctx context.Context, userId int, pointTransaction domains.CreatePointTransactions) (domains.PointModel, error)
	// "dueDate" format is "yyyy-mm-dd"
	FindUsersForSubscriptionDeduction(ctx context.Context, dueDate string) (domains.UsersDueForSubscription, error)
}
