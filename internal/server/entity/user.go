package entity

import (
	"bytes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/kripsy/GophKeeper/internal/utils"
)

type User struct {
	Username string `json:"login" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

func InitNewUser(data []byte) (User, error) {
	var buf bytes.Buffer
	buf.Write(data)
	decoder := json.NewDecoder(&buf)
	decoder.DisallowUnknownFields()

	u := User{}
	err := decoder.Decode(&u)
	if err != nil {
		return User{}, err
	}

	validate := validator.New()
	err = validate.Struct(&u)
	if err != nil {
		return User{}, err
	}

	err = utils.CheckDuplicateFields(data, &User{})
	if err != nil {
		return User{}, err
	}

	return u, nil
}
