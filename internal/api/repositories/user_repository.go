package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db    *pgxpool.Pool
	users []domains.User
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return UserRepository{
		users: []domains.User{},
		db:    pool,
	}
}

func (u UserRepository) GetAll() ([]domains.User, error) {
	rows, _ := u.db.Query(context.Background(), "SELECT * FROM users")
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.User])
	if err != nil {
		return []domains.User{}, err
	}

	return users, nil
}

func (u UserRepository) GetOne(id int) (domains.User, error) {
	query := "SELECT * FROM users WHERE id=$1"
	rows, err := u.db.Query(context.Background(), query, id)
	if err != nil {
		return domains.User{}, err
	}
	
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domains.User])
	if err != nil {
		return domains.User{}, err
	}

	return user, nil
}
