package domains

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
)

type TodoModel struct {
	Id          int
	Title       string
	IsCompleted bool
	UserId      int
	CreatedAt   time.Time
}

type TodoWithPoints struct {
	Id          int
	Title       string
	IsCompleted bool
	UserId      int
	CreatedAt   time.Time
	Point       domains.PointModel `json:"point"`
}
