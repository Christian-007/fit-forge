package domains

import "time"

type TodoModel struct {
	Id          int
	Title       string
	IsCompleted bool
	UserId      int
	CreatedAt   time.Time
}
