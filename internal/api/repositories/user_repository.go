package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db    *pgxpool.Pool
	users []domains.UserModel
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return UserRepository{
		users: []domains.UserModel{},
		db:    pool,
	}
}

func (u UserRepository) GetAll() ([]domains.UserModel, error) {
	rows, _ := u.db.Query(context.Background(), "SELECT * FROM users")
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.UserModel])
	if err != nil {
		return []domains.UserModel{}, err
	}

	return users, nil
}

func (u UserRepository) GetOne(id int) (domains.UserModel, error) {
	query := "SELECT * FROM users WHERE id=$1"
	rows, err := u.db.Query(context.Background(), query, id)
	if err != nil {
		return domains.UserModel{}, err
	}

	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domains.UserModel])
	if err != nil {
		return domains.UserModel{}, err
	}

	return user, nil
}

func (u UserRepository) Create(user domains.UserModel) (domains.UserModel, error) {
	query := "INSERT INTO users(name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, password, created_at"

	var insertedUser domains.UserModel
	err := u.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Password).Scan(
		&insertedUser.Id,
		&insertedUser.Name,
		&insertedUser.Email,
		&insertedUser.Password,
		&insertedUser.CreatedAt,
	)
	if err != nil {
		return domains.UserModel{}, err
	}

	return insertedUser, nil
}
