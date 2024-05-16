package services

import (
	"errors"

	"github.com/Christian-007/fit-forge/internal/api/apperrors"
	"github.com/Christian-007/fit-forge/internal/api/domains"
	"github.com/Christian-007/fit-forge/internal/api/dto"
	"github.com/Christian-007/fit-forge/internal/api/repositories"
	"github.com/jackc/pgx/v5"
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

func (u UserService) GetAll() ([]dto.UserResponse, error) {
	users, err := u.UserRepository.GetAll()
	if err != nil {
		return []dto.UserResponse{}, err
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = dto.UserResponse{
			Id:    user.Id,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return userResponses, nil
}

func (u UserService) GetOne(id int) (dto.UserResponse, error) {
	user, err := u.UserRepository.GetOne(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserResponse{}, apperrors.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	return dto.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (u UserService) Create(createUserRequest dto.CreateUserRequest) (dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 12)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user := domains.UserModel{
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
