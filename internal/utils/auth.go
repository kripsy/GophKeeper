package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func GetHash(ctx context.Context, password string, logger *zap.Logger) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error("error GetHash", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	return string(bytes), nil

}

func BuildJWTString(userID int, secretKey string, tokenExp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed in BuildJWTString: %w")
	}

	return tokenString, nil
}

func IsPasswordCorrect(ctx context.Context, password, hashPassowrd []byte, logger *zap.Logger) error {
	err := bcrypt.CompareHashAndPassword(hashPassowrd, password)

	if err != nil {
		logger.Error("error compare password and hash", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}

	return nil
}

// Placeholder function to validate the token
func IsValidToken(tokenString string, secret string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return false, fmt.Errorf("%w", err)
	}

	if !token.Valid {
		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

func GetUserIDFromToken(tokenString, secretKey string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("%w", err)
	}

	return claims.UserID, nil
}
