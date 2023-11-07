package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type minioRepository struct {
	client *minio.Client
	logger *zap.Logger
}

//nolint:revive,nolintlint
func NewMinioRepository(ctx context.Context,
	endpoint,
	accessKeyID,
	secretAccessKey,
	bucketName string,
	isUseSSL bool, l *zap.Logger) (*minioRepository, error) {
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

func (m *minioRepository) MultipartUploadFile(ctx context.Context,
	data *models.MultipartUploadFileData,
	partNum int,
	bucketName string) (*models.ObjectPart, error) {
	objectName := fmt.Sprintf("%s-part-%d.rc", data.FileName, partNum)
	opts := minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"Hash":     data.Hash,
			"Username": data.Username,
		},
	}

	_, err := m.client.PutObject(ctx,
		bucketName,
		objectName,
		bytes.NewReader(data.Content),
		int64(len(data.Content)), opts)
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

func (m *minioRepository) ListObjects(ctx context.Context, bucketName, prefix string) (*[]string, error) {
	opts := minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}

	objectCh := m.client.ListObjects(ctx, bucketName, opts)

	//nolint:prealloc
	var objectNames []string
	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing objects:", object.Err)

			return nil, object.Err
		}
		objectNames = append(objectNames, object.Key)
	}

	return &objectNames, nil
}

func (m *minioRepository) GetObject(ctx context.Context,
	bucketName, filename string) (*[]byte, string, error) {
	object, err := m.client.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		m.logger.Debug("Error in minio GetObject", zap.Error(err))

		return nil, "", fmt.Errorf("%w", err)
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		m.logger.Debug("Error reading object content", zap.Error(err))

		return nil, "", fmt.Errorf("%w", err)
	}

	s, err := object.Stat()
	if err != nil {
		m.logger.Debug("Error reading object stat", zap.Error(err))

		return nil, "", fmt.Errorf("%w", err)
	}

	hash := s.UserMetadata["Hash"]
	m.logger.Debug("Hash data", zap.String("msg", hash))

	return &content, hash, nil
}

func (m *minioRepository) ListFilesWithPostfix(bucketName, postfix string) ([]string, error) {
	var fileList []string

	// Создаем канал для получения объектов
	objectCh := m.client.ListObjects(context.Background(), bucketName, minio.ListObjectsOptions{})

	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing objects:", object.Err)

			return nil, fmt.Errorf("%w", object.Err)
		}

		// Проверяем, оканчивается ли имя файла на заданный постфикс
		if strings.HasSuffix(object.Key, postfix) {
			fileList = append(fileList, object.Key)
		}
	}

	return fileList, nil
}

//nolint:cyclop
func (m *minioRepository) ApplyChanges(ctx context.Context, bucketName string) error {
	postfix := ".rc"

	// Получить список файлов с расширением .rc
	//nolint:contextcheck
	rcFiles, err := m.ListFilesWithPostfix(bucketName, postfix)
	if err != nil {
		m.logger.Error("Error to get list files with rc prefix", zap.Error(err))

		return err
	}
	m.logger.Debug("List files with rc prefix", zap.Any("rcFiles", rcFiles))

	// Создать карту для файлов .rc
	rcMap := make(map[string]bool)
	for _, file := range rcFiles {
		// baseName := strings.TrimSuffix(file, postfix)      // Удалить .rc
		secretName := strings.Split(file, "-part-")[0] // Получить имя секрета без -part-<номер>
		rcMap[secretName] = true
	}
	m.logger.Debug("Current rcMap - files, that without rc, union by -path-", zap.Any("map", rcMap))

	// Получить список всех файлов
	//nolint:contextcheck
	allFiles, err := m.ListFilesWithPostfix(bucketName, "")
	if err != nil {
		m.logger.Error("Error to get list all files", zap.Error(err))

		return err
	}

	m.logger.Debug("List all files", zap.Any("allFiles", allFiles))

	for _, file := range allFiles {
		secretName := strings.Split(file, "-part-")[0]
		if rcMap[secretName] && !isInRCFiles(file, rcFiles) {
			// Если файл есть в карте, то удалить его
			m.logger.Debug("remove file for update", zap.String("fileName", file))
			err := m.client.RemoveObject(ctx, bucketName, file, minio.RemoveObjectOptions{})
			if err != nil {
				m.logger.Error("Error remove file, that should be updated", zap.Error(err))

				return fmt.Errorf("%w", err)
			}
		}
	}

	// Копировать файлы .rc без расширения .rc
	for _, rcFile := range rcFiles {
		newName := strings.TrimSuffix(rcFile, postfix)
		srcOpts := minio.CopySrcOptions{
			Bucket: bucketName,
			Object: rcFile,
		}

		dstOpts := minio.CopyDestOptions{
			Bucket: bucketName,
			Object: newName,
		}

		_, err := m.client.CopyObject(ctx, dstOpts, srcOpts)
		if err != nil {
			m.logger.Error("Error copy rc file", zap.Error(err))

			return fmt.Errorf("%w", err)
		}

		// Удалить исходный файл .rc после копирования
		err = m.client.RemoveObject(ctx, bucketName, rcFile, minio.RemoveObjectOptions{})
		if err != nil {
			m.logger.Error("Error remove rc file", zap.Error(err))

			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (m *minioRepository) DiscardChanges(ctx context.Context, bucketName string) error {
	postfix := ".rc"
	//nolint:contextcheck
	rcFiles, err := m.ListFilesWithPostfix(bucketName, postfix)
	if err != nil {
		return err
	}

	for _, file := range rcFiles {
		err := m.client.RemoveObject(ctx, bucketName, file, minio.RemoveObjectOptions{})
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

// Функция для проверки, находится ли файл в списке rcFiles.
func isInRCFiles(file string, rcFiles []string) bool {
	for _, rcFile := range rcFiles {
		if file == rcFile {
			return true
		}
	}

	return false
}
