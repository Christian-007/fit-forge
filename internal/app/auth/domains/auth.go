package domains

import (
	usersdomain "github.com/Christian-007/fit-forge/internal/app/users/domains"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int `json:"userId"`
	Uuid   string
	jwt.RegisteredClaims
}

type AuthToken struct {
	AccessToken          string
	AccessTokenUuid      string
	AccessTokenExpiresAt *jwt.NumericDate
}

type AuthData struct {
	UserId             int                            `json:"userId"`
	Role               int                            `json:"role"`
	SubscriptionStatus usersdomain.SubscriptionStatus `json:"subscriptionStatus"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
}
