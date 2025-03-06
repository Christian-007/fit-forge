package repositories

import "github.com/Christian-007/fit-forge/internal/app/points/domains"

type PointsRepository interface {
	GetOne(userId int) (domains.PointModel, error)
	Create(point domains.PointModel) (domains.PointModel, error)
}
