package dto

import validation "github.com/go-ozzo/ozzo-validation/v4"

type TodoResponse struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"isCompleted"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (c CreateTodoRequest) Validate() error {
	return validation.ValidateStruct(&c, validation.Field(&c.Title, validation.Required))
}
