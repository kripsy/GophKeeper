package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

const (
	AUTHORIZATIONHEADER = "authorization"
	TOKENPREFIX         = ""
	TOKENCONTEXTKEY     = "token"
	USERNAMECONTEXTKEY  = "username"
	USERIDCONTEXTKEY    = "userID"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID   int
	Username string
}

func GetHash(_ context.Context, password string, logger *zap.Logger) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		logger.Error("error GetHash", zap.String("msg", err.Error()))

		return "", fmt.Errorf("%w", err)
	}

	return string(bytes), nil
}

func BuildJWTString(userID int, username, secretKey string, tokenExp time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID:   userID,
		Username: username,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.Wrap(err, "failed in BuildJWTString: %w")
	}

	return tokenString, nil
}

func IsPasswordCorrect(_ context.Context, password, hashPassowrd []byte, logger *zap.Logger) error {
	err := bcrypt.CompareHashAndPassword(hashPassowrd, password)

	if err != nil {
		logger.Error("error compare password and hash", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}

	return nil
}

// Placeholder function to validate the token.
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

func GetUsernameFromToken(tokenString, secretKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("%w", err)
	}

	return claims.Username, nil
}
func GetUseIDFromToken(tokenString, secretKey string) (int, error) {
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

func DeriveKey(password, salt string) ([]byte, error) {
	const (
		N       = 32768
		r       = 8
		p       = 1
		saltLen = 16
		keyLen  = 32
	)

	key, err := scrypt.Key([]byte(password), []byte(salt), N, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return key, nil
}

func ExtractTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(TOKENCONTEXTKEY).(string)

	return token, ok
}

func ExtractUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(USERNAMECONTEXTKEY).(string)

	return username, ok
}

func ExtractUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(USERIDCONTEXTKEY).(int)

	return userID, ok
}
