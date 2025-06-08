package services

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/pkg/apperrors"
	"github.com/Christian-007/fit-forge/internal/pkg/cache"
	"github.com/Christian-007/fit-forge/internal/pkg/security"
)

type EmailService struct {
	EmailServiceOptions
}

type EmailServiceOptions struct {
	Host         string
	TokenService security.TokenService
	Cache        cache.Cache
}

func NewEmailService(options EmailServiceOptions) EmailService {
	return EmailService{options}
}

func (e EmailService) CreateVerificationLink(email string) (string, error) {
	verificationPath := "/email-verification?token="

	randomToken, err := e.TokenService.Generate()
	if err != nil {
		return "", err
	}

	err = e.Cache.Set(randomToken.Hashed, email, time.Hour*24)
	if err != nil {
		return "", err
	}

	link := e.Host + verificationPath + randomToken.Raw
	return link, nil
}

func (e EmailService) Verify(rawToken string) (string, string, error) {
	hashed, err := e.TokenService.HashWithSecret(rawToken)
	if err != nil {
		return "", "", err
	}

	email, err := e.Cache.Get(hashed)
	if err != nil {
		return "", "", apperrors.ErrRedisKeyNotFound
	}

	stringEmail, ok := email.(string)
	if !ok {
		return "", "", apperrors.ErrTypeAssertion
	}

	return stringEmail, hashed, nil
}
