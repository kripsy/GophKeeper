package models

import (
	"errors"
	"fmt"
)

type UserExistsError struct {
	Text string
	Err  error
}

func NewUserExistsError(username string) error {
	return &UserExistsError{
		Text: fmt.Sprintf("%v already exists", username),
		Err:  errors.New("user already exists"),
	}
}

func (ue *UserExistsError) Error() string {
	return ue.Err.Error()
}

type UserLoginError struct {
	Text string
	Err  error
}

func NewUserLoginError(username string) error {
	return &UserLoginError{
		Text: fmt.Sprintf("login failed for %v", username),
		Err:  errors.New("login failed"),
	}
}

func (ue *UserLoginError) Error() string {
	return ue.Err.Error()
}

type SecretNotFoundError struct {
	Text string
	Err  error
}

func NewSecretNotFoundError(ID int) error {
	return &SecretNotFoundError{
		Text: fmt.Sprintf("failed init secret number %d", ID),
		Err:  errors.New("failed init secret"),
	}
}

func (ue *SecretNotFoundError) Error() string {
	return ue.Err.Error()
}

type UnionError struct {
	Text string
	Err  error
}

func NewUnionError(text string) error {
	return &SecretNotFoundError{
		Err: errors.New(fmt.Sprintf("error %s", text)),
	}
}

func (ue *UnionError) Error() string {
	return ue.Err.Error()
}
