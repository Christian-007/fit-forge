package domains

import (
	"time"

	sharedmodel "github.com/Christian-007/fit-forge/internal/pkg/model"
)

type PointModel struct {
	UserId      int       `json:"userId"`
	TotalPoints int       `json:"totalPoints"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PointChange struct {
	Total  int    `json:"total"`
	Change string `json:"change"`
}

type UsersDueForSubscription struct {
	EligibleForDeduction []sharedmodel.UserWithPoints
	InsufficientPoints   []sharedmodel.UserWithPoints
}

const (
	SubscriptionDeductionAmount int = -25
)
