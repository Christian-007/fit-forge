package services

import (
	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	authdto "github.com/Christian-007/fit-forge/internal/app/auth/dto"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
)

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type AuthService interface {
	Authenticate(loginRequest authdto.LoginRequest) (userdto.UserResponse, error)
	CreateToken(userId int) (domains.AuthToken, error)
	ValidateToken(tokenString string) (*domains.Claims, error)
	InvalidateToken(accessTokenUuid string) error
	SaveToken(userResponse userdto.UserResponse, authToken domains.AuthToken) error
	GetHashAuthDataFromCache(accessTokenUuid string) (domains.AuthData, error)
	GetAuthDataFromCache(accessTokenUuid string) (int, error)
}
