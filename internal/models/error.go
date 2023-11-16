// Package models defines data structures and custom error types used in the GophKeeper application.
// It includes specific error types for user-related operations and secret management.
package models

import (
	"errors"
	"fmt"
)

// Predefined errors for common failure scenarios.
var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrLoginFailed       = errors.New("login failed")
	ErrInitSecretFailed  = errors.New("failed init secret")
)

// UserExistsError represents an error when a user already exists in the system.
type UserExistsError struct {
	Text string // Text contains the custom error message.
	Err  error  // Err holds the underlying error.
}

// NewUserExistsError creates a new UserExistsError with a formatted message including the username.
func NewUserExistsError(username string) error {
	return &UserExistsError{
		Text: fmt.Sprintf("%v already exists", username),
		Err:  ErrUserAlreadyExists,
	}
}

// Error returns the error message of UserExistsError.
func (ue *UserExistsError) Error() string {
	return ue.Err.Error()
}

// UserLoginError represents an error during the user login process.
type UserLoginError struct {
	Text string // Text contains the custom error message.
	Err  error  // Err holds the underlying error.
}

// NewUserLoginError creates a new UserLoginError with a formatted message including the username.
func NewUserLoginError(username string) error {
	return &UserLoginError{
		Text: fmt.Sprintf("login failed for %v", username),
		Err:  ErrLoginFailed,
	}
}

// Error returns the error message of UserLoginError.
func (ue *UserLoginError) Error() string {
	return ue.Err.Error()
}

// SecretNotFoundError represents an error when a specific secret is not found.
type SecretNotFoundError struct {
	Text string // Text contains the custom error message.
	Err  error  // Err holds the underlying error.
}

// NewSecretNotFoundError creates a new SecretNotFoundError with a formatted message including the secret ID.
func NewSecretNotFoundError(id int) error {
	return &SecretNotFoundError{
		Text: fmt.Sprintf("failed init secret number %d", id),
		Err:  ErrInitSecretFailed,
	}
}

// Error returns the error message of SecretNotFoundError.
func (ue *SecretNotFoundError) Error() string {
	return ue.Err.Error()
}

// UnionError represents a general error for union-type scenarios.
var (
	ErrUnionError = errors.New("union error")
)

// UnionError is a custom error type for general union errors.
type UnionError struct {
	Text string // Text contains the custom error message.
	Err  error  // Err holds the underlying error.
}

// NewUnionError creates a new UnionError with a specified error text.
func NewUnionError(text string) error {
	return &SecretNotFoundError{
		Text: fmt.Sprintf("error %s", text),
		Err:  ErrUnionError,
	}
}

// Error returns the error message of UnionError.
func (ue *UnionError) Error() string {
	return ue.Err.Error()
}
