package services

import (
	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/dto"
	"github.com/Christian-007/fit-forge/internal/api/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserServiceOptions
}

type UserServiceOptions struct {
	UserRepository repositories.UserRepository
}

func NewUserService(options UserServiceOptions) UserService {
	return UserService{
		options,
	}
}

func (u UserService) Create(createUserRequest dto.CreateUserRequest) (dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 12)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user := domains.User{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: hashedPassword,
	}

	userDb, err := u.UserRepository.Create(user)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Id:    userDb.Id,
		Name:  userDb.Name,
		Email: userDb.Email,
	}, nil
}
