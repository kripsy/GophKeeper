// This is package of business logic level.
// Here realized logic for register, login user.
package usecase

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"go.uber.org/zap"
)

type MinioRepository interface {
	MiltipartUploadFile(context.Context, *models.MiltipartUploadFileData, int, string) (*models.ObjectPart, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
	// SaveSecret(ctx context.Context, secret entity.Secret) (int, error)
	// GetSecretByID(ctx context.Context, secretID, userID int) (entity.Secret, error)
	// DeleteSecret(ctx context.Context, secretID int) error
	// GetSecretsByUserID(ctx context.Context, userID, limit, offset int) ([]entity.Secret, error)
}

type secretUseCase struct {
	ctx        context.Context
	userRepo   UserRepository
	secretRepo MinioRepository
	logger     *zap.Logger
}

func InitSecretUseCases(ctx context.Context, userRepo UserRepository, secretRepo MinioRepository, l *zap.Logger) (*secretUseCase, error) {
	uc := &secretUseCase{
		ctx:        ctx,
		userRepo:   userRepo,
		secretRepo: secretRepo,
		logger:     l,
	}
	return uc, nil
}

func (uc *secretUseCase) MiltipartUploadFile(ctx context.Context, dataChan <-chan *models.MiltipartUploadFileData, bucketName string) (bool, error) {

	parts := []models.ObjectPart{}
	var partNum int
	var fileName string

	uc.logger.Debug("MiltipartUploadFile", zap.String("bucketname", bucketName))
	var once sync.Once

loop:
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				uc.logger.Debug("loop getting data ended")
				break loop
			}
			uc.logger.Debug("we got simple data", zap.Any("context", data))

			partNum++
			part, err := uc.secretRepo.MiltipartUploadFile(ctx, data, partNum, bucketName)
			if err != nil {
				uc.logger.Error("Error upload in usecase", zap.Error(err))

				return false, fmt.Errorf("%w", err)
			}
			uc.logger.Debug("success upload part", zap.Int("part number", partNum))
			parts = append(parts, *part)
			once.Do(func() {
				fileName = data.FileName
			})

		case <-ctx.Done():
			uc.logger.Debug("ctx in MiltipartUploadFile usecase exeed")

			uc.logger.Debug("send empty to dataIdChan from usecase")
			return false, ctx.Err()
		}
	}

	uc.logger.Debug("Multipart upload success", zap.String("filename", fileName))
	return true, nil
}

func (uc *secretUseCase) FinishSaveMultipartSecret(ctx context.Context, secret entity.Secret) (uuid.UUID, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		uc.logger.Error(err.Error())

		return uuid.UUID{}, fmt.Errorf("%w", err)
	}
	return id, nil
}

func (uc *secretUseCase) CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error) {
	return uc.secretRepo.CreateBucketSecret(ctx, username, userID)
}

// // SaveSecret saves the provided secret to the database.
// // Returns the ID of the saved secret.
// func (uc *secretUseCase) SaveTextSecret(ctx context.Context, secret entity.Secret) (int, error) {
// 	encryptedData, err := utils.Encrypt(secret.Data, uc.cipherSecret)
// 	if err != nil {
// 		return 0, err
// 	}
// 	secret.Data = encryptedData

// 	id, err := uc.db.SaveSecret(ctx, secret)
// 	if err != nil {
// 		uc.logger.Error("Error save secret in db", zap.Error(err))

// 		return 0, fmt.Errorf("%w", err)
// 	}
// 	// Сохранение секрета в репозитории

// 	return id, nil
// }

// func (uc *secretUseCase) GetSecretsByUserID(ctx context.Context, userID, limit, offset int) ([]entity.Secret, error) {
// 	secrets, err := uc.db.GetSecretsByUserID(ctx, userID, limit, offset)
// 	if err != nil {
// 		return []entity.Secret{}, fmt.Errorf("%w", err)
// 	}
// 	return secrets, nil
// }

// // d, _ := utils.Decrypt(encryptedData, uc.cipherSecret)

// // GetSecret retrieves a secret based on the provided secretID.
// func (uc *secretUseCase) GetSecretByID(ctx context.Context, secretID, userID int) (entity.Secret, error) {
// 	secret, err := uc.db.GetSecretByID(ctx, secretID, userID)
// 	if err != nil {
// 		return entity.Secret{}, fmt.Errorf("%w", err)
// 	}
// 	decryptedData, err := utils.Decrypt(secret.Data, uc.cipherSecret)
// 	if err != nil {
// 		return entity.Secret{}, err
// 	}
// 	secret.Data = decryptedData

// 	return secret, nil
// }

// // DeleteSecret deletes a secret based on the provided secretID.
// func (uc *secretUseCase) DeleteSecret(ctx context.Context, secretID int) error {
// 	err := uc.db.DeleteSecret(ctx, secretID)
// 	if err != nil {
// 		return fmt.Errorf("%w", err)
// 	}
// 	return nil
// }
