package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type RandomToken struct {
	Raw string
	Hashed string
}

type TokenService struct {
	TokenServiceOptions
}

type TokenServiceOptions struct {
	SecretKey string
}

func NewTokenService(options TokenServiceOptions) TokenService {
	return TokenService{options}
}

func (t TokenService) Generate() (RandomToken, error) {
	byteLength := 16
	token := make([]byte, byteLength)
	_, err := rand.Read(token)
	if err != nil {
		// The generated token doesn't have a length of 16
		return RandomToken{}, err
	}

	rawToken := hex.EncodeToString(token)

	hashedToken, err := t.HashWithSecret(rawToken)
	if err != nil {
		return RandomToken{}, err
	}

	return RandomToken{
		Raw: rawToken,
		Hashed: hashedToken,
	}, nil
}

func (t TokenService) HashWithSecret(token string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(token + t.SecretKey))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
