package entity

import (
	"fmt"

	"github.com/kripsy/GophKeeper/internal/models"
)

// User represents a user in the GophKeeper application.
type User struct {
	Username string `json:"login"`
	Password string `json:"password,omitempty"`
}

// InitNewUser initializes a new User instance.
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
