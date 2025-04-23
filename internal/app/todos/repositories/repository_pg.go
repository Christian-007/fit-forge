package repositories

import (
	"context"

	pointdomains "github.com/Christian-007/fit-forge/internal/app/points/domains"
	"github.com/Christian-007/fit-forge/internal/app/todos/domains"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/google/uuid"
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

func (t TodoRepositoryPg) GetAllByUserId(userId int) ([]domains.TodoModel, error) {
	query := "SELECT * FROM todos WHERE user_id = $1"

	rows, _ := t.db.Query(context.Background(), query, userId)
	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.TodoModel])
	if err != nil {
		return []domains.TodoModel{}, err
	}

	defer rows.Close()

	return todos, nil
}

func (t TodoRepositoryPg) GetOneByUserId(userId int, todoId int) (domains.TodoModel, error) {
	query := "SELECT * from todos WHERE id = $1 AND user_id = $2"
	row, err := t.db.Query(context.Background(), query, todoId, userId)
	if err != nil {
		return domains.TodoModel{}, err
	}

	defer row.Close()

	todo, err := pgx.CollectOneRow(row, pgx.RowToStructByName[domains.TodoModel])
	if err != nil {
		return domains.TodoModel{}, err
	}

	return todo, nil
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

func (t TodoRepositoryPg) CreateWithPoints(ctx context.Context, userId int, todo domains.TodoModel) (domains.TodoWithPoints, error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return domains.TodoWithPoints{}, err
	}
	defer tx.Rollback(ctx)

	var row pgx.Row

	// Step 1: Create todo
	var insertedTodo domains.TodoWithPoints
	query := "INSERT INTO todos(user_id, title) VALUES ($1, $2) RETURNING *"
	row = tx.QueryRow(ctx, query, userId, todo.Title)
	err = row.Scan(
		&insertedTodo.Id,
		&insertedTodo.Title,
		&insertedTodo.IsCompleted,
		&insertedTodo.UserId,
		&insertedTodo.CreatedAt,
	)
	if err != nil {
		return domains.TodoWithPoints{}, err
	}

	// Step 2: Insert points
	var insertedPoint pointdomains.PointModel
	earnedPoints := 2
	query = "INSERT INTO points(user_id, total_points) VALUES ($1, $2) RETURNING user_id, total_points, created_at, updated_at"
	row = tx.QueryRow(ctx, query, insertedTodo.UserId, earnedPoints)
	err = row.Scan(
		&insertedPoint.UserId,
		&insertedPoint.TotalPoints,
		&insertedPoint.CreatedAt,
		&insertedPoint.UpdatedAt,
	)
	if err != nil {
		return domains.TodoWithPoints{}, err
	}

	// Step 3: Log to the point transaction
	pointTransactions := pointdomains.PointTransactionsModel{
		ID:              uuid.New(),
		TransactionType: pointdomains.EarnTransactionType,
		Points:          earnedPoints,
		Reason:          pointdomains.CreateTodoReason,
		UserID:          insertedTodo.UserId,
	}
	query = "INSERT INTO point_transactions(id, transaction_type, points, reason, user_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, pointTransactions.ID, pointTransactions.TransactionType, pointTransactions.Points, pointTransactions.Reason, pointTransactions.UserID)
	if err != nil {
		return domains.TodoWithPoints{}, err
	}

	// Step 4: Add the inserted points to the user model
	insertedTodo.Point = insertedPoint

	err = tx.Commit(ctx)
	if err != nil {
		return domains.TodoWithPoints{}, err
	}

	return insertedTodo, nil
}

func (t TodoRepositoryPg) Delete(todoId int, userId int) error {
	query := "DELETE FROM todos WHERE id = $1 AND user_id = $2"

	cmdTag, err := t.db.Exec(context.Background(), query, todoId, userId)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrUserOrTodoNotFound
	}

	return nil
}
