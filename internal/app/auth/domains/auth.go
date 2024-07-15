package domains

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID int `json:"userId"`
	jwt.RegisteredClaims
}
