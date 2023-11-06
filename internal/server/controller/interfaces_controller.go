package controller

import (
	"context"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
)

type SecretUseCase interface {
	MultipartUploadFile(context.Context, <-chan *models.MultipartUploadFileData, string) (bool, error)
	CreateBucketSecret(ctx context.Context, username string, userID int) (bool, error)
	MultipartDownloadFile(context.Context,
		*models.MultipartDownloadFileRequest,
		string) (chan *models.MultipartDownloadFileResponse, chan error)
	ApplyChanges(ctx context.Context, bucketName string) (bool, error)
	DiscardChanges(ctx context.Context, bucketName string) (bool, error)
}

type UserUseCase interface {
	RegisterUser(ctx context.Context, user entity.User) (string, int, error)
	LoginUser(ctx context.Context, user entity.User) (string, int, error)
}

type SyncStatus interface {
	AddSync(userID int, syncID uuid.UUID) (bool, error)
	RemoveClientSync(userID int, syncID uuid.UUID) error
	IsSyncExists(userID int, syncID uuid.UUID) (bool, error)
}
