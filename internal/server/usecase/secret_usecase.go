// Package usecase contains the business logic for the GophKeeper application.
// It primarily handles operations related to file upload, download, and managing secret data.
package usecase

import (
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"go.uber.org/zap"
)

// MinioRepository interface defines the set of methods required to interact with Minio storage.
type MinioRepository interface {
	MultipartUploadFile(ctx context.Context,
		data *models.MultipartUploadFileData,
		partNum int,
		bucketName string,
	) (*models.ObjectPart, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
	GetObject(ctx context.Context, bucketName, filename string) (*[]byte, string, error)
	ListObjects(ctx context.Context, bucketName, prefix string) (*[]string, error)
	ApplyChanges(ctx context.Context, bucketName string) error
	DiscardChanges(ctx context.Context, bucketName string) error
}

type secretUseCase struct {
	//nolint:containedctx
	ctx        context.Context
	userRepo   UserRepository
	secretRepo MinioRepository
	logger     *zap.Logger
}

// InitSecretUseCases initializes a new secretUseCase instance.
// It prepares the use case with necessary repositories and a logger.
//
//nolint:revive,nolintlint
func InitSecretUseCases(ctx context.Context,
	userRepo UserRepository,
	secretRepo MinioRepository,
	l *zap.Logger) (*secretUseCase, error) {
	uc := &secretUseCase{
		ctx:        ctx,
		userRepo:   userRepo,
		secretRepo: secretRepo,
		logger:     l,
	}

	return uc, nil
}

// MultipartUploadFile handles the uploading of a file in multiple parts.
// It processes data from a channel and uploads each part to Minio storage.
func (uc *secretUseCase) MultipartUploadFile(ctx context.Context,
	dataChan <-chan *models.MultipartUploadFileData,
	bucketName string) (bool, error) {
	// parts := []models.ObjectPart{}
	var partNum int
	var fileName string

	uc.logger.Debug("MultipartUploadFile", zap.String("bucketname", bucketName))
	var once sync.Once

loop:
	for {
		select {
		case data, ok := <-dataChan:
			if !ok {
				uc.logger.Debug("loop getting data ended")

				break loop
			}
			//	uc.logger.Debug("we got simple data",
			// zap.Any("context", data))
			// перегружает лог при получении большого пакета данных.

			partNum++
			_, err := uc.secretRepo.MultipartUploadFile(ctx, data, partNum, bucketName)
			if err != nil {
				uc.logger.Error("Error upload in usecase", zap.Error(err))

				return false, fmt.Errorf("%w", err)
			}
			uc.logger.Debug("success upload part", zap.Int("part number", partNum))
			// parts = append(parts, *part)
			once.Do(func() {
				fileName = data.FileName
			})

		case <-ctx.Done():
			uc.logger.Debug("ctx in MultipartUploadFile usecase exeed")

			uc.logger.Debug("send empty to dataIdChan from usecase")

			return false, fmt.Errorf("%w", ctx.Err())
		}
	}

	uc.logger.Debug("Multipart upload success", zap.String("filename", fileName))

	return true, nil
}

// CreateBucketSecret creates a new bucket in Minio storage.
// It is used during the user registration process to allocate storage for the user's secrets.
func (uc *secretUseCase) CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error) {
	//nolint:wrapcheck
	return uc.secretRepo.CreateBucketSecret(ctx, username, userID)
}

// MultipartDownloadFile manages the downloading of a file in multiple parts.
// It retrieves file parts from Minio storage and sends them through a channel.
func (uc *secretUseCase) MultipartDownloadFile(ctx context.Context, req *models.MultipartDownloadFileRequest,
	bucketName string) (chan *models.MultipartDownloadFileResponse, chan error) {
	dataChan := make(chan *models.MultipartDownloadFileResponse)
	errChan := make(chan error)

	go func() {
		defer close(dataChan)
		defer close(errChan)

		// 1. Get list files by `filename-part-*`.
		prefix := fmt.Sprintf("%s-part-", req.FileName)
		objectsCh, err := uc.secretRepo.ListObjects(ctx, bucketName, prefix)
		if err != nil {
			uc.logger.Error("Error in uc.secretRepo.ListObjects")
			errChan <- err

			return
		}

		var objectNames []string
		for _, object := range *objectsCh {
			// remove files with postfix .rc
			if !strings.HasSuffix(object, ".rc") {
				uc.logger.Debug("added into objectNames", zap.String("object", object))
				objectNames = append(objectNames, object)
			}
		}

		// sort list
		sort.Slice(objectNames, func(i, j int) bool {
			iPartNum := utils.ExtractPartNumber(objectNames[i])
			jPartNum := utils.ExtractPartNumber(objectNames[j])

			return iPartNum < jPartNum
		})

		uc.logger.Debug("sorted slice", zap.Any("msg", objectNames))
		// read list and send into chan
		for _, objectName := range objectNames {
			select {
			case <-ctx.Done():
				// If the context is done, exit the goroutine
				uc.logger.Debug("Context is done")
				errChan <- ctx.Err()

				return
			default:
				objectContent, hash, err := uc.secretRepo.GetObject(ctx, bucketName, objectName)
				if err != nil {
					uc.logger.Error("Error in uc.secretRepo.GetObject", zap.Error(err))
					errChan <- err

					return
				}

				data := &models.MultipartDownloadFileResponse{
					Content:  *objectContent,
					FileName: req.FileName,
					Guid:     req.Guid,
					Hash:     hash,
				}
				dataChan <- data
			}
		}
		errChan <- io.EOF
	}()

	return dataChan, errChan
}

// ApplyChanges commits the changes made in a Minio bucket.
// It finalizes the upload process by merging temporary files into a single file.
func (uc *secretUseCase) ApplyChanges(ctx context.Context, bucketName string) (bool, error) {
	uc.logger.Debug("Start ApplyChanges in usecase")
	err := uc.secretRepo.ApplyChanges(ctx, bucketName)
	uc.logger.Debug("End ApplyChanges in usecase")
	if err != nil {
		uc.logger.Error("Error in ApplyChanges", zap.Error(err))

		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}

// DiscardChanges rolls back any changes made in a Minio bucket.
// It is used to cancel an ongoing upload process and remove temporary files.
func (uc *secretUseCase) DiscardChanges(ctx context.Context, bucketName string) (bool, error) {
	err := uc.secretRepo.DiscardChanges(ctx, bucketName)
	if err != nil {
		uc.logger.Error("Error in discard changes - usecase", zap.Error(err))

		return false, fmt.Errorf("%w", err)
	}

	return true, nil
}
