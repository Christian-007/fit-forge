package services

import (
	"strconv"
	"time"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	authdto "github.com/Christian-007/fit-forge/internal/app/auth/dto"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/app/users/services"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/envvariable"
	"github.com/golang-jwt/jwt/v5"
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthServiceOptions
}

type AuthServiceOptions struct {
	UserService        services.UserService
	Cache              cache.Cache
	EnvVariableService envvariable.EnvVariableService
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
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}
	return response, nil
}

func (a AuthService) CreateToken(userId int) (domains.AuthToken, error) {
	authToken := domains.AuthToken{}
	uuid := uuid.NewV4().String()
	secretKey := []byte(a.EnvVariableService.Get("AUTH_SECRET_KEY"))
	expiresAt := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims := domains.Claims{
		UserID: userId,
		Uuid:   uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return domains.AuthToken{}, err
	}

	authToken.AccessToken = tokenString
	authToken.AccessTokenUuid = uuid
	authToken.AccessTokenExpiresAt = expiresAt

	return authToken, nil
}

func (a AuthService) ValidateToken(tokenString string) (*domains.Claims, error) {
	secretKey := []byte(a.EnvVariableService.Get("AUTH_SECRET_KEY"))
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

func (a AuthService) InvalidateToken(accessTokenUuid string) error {
	err := a.Cache.Delete(accessTokenUuid)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthService) SaveToken(userResponse userdto.UserResponse, authToken domains.AuthToken) error {
	accessTokenExpiration := time.Until(authToken.AccessTokenExpiresAt.Time)
	err := a.Cache.SetHash(authToken.AccessTokenUuid, "userId", userResponse.Id, "role", userResponse.Role)
	if err != nil {
		return err
	}

	err = a.Cache.SetExpire(authToken.AccessTokenUuid, accessTokenExpiration)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthService) GetHashAuthDataFromCache(accessTokenUuid string) (domains.AuthData, error) {
	result, err := a.Cache.GetAllHashFields(accessTokenUuid)
	if err != nil {
		if len(result) == 0 {
			return domains.AuthData{}, apperrors.ErrRedisValueNotInHash
		}

		return domains.AuthData{}, err
	}

	userIdInt, err := strconv.Atoi(result["userId"])
	if err != nil {
		return domains.AuthData{}, err
	}

	roleInt, err := strconv.Atoi(result["role"])
	if err != nil {
		return domains.AuthData{}, err
	}

	return domains.AuthData{
		UserId: userIdInt,
		Role:   roleInt,
	}, nil
}

// (Obsolete) Get auth data from Cache with a single value.
// Use `.GetHashAuthDataFromCache(accessTokenUuid string)` instead.
func (a AuthService) GetAuthDataFromCache(accessTokenUuid string) (int, error) {
	value, err := a.Cache.Get(accessTokenUuid)
	if err != nil {
		return -1, err
	}

	stringUserId, ok := value.(string)
	if !ok {
		return -1, apperrors.ErrTypeAssertion
	}

	userId, err := strconv.Atoi(stringUserId)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
