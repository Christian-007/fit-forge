package services

import (
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
)

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type UserService interface {
	GetAll() ([]dto.UserResponse, error)
	GetOne(id int) (dto.UserResponse, error)
	GetOneByEmail(email string) (dto.GetUserByEmailResponse, error)
	Create(createUserRequest dto.CreateUserRequest) (dto.UserResponse, error)
	Delete(id int) error
	UpdateOne(id int, updateUserRequest dto.UpdateUserRequest) (dto.UserResponse, error)
}
