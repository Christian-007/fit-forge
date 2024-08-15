package domains

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}

type AuthToken struct {
	AccessToken          string
	AccessTokenUuid      string
	AccessTokenExpiresAt *jwt.NumericDate
}
