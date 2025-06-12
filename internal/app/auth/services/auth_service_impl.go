package services

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Christian-007/fit-forge/internal/app/auth/domains"
	usersdomain "github.com/Christian-007/fit-forge/internal/app/users/domains"
	userdto "github.com/Christian-007/fit-forge/internal/app/users/dto"
	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	AuthServiceOptions
}

type AuthServiceOptions struct {
	Cache       cache.Cache
}

func NewAuthServiceImpl(options AuthServiceOptions) AuthServiceImpl {
	return AuthServiceImpl{
		options,
	}
}

func (a AuthServiceImpl) Authenticate(inputtedPassword string, userPassword []byte) error {
	err := bcrypt.CompareHashAndPassword(userPassword,[]byte(inputtedPassword))
	if err != nil {
		return apperrors.ErrInvalidCredentials
	}

	return nil
}

func (a AuthServiceImpl) CreateToken(privateKey *rsa.PrivateKey, userId int) (domains.AuthToken, error) {
	authToken := domains.AuthToken{}
	uuid := uuid.NewString()
	expiresAt := jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
	claims := domains.Claims{
		UserID: userId,
		Uuid:   uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER_CLAIM"),
			Audience:  jwt.ClaimStrings{os.Getenv("JWT_AUDIENCE_CLAIM")},
			Subject:   fmt.Sprintf("%d", userId),
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Set the 'kid' (Key ID) header so verifiers know which public key to use from JWKS
	token.Header["kid"] = os.Getenv("JWK_KEY_ID")

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return domains.AuthToken{}, err
	}

	authToken.AccessToken = tokenString
	authToken.AccessTokenUuid = uuid
	authToken.AccessTokenExpiresAt = expiresAt

	return authToken, nil
}

func (a AuthServiceImpl) ValidateToken(privateKey *rsa.PrivateKey, tokenString string) (*domains.Claims, error) {
	publicKey := &privateKey.PublicKey
	claims := &domains.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		return publicKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg()}))
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, apperrors.ErrInvalidSignature
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperrors.ErrExpiredToken
		}

		return nil, err
	}

	if !token.Valid {
		return nil, apperrors.ErrInvalidToken
	}

	return claims, nil
}

func (a AuthServiceImpl) InvalidateToken(accessTokenUuid string) error {
	err := a.Cache.Delete(accessTokenUuid)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthServiceImpl) SaveToken(userResponse userdto.UserResponse, authToken domains.AuthToken) error {
	accessTokenExpiration := time.Until(authToken.AccessTokenExpiresAt.Time)
	err := a.Cache.SetHash(authToken.AccessTokenUuid,
		"userId", userResponse.Id,
		"role", userResponse.Role,
		"name", userResponse.Name,
		"email", userResponse.Email,
		"subscriptionStatus", string(userResponse.SubscriptionStatus),
	)
	if err != nil {
		return err
	}

	err = a.Cache.SetExpire(authToken.AccessTokenUuid, accessTokenExpiration)
	if err != nil {
		return err
	}

	return nil
}

func (a AuthServiceImpl) GetHashAuthDataFromCache(accessTokenUuid string) (domains.AuthData, error) {
	result, err := a.Cache.GetAllHashFields(accessTokenUuid)
	if len(result) == 0 {
		return domains.AuthData{}, apperrors.ErrRedisValueNotInHash
	}

	if err != nil {
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

	subscriptionStatus := usersdomain.SubscriptionStatus(result["subscriptionStatus"])

	return domains.AuthData{
		UserId:             userIdInt,
		Role:               roleInt,
		SubscriptionStatus: subscriptionStatus,
	}, nil
}

// (Obsolete) Get auth data from Cache with a single value.
// Use `.GetHashAuthDataFromCache(accessTokenUuid string)` instead.
func (a AuthServiceImpl) GetAuthDataFromCache(accessTokenUuid string) (int, error) {
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
