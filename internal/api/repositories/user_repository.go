package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/api/apperrors"
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

	defer rows.Close()

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

func (u UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = $1"

	cmdTag, err := u.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrUserNotFound
	}

	return nil
}

func (u UserRepository) UpdateOne(id int, updateUser domains.UserModel) (domains.UserModel, error) {
	args := createUpdateUserPgxArgs(id, updateUser)
	query := `
		UPDATE users 
		SET name = COALESCE(@name, name), email = COALESCE(@email, email), password = COALESCE(@password, password) 
		WHERE id = @id
		RETURNING id, name, email, password, created_at
	`

	var user domains.UserModel
	err := u.db.QueryRow(context.Background(), query, args).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return domains.UserModel{}, err
	}

	return user, nil
}

func createUpdateUserPgxArgs(id int, updateUser domains.UserModel) pgx.NamedArgs {
	result := pgx.NamedArgs{
		"id": id,
	}

	if updateUser.Name != "" {
		result["name"] = updateUser.Name
	}
	if updateUser.Email != "" {
		result["email"] = updateUser.Email
	}
	if len(updateUser.Password) > 0 {
		result["password"] = updateUser.Password
	}

	return result
}
