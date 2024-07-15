package services

import (
	"os"
	"time"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	authdto "github.com/Christian-007/fit-forge/internal/app/auth/dto"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthServiceOptions
}

type AuthServiceOptions struct {
	UserService services.UserService
}

func NewAuthService(options AuthServiceOptions) AuthService {
	return AuthService{
		options,
	}
}

func (a AuthService) Authenticate(loginRequest authdto.LoginRequest) (userdto.UserResponse, error) {
	user, err := a.UserService.GetOneByEmail(loginRequest.Username)
	if err != nil {
		return userdto.UserResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(loginRequest.Password))
	if err != nil {
		return userdto.UserResponse{}, apperrors.ErrInvalidCredentials
	}

	response := userdto.UserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
	}
	return response, nil
}

func (a AuthService) CreateToken(userId int) (string, error) {
	secretKey := []byte(os.Getenv("AUTH_SECRET_KEY"))
	claims := domains.Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a AuthService) ValidateToken(tokenString string) (*domains.Claims, error){
	secretKey := []byte(os.Getenv("AUTH_SECRET_KEY"))
	claims := &domains.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, apperrors.ErrInvalidSignature
		}
		return nil, err
	}

	if !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	return claims, nil
}
