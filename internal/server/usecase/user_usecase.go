// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
)

// UserRepository interface defines the required methods for user data operations.
type UserRepository interface {
	RegisterUser(ctx context.Context, user entity.User) (int, error)
	LoginUser(ctx context.Context, user entity.User) (int, error)
}

type userUseCase struct {
	//nolint:containedctx
	ctx      context.Context
	db       UserRepository
	logger   *zap.Logger
	secret   string
	tokenExp time.Duration
}

// InitUseCases initializes a new userUseCase instance.
// It prepares the use case with the user repository, JWT secret, token expiration, and a logger.
//
//nolint:revive,nolintlint
func InitUseCases(ctx context.Context,
	db UserRepository,
	secret string,
	tokenExp time.Duration,
	l *zap.Logger) (*userUseCase, error) {
	uc := &userUseCase{
		ctx:      ctx,
		db:       db,
		logger:   l,
		secret:   secret,
		tokenExp: tokenExp,
	}

	return uc, nil
}

// RegisterUser handles the registration of a new user.
// It checks if the user exists, registers them if not, and then generates a JWT token.
func (uc *userUseCase) RegisterUser(ctx context.Context, user entity.User) (string, int, error) {
	userID, err := uc.db.RegisterUser(ctx, user)
	if err != nil {
		return "", 0, fmt.Errorf("%w", err)
	}

	token, err := utils.BuildJWTString(userID, user.Username, uc.secret, uc.tokenExp)
	if err != nil {
		return "", 0, fmt.Errorf("%w", err)
	}

	return token, userID, nil
}

// LoginUser authenticates a user based on the provided credentials.
// It verifies the user's existence and, upon successful authentication, generates a JWT token.
func (uc *userUseCase) LoginUser(ctx context.Context, user entity.User) (string, int, error) {
	userID, err := uc.db.LoginUser(ctx, user)
	if err != nil {
		return "", 0, fmt.Errorf("%w", err)
	}

	token, err := utils.BuildJWTString(userID, user.Username, uc.secret, uc.tokenExp)
	if err != nil {
		return "", 0, fmt.Errorf("%w", err)
	}

	return token, userID, nil
}
