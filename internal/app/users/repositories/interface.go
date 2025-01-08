package repositories

import "github.com/Christian-007/fit-forge/internal/app/users/domains"

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type UserRepository interface {
	GetAll() ([]domains.UserModel, error)
	GetOne(id int) (domains.UserModel, error)
	GetOneByEmail(email string) (domains.UserModel, error)
	Create(user domains.UserModel) (domains.UserModel, error)
	Delete(id int) error
	UpdateOne(id int, updateUser domains.UserModel) (domains.UserModel, error)
	UpdateOneByEmail(email string, updateUser domains.UserModel) (domains.UserModel, error)
}
