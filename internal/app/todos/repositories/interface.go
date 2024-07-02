package repositories

import "github.com/Christian-007/fit-forge/internal/app/todos/domains"

type TodoRepository interface {
	GetAll() ([]domains.TodoModel, error)
}
