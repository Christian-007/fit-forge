package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(2, 200)),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required),
	)
}
