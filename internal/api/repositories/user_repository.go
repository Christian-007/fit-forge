package repositories

import (
	"time"

	"github.com/Christian-007/fit-forge/internal/api/domains"
)

type UserRepository struct {
	users []domains.User
}

func NewUserRepository() UserRepository {
	return UserRepository{
		users: []domains.User{
			{
				Id:        12,
				Name:      "John Doe",
				Email:     "jdoe@example.com",
				Password:  "test",
				CreatedAt: time.Time{},
			},
		},
	}
}

func (u UserRepository) GetAll() []domains.User {
	return u.users
}
