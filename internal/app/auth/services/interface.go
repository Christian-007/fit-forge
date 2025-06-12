package services

import (
	"crypto/rsa"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
)

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type AuthService interface {
	Authenticate(inputtedPassword string, userPassword []byte) error
	CreateToken(privateKey *rsa.PrivateKey, userId int) (domains.AuthToken, error)
	ValidateToken(privateKey *rsa.PrivateKey, tokenString string) (*domains.Claims, error)
	InvalidateToken(accessTokenUuid string) error
	SaveToken(userResponse userdto.UserResponse, authToken domains.AuthToken) error
	GetHashAuthDataFromCache(accessTokenUuid string) (domains.AuthData, error)
	GetAuthDataFromCache(accessTokenUuid string) (int, error)
}
