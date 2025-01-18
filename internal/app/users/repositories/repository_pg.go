package repositories

import (
	"context"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPg struct {
	db    *pgxpool.Pool
	users []domains.UserModel
}

func NewUserRepositoryPg(pool *pgxpool.Pool) UserRepositoryPg {
	return UserRepositoryPg{
		users: []domains.UserModel{},
		db:    pool,
	}
}

func (u UserRepositoryPg) GetAll() ([]domains.UserModel, error) {
	rows, _ := u.db.Query(context.Background(), "SELECT * FROM users")
	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domains.UserModel])
	if err != nil {
		return []domains.UserModel{}, err
	}

	defer rows.Close()

	return users, nil
}

func (u UserRepositoryPg) GetOne(id int) (domains.UserModel, error) {
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

func (u UserRepositoryPg) GetOneByEmail(email string) (domains.UserModel, error) {
	query := "SELECT * FROM users WHERE email=$1"
	rows, err := u.db.Query(context.Background(), query, email)
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

func (u UserRepositoryPg) Create(user domains.UserModel) (domains.UserModel, error) {
	query := "INSERT INTO users(name, email, password) VALUES ($1, $2, $3) RETURNING id, name, email, password, role, email_verified_at, created_at"

	var insertedUser domains.UserModel
	err := u.db.QueryRow(context.Background(), query, user.Name, user.Email, user.Password).Scan(
		&insertedUser.Id,
		&insertedUser.Name,
		&insertedUser.Email,
		&insertedUser.Password,
		&insertedUser.Role,
		&insertedUser.EmailVerifiedAt,
		&insertedUser.CreatedAt,
	)
	if err != nil {
		return domains.UserModel{}, err
	}

	return insertedUser, nil
}

func (u UserRepositoryPg) Delete(id int) error {
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

func (u UserRepositoryPg) UpdateOne(id int, updateUser domains.UserModel) (domains.UserModel, error) {
	args := createUpdateUserPgxArgs(id, updateUser)
	query := `
		UPDATE users 
		SET
			name = COALESCE(@name, name),
			email = COALESCE(@email, email),
			password = COALESCE(@password, password),
			role = COALESCE(@role, role),
			email_verified_at = COALESCE(@email_verified_at, email_verified_at)
		WHERE id = @id
		RETURNING id, name, email, password, role, created_at, email_verified_at
	`

	var user domains.UserModel
	err := u.db.QueryRow(context.Background(), query, args).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.EmailVerifiedAt,
	)
	if err != nil {
		return domains.UserModel{}, err
	}

	return user, nil
}

func (u UserRepositoryPg) UpdateOneByEmail(email string, updateUser domains.UserModel) (domains.UserModel, error) {
	args := createUpdateUserPgxArgsByEmail(email, updateUser)
	query := `
		UPDATE users 
		SET
			name = COALESCE(@name, name),
			email = COALESCE(@email, email),
			password = COALESCE(@password, password),
			role = COALESCE(@role, role),
			email_verified_at = COALESCE(@email_verified_at, email_verified_at)
		WHERE email = @email
		RETURNING id, name, email, password, role, created_at, email_verified_at
	`

	var user domains.UserModel
	err := u.db.QueryRow(context.Background(), query, args).Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.EmailVerifiedAt,
	)
	if err != nil {
		return domains.UserModel{}, err
	}

	return user, nil
}

func createUpdateUserPgxArgsByEmail(email string, updateUser domains.UserModel) pgx.NamedArgs {
	result := pgx.NamedArgs{
		"email": email,
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
	if updateUser.Role != 0 {
		result["role"] = updateUser.Role
	}
	if updateUser.EmailVerifiedAt != nil {
		result["email_verified_at"] = updateUser.EmailVerifiedAt
	}

	return result
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
	if updateUser.Role != 0 {
		result["role"] = updateUser.Role
	}
	if updateUser.EmailVerifiedAt.IsZero() {
		result["email_verified_at"] = updateUser.EmailVerifiedAt
	}

	return result
}
