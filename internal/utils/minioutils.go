package utils

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

func FromUser2BucketName(ctx context.Context, username string, userID int) (string, error) {
	str := strconv.FormatInt(int64(userID), 10)
	username = strings.ToLower(username)
	bucketName := username + str
	if bucketName == "" {
		return "", errors.New("incorrect bucket name")
	}

	return ("ilovesber" + username + str), nil
}
