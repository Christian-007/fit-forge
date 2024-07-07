package dto

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TodoResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

type GetAllTodosResponse struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"isCompleted"`
	UserId      int       `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (c CreateTodoRequest) Validate() error {
	return validation.ValidateStruct(&c, validation.Field(&c.Title, validation.Required))
}
