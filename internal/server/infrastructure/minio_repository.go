package infrastructure

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinioRepository interface {
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
