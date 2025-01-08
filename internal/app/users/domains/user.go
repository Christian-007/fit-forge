package domains

import "time"

type UserModel struct {
	Id              int        `json:"id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Password        []byte     `json:"password"`
	Role            int        `json:"role"` // 1 is admin and, 2 is user
	CreatedAt       time.Time  `json:"createdAt"`
	EmailVerifiedAt *time.Time `json:"emailVerifiedAt"`
}
