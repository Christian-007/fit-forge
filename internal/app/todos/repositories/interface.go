package repositories

import "github.com/Christian-007/fit-forge/internal/app/todos/domains"

type TodoRepository interface {
	GetAll() ([]domains.TodoModel, error)
	GetOne(userId int, todoId int) (domains.TodoModel, error)
	Create(userId int, todo domains.TodoModel) (domains.TodoModel, error)
	Delete(todoId int, userId int) error
}
