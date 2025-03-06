package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
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

func (p PointsRepositoryPg) GetOne(userId int) (domains.PointModel, error) {
	query := "SELECT user_id, total_points, created_at, updated_at FROM points WHERE user_id=$1"
	rows, err := p.db.Query(context.Background(), query, userId)
	if err != nil {
		return domains.PointModel{}, err
	}

	defer rows.Close()

	result, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domains.PointModel])
	if err != nil {
		return domains.PointModel{}, err
	}

	return result, nil
}

func (p PointsRepositoryPg) Create(point domains.PointModel) (domains.PointModel, error) {
	query := "INSERT INTO points(user_id, total_points) VALUES ($1, $2) RETURNING user_id, total_points, created_at, updated_at"

	var insertedPoint domains.PointModel
	err := p.db.QueryRow(context.Background(), query, point.UserId, point.TotalPoints).Scan(
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
