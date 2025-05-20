package dto

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/app/points/domains"
	usersdomain "github.com/Christian-007/fit-forge/internal/app/users/domains"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserResponse struct {
	Id                 int                            `json:"id"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
	Role               int                            `json:"role"`
	SubscriptionStatus usersdomain.SubscriptionStatus `json:"subscriptionStatus"`
	EmailVerifiedAt    *time.Time                     `json:"emailVerifiedAt"`
}

type UserWithPointsResponse struct {
	Id                 int                            `json:"id"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
	Role               int                            `json:"role"`
	SubscriptionStatus usersdomain.SubscriptionStatus `json:"subscriptionStatus"`
	EmailVerifiedAt    *time.Time                     `json:"emailVerifiedAt"`
	Point              domains.PointModel             `json:"point"`
}

type GetUserByEmailResponse struct {
	Id                 int                            `json:"id"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
	Role               int                            `json:"role"`
	SubscriptionStatus usersdomain.SubscriptionStatus `json:"subscriptionStatus"`
	Password           []byte                         `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required, validation.Length(2, 200)),
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required),
	)
}

type UpdateUserRequest struct {
	Name               *string                         `json:"name"`
	Email              *string                         `json:"email"`
	Password           *string                         `json:"password"`
	Role               *int                            `json:"role"`
	SubscriptionStatus *usersdomain.SubscriptionStatus `json:"subscriptionStatus"`
	EmailVerifiedAt    *time.Time                      `json:"emailVerifiedAt"`
}

// As of 28 May '24 does not support an empty string update
func (u UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.NilOrNotEmpty, validation.Length(2, 200)),
		validation.Field(&u.Email, validation.NilOrNotEmpty, is.Email),
		validation.Field(&u.Role, validation.NilOrNotEmpty, validation.In(1, 2)),
		validation.Field(&u.Role, validation.NilOrNotEmpty, validation.In(usersdomain.InactiveSubscriptionStatus, usersdomain.ActiveSubscriptionStatus)),
	)
}
