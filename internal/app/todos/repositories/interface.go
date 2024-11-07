package repositories

import "github.com/Christian-007/fit-forge/internal/app/todos/domains"

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type TodoRepository interface {
	GetAll() ([]domains.TodoModel, error)
	GetAllByUserId(userId int) ([]domains.TodoModel, error)
	GetOneByUserId(userId int, todoId int) (domains.TodoModel, error)
	Create(userId int, todo domains.TodoModel) (domains.TodoModel, error)
	Delete(todoId int, userId int) error
}
