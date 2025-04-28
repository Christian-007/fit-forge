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

type TodoModelWithPoints struct {
	Id          int                 `json:"id"`
	Title       string              `json:"title"`
	IsCompleted bool                `json:"isCompleted"`
	UserId      int                 `json:"userId"`
	CreatedAt   time.Time           `json:"createdAt"`
	Points      domains.PointChange `json:"points"`
}
