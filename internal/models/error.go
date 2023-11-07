package models

import (
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrLoginFailed       = errors.New("login failed")
	ErrInitSecretFailed  = errors.New("failed init secret")
)

type UserExistsError struct {
	Text string
	Err  error
}

func NewUserExistsError(username string) error {
	return &UserExistsError{
		Text: fmt.Sprintf("%v already exists", username),
		Err:  ErrUserAlreadyExists,
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
		Err:  ErrLoginFailed,
	}
}

func (ue *UserLoginError) Error() string {
	return ue.Err.Error()
}

type SecretNotFoundError struct {
	Text string
	Err  error
}

func NewSecretNotFoundError(id int) error {
	return &SecretNotFoundError{
		Text: fmt.Sprintf("failed init secret number %d", id),
		Err:  ErrInitSecretFailed,
	}
}

func (ue *SecretNotFoundError) Error() string {
	return ue.Err.Error()
}

var (
	ErrUnionError = errors.New("union error")
)

type UnionError struct {
	Text string
	Err  error
}

func NewUnionError(text string) error {
	return &SecretNotFoundError{
		Text: fmt.Sprintf("error %s", text),
		Err:  ErrUnionError,
	}
}

func (ue *UnionError) Error() string {
	return ue.Err.Error()
}
