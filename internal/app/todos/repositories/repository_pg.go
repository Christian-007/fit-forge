package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/todos/domains"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepositoryPg struct {
	db *pgxpool.Pool
}

func NewTodoRepositoryPg(pool *pgxpool.Pool) TodoRepositoryPg {
	return TodoRepositoryPg{
		db: pool,
	}
}

func (t TodoRepositoryPg) GetAll() ([]domains.TodoModel, error) {
	query := "SELECT * from todos"
	rows, _ := t.db.Query(context.Background(), query)
	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.TodoModel])
	if err != nil {
		return []domains.TodoModel{}, err
	}

	defer rows.Close()

	return todos, nil
}

func (t TodoRepositoryPg) Create(userId int, todo domains.TodoModel) (domains.TodoModel, error) {
	query := "INSERT INTO todos(user_id, title) VALUES ($1, $2) RETURNING *"

	var insertedTodo domains.TodoModel
	err := t.db.QueryRow(context.Background(), query, userId, todo.Title).Scan(
		&insertedTodo.Id,
		&insertedTodo.Title,
		&insertedTodo.IsCompleted,
		&insertedTodo.UserId,
		&insertedTodo.CreatedAt,
	)
	if err != nil {
		return domains.TodoModel{}, err
	}

	return insertedTodo, nil
}
