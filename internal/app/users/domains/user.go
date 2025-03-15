package domains

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
)

type UserModel struct {
	Id              int        `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Password        []byte     `json:"password"`
	Role            int        `json:"role"` // 1 is admin and, 2 is user
	CreatedAt       time.Time  `json:"createdAt"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt"`
}

type UserWithPoints struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Password        []byte             `json:"password"`
	Role            int                `json:"role"` // 1 is admin and, 2 is user
	CreatedAt       time.Time          `json:"createdAt"`
	EmailVerifiedAt *time.Time         `json:"emailVerifiedAt"`
	Point           domains.PointModel `json:"point"`
}
