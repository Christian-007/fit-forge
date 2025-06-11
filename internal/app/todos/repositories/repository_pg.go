package repositories

import (
	"context"
	"fmt"
	"strings"

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
	query := "SELECT * FROM todos WHERE user_id = $1 ORDER BY created_at DESC"

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

func (t TodoRepositoryPg) CreateWithPoints(ctx context.Context, todo domains.TodoModel, userId int) (domains.TodoModelWithPoints, error) {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return domains.TodoModelWithPoints{}, err
	}
	defer tx.Rollback(ctx)

	var row pgx.Row

	// Step 1: Create todo
	var insertedTodo domains.TodoModelWithPoints
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
		return domains.TodoModelWithPoints{}, err
	}

	// Step 2: Add points by updating the existing point
	var insertedPoint pointdomains.PointModel
	addedPoints := 2
	query = `
		UPDATE points 
		SET
			total_points = total_points + $1
		WHERE user_id = $2
		RETURNING user_id, total_points
	`
	row = tx.QueryRow(ctx, query, addedPoints, userId)
	err = row.Scan(&insertedPoint.UserId, &insertedPoint.TotalPoints)
	if err != nil {
		return domains.TodoModelWithPoints{}, err
	}

	// Step 3: Log to the point transaction
	pointTransactions := pointdomains.PointTransactionsModel{
		ID:              uuid.New(),
		TransactionType: pointdomains.EarnTransactionType,
		Points:          addedPoints,
		Reason:          "created a todo",
		UserID:          userId,
	}
	query = "INSERT INTO point_transactions(id, transaction_type, points, reason, user_id) VALUES ($1, $2, $3, $4, $5)"
	_, err = tx.Exec(ctx, query, pointTransactions.ID, pointTransactions.TransactionType, pointTransactions.Points, pointTransactions.Reason, pointTransactions.UserID)
	if err != nil {
		return domains.TodoModelWithPoints{}, err
	}

	// Step 4: Add the inserted points to the todo model
	insertedTodo.Points = pointdomains.PointChange{
		Total:  insertedPoint.TotalPoints,
		Change: fmt.Sprintf("+%d", addedPoints),
	}
	err = tx.Commit(ctx)
	if err != nil {
		return domains.TodoModelWithPoints{}, err
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

func (t TodoRepositoryPg) Update(ctx context.Context, todoId int, updates map[string]any) error {
	var (
		setClauses []string
		args       []any
		i          = 1
	)

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, i))
		args = append(args, value)
		i++
	}

	query := fmt.Sprintf("UPDATE todos SET %s WHERE id = %d", strings.Join(setClauses, ", "), todoId)
	_, err := t.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
