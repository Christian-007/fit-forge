package services

import (
	"github.com/Christian-007/fit-forge/internal/app/todos/domains"
	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	"github.com/Christian-007/fit-forge/internal/app/todos/repositories"
)

type TodoService struct {
	TodoServiceOptions
}

type TodoServiceOptions struct {
	TodoRepository repositories.TodoRepository
}

func NewTodoService(options TodoServiceOptions) TodoService {
	return TodoService{
		options,
	}
}

func (t TodoService) GetAll() ([]dto.TodoResponse, error) {
	todos, err := t.TodoRepository.GetAll()
	if err != nil {
		return []dto.TodoResponse{}, err
	}

	todoResponse := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponse[i] = toTodoResponse(todo)
	}

	return todoResponse, nil
}

func (t TodoService) GetOne(userId int, todoId int) (dto.TodoResponse, error) {
	todoModel, err := t.TodoRepository.GetOne(userId, todoId)
	if err != nil {
		return dto.TodoResponse{}, err
	}

	return toTodoResponse(todoModel), nil
}

func (t TodoService) Create(userId int, createTodoRequest dto.CreateTodoRequest) (dto.TodoResponse, error) {
	todoModel := domains.TodoModel{
		Title: createTodoRequest.Title,
	}

	createdTodoModel, err := t.TodoRepository.Create(userId, todoModel)
	if err != nil {
		return dto.TodoResponse{}, err
	}

	return toTodoResponse(createdTodoModel), nil
}

func (t TodoService) Delete(todoId int, userId int) error {
	err := t.TodoRepository.Delete(todoId, userId)
	if err != nil {
		return err
	}

	return nil
}

func toTodoResponse(todoModel domains.TodoModel) dto.TodoResponse {
	return dto.TodoResponse{
		Id:          todoModel.Id,
		Title:       todoModel.Title,
		IsCompleted: todoModel.IsCompleted,
	}
}
