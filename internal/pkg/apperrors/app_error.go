package apperrors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTodoNotFound       = errors.New("todo not found")
	ErrUserOrTodoNotFound = errors.New("user or todo not found")
)
