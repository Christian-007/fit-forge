package model

import "time"

type UserWithPoints struct {
	Id              int                `json:"id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	Password        []byte             `json:"password"`
	Role            int                `json:"role"` // 1 is admin and, 2 is user
	CreatedAt       time.Time          `json:"createdAt"`
	EmailVerifiedAt *time.Time         `json:"emailVerifiedAt"`
	Point           PointModel `json:"point"`
}

type PointModel struct {
	UserId      int       `json:"userId"`
	TotalPoints int       `json:"totalPoints"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
