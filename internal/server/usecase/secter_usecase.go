// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/server/infrastructure"
	"go.uber.org/zap"
)

type SecretUseCase interface {
	SaveSecret(ctx context.Context, user entity.User) error
	GetSecret(ctx context.Context, secretID int) (entity.Secret, error)
}

type secretUseCase struct {
	ctx          context.Context
	db           infrastructure.SecretRepository
	logger       *zap.Logger
	cipherSecret string
}

func InitSecretUseCases(ctx context.Context, db infrastructure.SecretRepository, cipherSecret string, l *zap.Logger) (SecretUseCase, error) {
	uc := &secretUseCase{
		ctx:          ctx,
		db:           db,
		logger:       l,
		cipherSecret: cipherSecret,
	}
	return uc, nil
}

// SaveSecret saves the provided secret to the database.
// Returns the ID of the saved secret.
func (uc *secretUseCase) SaveSecret(ctx context.Context, user entity.User) error {
	return errors.New("not implemented")
}

// GetSecret retrieves a secret based on the provided secretID.
func (uc *secretUseCase) GetSecret(ctx context.Context, secretID int) (entity.Secret, error) {
	secret, err := uc.db.GetSecret(ctx, secretID)
	if err != nil {
		return entity.Secret{}, fmt.Errorf("%w", err)
	}
	return secret, nil
}

// DeleteSecret deletes a secret based on the provided secretID.
func (uc *secretUseCase) DeleteSecret(ctx context.Context, secretID int) error {
	err := uc.db.DeleteSecret(ctx, secretID)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
