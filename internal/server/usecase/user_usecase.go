// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/server/infrastructure"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
)

type UserUseCase interface {
	RegisterUser(ctx context.Context, user entity.User) (string, error)
	LoginUser(ctx context.Context, user entity.User) (string, error)
}

type userUseCase struct {
	ctx      context.Context
	db       infrastructure.UserRepository
	logger   *zap.Logger
	secret   string
	tokenExp time.Duration
}

func InitUseCases(ctx context.Context, db infrastructure.UserRepository, secret string, tokenExp time.Duration, l *zap.Logger) (UserUseCase, error) {
	uc := &userUseCase{
		ctx:      ctx,
		db:       db,
		logger:   l,
		secret:   secret,
		tokenExp: tokenExp,
	}
	return uc, nil
}

// RegisterUser get context, user and return token, expired time, error.
// At the first step we check is user exists. If exists - return error conflict.
// If user not exists we get new user ID.
// After register new user we generate new jwt token.
func (uc *userUseCase) RegisterUser(ctx context.Context, user entity.User) (string, error) {
	userID, err := uc.db.RegisterUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	token, err := utils.BuildJWTString(userID, uc.secret, uc.tokenExp)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return token, nil
}

// LoginUser authenticates a user based on the provided credentials.
// Upon successful authentication, the function generates a JWT token for the user.
//
// Parameters:
// - ctx: The context for the operation, which can be used for cancellation or timeout.
// - user: The User structure containing the user's credentials for authentication.
//
// Returns:
// - string: The generated JWT token that can be used for subsequent authorization.
// - error: An error that might occur during authentication or token generation.
//
// Usage example:
// token, err := uc.LoginUser(ctx, user)
//
//	if err != nil {
//	    log.Fatalf("Failed to login user: %v", err)
//	}
//
// fmt.Println("Generated JWT token:", token)
func (uc *userUseCase) LoginUser(ctx context.Context, user entity.User) (string, error) {
	userID, err := uc.db.LoginUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	token, err := utils.BuildJWTString(userID, uc.secret, uc.tokenExp)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return token, nil
}
