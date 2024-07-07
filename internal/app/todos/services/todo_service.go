package services

import (
	"errors"

	"github.com/Christian-007/fit-forge/internal/app/todos/domains"
	"github.com/Christian-007/fit-forge/internal/app/todos/dto"
	"github.com/Christian-007/fit-forge/internal/app/todos/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
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

func (t TodoService) GetAll() ([]dto.GetAllTodosResponse, error) {
	todos, err := t.TodoRepository.GetAll()
	if err != nil {
		return []dto.GetAllTodosResponse{}, err
	}

	getAllTodosResponse := make([]dto.GetAllTodosResponse, len(todos))
	for i, todo := range todos {
		getAllTodosResponse[i] = toGetAllTodosResponse(todo)
	}

	return getAllTodosResponse, nil
}

func (t TodoService) GetAllByUserId(userId int) ([]dto.TodoResponse, error) {
	todos, err := t.TodoRepository.GetAllByUserId(userId)
	if err != nil {
		return []dto.TodoResponse{}, err
	}

	todoResponse := make([]dto.TodoResponse, len(todos))
	for i, todo := range todos {
		todoResponse[i] = toTodoResponse(todo)
	}

	return todoResponse, nil
}

func (t TodoService) GetOneByUserId(userId int, todoId int) (dto.TodoResponse, error) {
	todoModel, err := t.TodoRepository.GetOneByUserId(userId, todoId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.TodoResponse{}, apperrors.ErrTodoNotFound
		}
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

func toGetAllTodosResponse(todoModel domains.TodoModel) dto.GetAllTodosResponse {
	return dto.GetAllTodosResponse{
		Id:          todoModel.Id,
		Title:       todoModel.Title,
		IsCompleted: todoModel.IsCompleted,
		UserId:      todoModel.UserId,
		CreatedAt:   todoModel.CreatedAt,
	}
}
