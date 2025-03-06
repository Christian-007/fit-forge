package domains

import "time"

type PointModel struct {
	UserId      int       `json:"userId"`
	TotalPoints int       `json:"totalPoints"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
