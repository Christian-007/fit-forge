package apperrors

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrTodoNotFound       = errors.New("todo not found")
	ErrUserOrTodoNotFound = errors.New("user or todo not found")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidSignature   = errors.New("invalid token signature")
	ErrInvalidToken       = errors.New("invalid token")
	ErrEmailDuplicate     = errors.New("email is duplicated")

	ErrRedisKeyNotFound = errors.New("key does not exist in Redis")
	ErrTypeAssertion    = errors.New("type assertion failed")
)
