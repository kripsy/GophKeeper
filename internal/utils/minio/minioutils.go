package minio

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrIncorrectBucket   = errors.New("incorrect bucket name")
	ErrIncorrectUserName = errors.New("incorrect user name")
	ErrIncorrectUserID   = errors.New("incorrect user id")
)

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
	bucketName := username + str
	if bucketName == "" {
		return "", fmt.Errorf("%w", ErrIncorrectBucket)
	}

	return (prefix + username + str), nil
}
