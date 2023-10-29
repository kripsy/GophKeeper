// This is package of business logic level.
// Here realized logic for register, login user.
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

type MinioRepository interface {
	MultipartUploadFile(context.Context, *models.MultipartUploadFileData, int, string) (*models.ObjectPart, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
	GetObject(ctx context.Context, bucketName, filename string) (*[]byte, string, error)
	ListObjects(ctx context.Context, bucketName, prefix string) (*[]string, error)
	CopyRCFiles(ctx context.Context, bucketName string) error
	DeleteFilesWithoutRC(ctx context.Context, bucketName string) error
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

func (uc *secretUseCase) MultipartUploadFile(ctx context.Context, dataChan <-chan *models.MultipartUploadFileData, bucketName string) (bool, error) {

	parts := []models.ObjectPart{}
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
			uc.logger.Debug("we got simple data", zap.Any("context", data))

			partNum++
			part, err := uc.secretRepo.MultipartUploadFile(ctx, data, partNum, bucketName)
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
			uc.logger.Debug("ctx in MultipartUploadFile usecase exeed")

			uc.logger.Debug("send empty to dataIdChan from usecase")
			return false, ctx.Err()
		}
	}

	uc.logger.Debug("Multipart upload success", zap.String("filename", fileName))
	return true, nil
}

func (uc *secretUseCase) CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error) {
	return uc.secretRepo.CreateBucketSecret(ctx, username, userID)
}

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

func (uc *secretUseCase) ApplyChanges(ctx context.Context, bucketName string) (bool, error) {
	uc.logger.Debug("Start ApplyChanges")
	err := uc.secretRepo.CopyRCFiles(ctx, bucketName)
	uc.logger.Debug("End CopyRCFiles")
	if err != nil {
		uc.logger.Error("Error in ApplyChanges", zap.Error(err))

		return false, err
	}

	uc.logger.Debug("Start DeleteFilesWithoutRC")
	err = uc.secretRepo.DeleteFilesWithoutRC(ctx, bucketName)
	uc.logger.Debug("End DeleteFilesWithoutRC")
	if err != nil {
		uc.logger.Error("Error in DeleteFilesWithoutRC", zap.Error(err))

		return false, err
	}

	return true, nil
}
