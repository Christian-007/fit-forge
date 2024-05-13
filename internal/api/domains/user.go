package domains

import "time"

type User struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  []byte    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}
