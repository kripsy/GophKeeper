package infrastructure

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinioRepository interface {
	MiltipartUploadFile(context.Context, *models.MiltipartUploadFileData, int, string) (*models.ObjectPart, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
}

type minioRepository struct {
	client *minio.Client
	logger *zap.Logger
}

func NewMinioRepository(ctx context.Context, endpoint, accessKeyID, secretAccessKey, bucketName string, isUseSSL bool, l *zap.Logger) (MinioRepository, error) {
	l.Debug("start init minio repository")
	client, err := initMinioClient(endpoint, accessKeyID, secretAccessKey, isUseSSL)

	if err != nil {
		l.Error("Error init minio client", zap.Error(err))

		return nil, fmt.Errorf("%w", err)
	}
	l.Debug("minio client created")

	l.Debug("check init bucket")
	err = initBucket(ctx, bucketName, client)
	if err != nil {
		l.Error("Error init bucket", zap.Error(err))

		return nil, fmt.Errorf("%w", err)
	}

	l.Debug("check init bucket finished")

	return &minioRepository{
		client: client,
		logger: l,
	}, nil
}

func initMinioClient(endpoint, accessKeyID, secretAccessKey string, isUseSSL bool) (*minio.Client, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: isUseSSL,
	})
	if err != nil {
		fmt.Printf("Error init client minio: %s", err.Error())

		return nil, fmt.Errorf("%w", err)
	}

	return minioClient, nil
}

func initBucket(ctx context.Context, bucketName string, client *minio.Client) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (m *minioRepository) MiltipartUploadFile(ctx context.Context, data *models.MiltipartUploadFileData, partNum int, bucketName string) (*models.ObjectPart, error) {
	objectName := fmt.Sprintf("%s-part-%d", data.FileName, partNum)
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"Hash":     data.Hash,
			"Username": data.Username,
		},
	}

	_, err := m.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(data.Content), int64(len(data.Content)), opts)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &models.ObjectPart{
		PartNumber: partNum,
	}, nil
}

func (m *minioRepository) CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error) {
	bucketName, err := utils.FromUser2BucketName(ctx, username, userID)
	if err != nil {
		m.logger.Error("Error in CreateBucketSecret", zap.Error(err))

		return false, fmt.Errorf("%w", err)
	}
	m.logger.Debug("check bucket exists", zap.String("bucketName", bucketName))
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		m.logger.Error("Error check bucket exists", zap.Error(err))

		return false, fmt.Errorf("%w", err)
	}

	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			m.logger.Error("Error create bucket", zap.Error(err))

			return false, fmt.Errorf("%w", err)
		}
	}

	return false, nil
}
