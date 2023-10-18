// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"fmt"

	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
)

type SecretRepository interface {
	SaveSecret(ctx context.Context, secret entity.Secret) (int, error)
	GetSecretByID(ctx context.Context, secretID, userID int) (entity.Secret, error)
	DeleteSecret(ctx context.Context, secretID int) error
	GetSecretsByUserID(ctx context.Context, userID, limit, offset int) ([]entity.Secret, error)
}

type secretUseCase struct {
	ctx          context.Context
	db           SecretRepository
	logger       *zap.Logger
	cipherSecret string
}

func InitSecretUseCases(ctx context.Context, db SecretRepository, cipherSecret string, l *zap.Logger) (*secretUseCase, error) {
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
func (uc *secretUseCase) SaveTextSecret(ctx context.Context, secret entity.Secret) (int, error) {
	encryptedData, err := utils.Encrypt(secret.Data, uc.cipherSecret)
	if err != nil {
		return 0, err
	}
	secret.Data = encryptedData

	id, err := uc.db.SaveSecret(ctx, secret)
	if err != nil {
		uc.logger.Error("Error save secret in db", zap.Error(err))

		return 0, fmt.Errorf("%w", err)
	}
	// Сохранение секрета в репозитории

	return id, nil
}

func (uc *secretUseCase) GetSecretsByUserID(ctx context.Context, userID, limit, offset int) ([]entity.Secret, error) {
	secrets, err := uc.db.GetSecretsByUserID(ctx, userID, limit, offset)
	if err != nil {
		return []entity.Secret{}, fmt.Errorf("%w", err)
	}
	return secrets, nil
}

// d, _ := utils.Decrypt(encryptedData, uc.cipherSecret)

// GetSecret retrieves a secret based on the provided secretID.
func (uc *secretUseCase) GetSecretByID(ctx context.Context, secretID, userID int) (entity.Secret, error) {
	secret, err := uc.db.GetSecretByID(ctx, secretID, userID)
	if err != nil {
		return entity.Secret{}, fmt.Errorf("%w", err)
	}
	decryptedData, err := utils.Decrypt(secret.Data, uc.cipherSecret)
	if err != nil {
		return entity.Secret{}, err
	}
	secret.Data = decryptedData

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
