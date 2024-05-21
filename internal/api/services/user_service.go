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
		userResponses[i] = toUserResponse(user)
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

	return toUserResponse(user), nil
}

func (u UserService) Create(createUserRequest dto.CreateUserRequest) (dto.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserRequest.Password), 12)
	if err != nil {
		return dto.UserResponse{}, err
	}

	userModel := domains.UserModel{
		Name:     createUserRequest.Name,
		Email:    createUserRequest.Email,
		Password: hashedPassword,
	}

	createdUser, err := u.UserRepository.Create(userModel)
	if err != nil {
		return dto.UserResponse{}, err
	}

	return toUserResponse(createdUser), nil
}

func (u UserService) Delete(id int) error {
	err := u.UserRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (u UserService) UpdateOne(id int, updateUserRequest dto.UpdateUserRequest) (dto.UserResponse, error) {
	userModel, err := toUserModel(updateUserRequest)
	if err != nil {
		return dto.UserResponse{}, err
	}

	user, err := u.UserRepository.UpdateOne(id, userModel)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserResponse{}, apperrors.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func toUserModel(updateUserRequest dto.UpdateUserRequest) (domains.UserModel, error) {
	var userModel domains.UserModel

	if updateUserRequest.Email != nil {
		userModel.Email = *updateUserRequest.Email
	}

	if updateUserRequest.Name != nil {
		userModel.Name = *updateUserRequest.Name
	}

	if updateUserRequest.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*updateUserRequest.Password), 12)
		if err != nil {
			return domains.UserModel{}, err
		}

		userModel.Password = hashedPassword
	}

	return userModel, nil
}

func toUserResponse(userModel domains.UserModel) dto.UserResponse {
	return dto.UserResponse{
		Id:    userModel.Id,
		Name:  userModel.Name,
		Email: userModel.Email,
	}
}
