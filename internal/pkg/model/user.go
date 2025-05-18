package model

type UserWithPoints struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	TotalPoints int    `json:"totalPoints"`
}
