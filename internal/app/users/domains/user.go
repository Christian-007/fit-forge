package domains

import "time"

type UserModel struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}
