package dto

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type TodoResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

type TodoWithPointsResponse struct {
	Id          int                 `json:"id"`
	Title       string              `json:"title"`
	IsCompleted bool                `json:"isCompleted"`
	Points      domains.PointChange `json:"points"`
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

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	IsCompleted *bool   `json:"isCompleted"`
}

func (u UpdateTodoRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Title,
			validation.NilOrNotEmpty,
			validation.When(u.Title != nil, validation.Length(1, 200)),
		),
		validation.Field(&u.IsCompleted),
	)
}

func (u UpdateTodoRequest) IsCompletedTrue() bool {
	return u.IsCompleted != nil && *u.IsCompleted
}
