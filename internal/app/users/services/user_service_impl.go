package services

import (
	"errors"

	"github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/app/users/repositories"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	UserServiceOptions
}

type UserServiceOptions struct {
	UserRepository repositories.UserRepository
}

func NewUserService(options UserServiceOptions) UserServiceImpl {
	return UserServiceImpl{
		options,
	}
}

func (u UserServiceImpl) GetAll() ([]dto.UserResponse, error) {
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

func (u UserServiceImpl) GetOne(id int) (dto.UserResponse, error) {
	user, err := u.UserRepository.GetOne(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserResponse{}, apperrors.ErrUserNotFound
		}

		return dto.UserResponse{}, err
	}

	return toUserResponse(user), nil
}

func (u UserServiceImpl) GetOneByEmail(email string) (dto.GetUserByEmailResponse, error) {
	userModel, err := u.UserRepository.GetOneByEmail(email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.GetUserByEmailResponse{}, apperrors.ErrUserNotFound
		}

		return dto.GetUserByEmailResponse{}, err
	}

	response := dto.GetUserByEmailResponse{
		Id:       userModel.Id,
		Name:     userModel.Name,
		Email:    userModel.Email,
		Password: userModel.Password,
	}
	return response, nil
}

func (u UserServiceImpl) Create(createUserRequest dto.CreateUserRequest) (dto.UserResponse, error) {
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

func (u UserServiceImpl) Delete(id int) error {
	err := u.UserRepository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (u UserServiceImpl) UpdateOne(id int, updateUserRequest dto.UpdateUserRequest) (dto.UserResponse, error) {
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