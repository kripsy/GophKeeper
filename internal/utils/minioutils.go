package utils

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Predefined errors for validation failures in utility functions.
var (
	ErrIncorrectBucket   = errors.New("incorrect bucket name")
	ErrIncorrectUserName = errors.New("incorrect user name")
	ErrIncorrectUserID   = errors.New("incorrect user id")
)

// FromUser2BucketName generates a bucket name based on the provided username and user ID.
//
// Parameters:
// - ctx: A context that may be used for cancellation or deadline control (not currently used in function).
// - username: The username of the user.
// - userID: The ID of the user.
//
// Returns:
// - A string representing the generated bucket name.
// - An error if the username is empty or the user ID is invalid.
func FromUser2BucketName(_ context.Context, username string, userID int) (string, error) {
	prefix := "ilovesber"
	if userID < 0 {
		return "", fmt.Errorf("%w", ErrIncorrectUserID)
	}
	str := strconv.FormatInt(int64(userID), 10)
	username = strings.ToLower(username)
	if username == "" {
		return "", fmt.Errorf("%w", ErrIncorrectUserName)
	}

	return (prefix + username + str), nil
}
