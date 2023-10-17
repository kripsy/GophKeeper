// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"errors"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/server/infrastructure"
	"go.uber.org/zap"
)

type SecretUseCase interface {
	AddSecret(ctx context.Context, user entity.User) error
}

type secretUseCase struct {
	ctx          context.Context
	db           infrastructure.SecretRepository
	logger       *zap.Logger
	cipherSecret string
}

func InitSecretUseCases(ctx context.Context, db infrastructure.UserRepository, cipherSecret string, l *zap.Logger) (SecretUseCase, error) {
	uc := &secretUseCase{
		ctx:          ctx,
		db:           db,
		logger:       l,
		cipherSecret: cipherSecret,
	}
	return uc, nil
}

// RegisterUser get context, user and return token, expired time, error.
// At the first step we check is user exists. If exists - return error conflict.
// If user not exists we get new user ID.
// After register new user we generate new jwt token.
func (uc *secretUseCase) AddSecret(ctx context.Context, user entity.User) error {
	return errors.New("not implemented")
}
