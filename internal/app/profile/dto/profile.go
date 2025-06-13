package dto

import userdomains "github.com/Christian-007/fit-forge/internal/app/users/domains"

type ProfileResponse struct {
	UserId             int                            `json:"userId"`
	Role               int                            `json:"role"`
	SubscriptionStatus userdomains.SubscriptionStatus `json:"subscriptionStatus"`
	Name               string                         `json:"name"`
	Email              string                         `json:"email"`
}
