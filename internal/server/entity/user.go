package entity

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/models"
)

type User struct {
	Username string `json:"login"`
	Password string `json:"password,omitempty"`
}

func InitNewUser(username, password string) (User, error) {
	if len(username) < 3 || len(password) < 3 {
		return User{}, fmt.Errorf("%w", models.NewUnionError("error init user object"))
	}
	u := User{
		Username: username,
		Password: password,
	}

	return u, nil
}
