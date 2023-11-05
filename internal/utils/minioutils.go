package utils

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

func FromUser2BucketName(ctx context.Context, username string, userID int) (string, error) {
	prefix := "ilovesber"
	if userID < 0 {
		return "", errors.New("incorrect userID")
	}
	str := strconv.FormatInt(int64(userID), 10)
	username = strings.ToLower(username)
	if username == "" {
		return "", errors.New("incorrect username")
	}
	bucketName := username + str
	if bucketName == "" {
		return "", errors.New("incorrect bucket name")
	}

	return (prefix + username + str), nil
}
