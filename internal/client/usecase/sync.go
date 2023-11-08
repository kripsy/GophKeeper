//nolint:nonamedreturns,durationcheck,nolintlint
package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
)

func (c *ClientUsecase) sync() {
	defer c.InMenu()
	if c.grpc.IsNotAvailable() || !c.grpc.TryToConnect() {
		c.ui.PrintErr("Failed connect to server")

		return
	}

	done := make(chan struct{})
	errSync := make(chan struct{})
	defer close(errSync)
	defer close(done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go c.ui.Sync(done)
	syncKey := uuid.New().String()

	if err := c.blockSync(ctx, syncKey); err != nil {
		return
	}

	serverMeta, err := c.downloadServerMeta(ctx, syncKey)
	if err != nil {
		c.ui.PrintErr(ui.SyncErr)
		c.log.Err(err).Msg("failed download server meta data")

		return
	}

	if err = c.synchronizeWithServer(ctx, syncKey, serverMeta); err != nil {
		return
	}

	if err = c.uploadMeta(ctx, syncKey); err != nil {
		c.ui.PrintErr(ui.SyncErr)
		c.log.Err(err).Msg("error upload meta")

		return
	}

	if err := c.grpc.ApplyChanges(ctx, syncKey); err != nil {
		c.ui.PrintErr(ui.SyncErr)
		c.log.Err(err).Msg("failed apply changes")
	}
}

func (c *ClientUsecase) synchronizeWithServer(ctx context.Context, syncKey string, serverMeta *models.UserMeta) error {
	if nil == serverMeta {
		return nil
	}

	toDownload, toUpload := findDifferences(c.userData.Meta, *serverMeta)
	if len(toUpload) == 0 && len(toDownload) == 0 {
		return nil
	}

	if err := c.downloadSecrets(ctx, syncKey, toDownload); err != nil {
		c.ui.PrintErr(ui.SyncErr)
		c.log.Err(err).Msg("error upload secrets")

		return err
	}

	if err := c.uploadSecrets(ctx, syncKey, toUpload); err != nil {
		c.ui.PrintErr(ui.SyncErr)
		c.log.Err(err).Msg("error upload secrets")

		return err
	}

	return nil
}

func (c *ClientUsecase) uploadMeta(ctx context.Context, syncKey string) error {
	data, err := c.fileManager.ReadEncryptedByName(c.userData.User.GetMetaFileName())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = c.grpc.UploadFile(ctx, c.userData.User.GetMetaFileName(), c.userData.Meta.HashData, syncKey, data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (c *ClientUsecase) syncDeletedSecret(deleted models.Deleted) {
	//	var needDeleted []string
	for dataID := range deleted {
		if _, ok := c.userData.Meta.DeletedData[dataID]; !ok {
			c.userData.Meta.DeletedData[dataID] = struct{}{}
			//	needDeleted = append(needDeleted, dataID)
		}
	}

}

func (c *ClientUsecase) uploadSecrets(ctx context.Context, syncKey string, toUpload models.MetaData) error {
	var wg sync.WaitGroup

	for dataID, info := range toUpload {
		wg.Add(1)

		data, err := c.fileManager.ReadEncryptedByName(dataID)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		err = c.grpc.UploadFile(ctx, dataID, info.Hash, syncKey, data)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		wg.Done()
	}

	wg.Wait()

	return nil
}

func (c *ClientUsecase) downloadSecrets(ctx context.Context, syncKey string, toDownload models.MetaData) error {
	var wg sync.WaitGroup
	errors := make(chan error, 1)

	for _, info := range toDownload {
		wg.Add(1)

		go func(info models.DataInfo) {
			defer wg.Done()

			data, err := c.grpc.DownloadFile(ctx, info.DataID, c.userData.Meta.HashData, syncKey)
			if err != nil {
				errors <- err
				c.log.Err(err).Msg("failed download secret")

				return
			}

			err = c.fileManager.AddEncryptedToStorage(info.Name, data, info)
			if err != nil {
				errors <- err
				c.log.Err(err).Msg("AddEncryptedToStorage")

				return
			}
		}(info)
	}

	wg.Wait()
	if len(errors) != 0 {
		return <-errors
	}

	return nil
}

func (c *ClientUsecase) downloadServerMeta(ctx context.Context, syncKey string) (*models.UserMeta, error) {
	data, err := c.grpc.DownloadFile(ctx, c.userData.User.GetMetaFileName(), c.userData.Meta.HashData, syncKey)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	var concatenatedData []byte
	var serverData models.UserMeta
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for chunk := range data {
			concatenatedData = append(concatenatedData, chunk...)
		}
	}()
	wg.Wait()

	if len(concatenatedData) == 0 {
		return &serverData, nil
	}

	key, err := c.userData.User.GetUserKey()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	metaData, err := utils.Decrypt(concatenatedData, key)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if err = json.Unmarshal(metaData, &serverData); err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return &serverData, nil
}

func (c *ClientUsecase) blockSync(ctx context.Context, syncKey string) error {
	guidChan := make(chan string, 1)
	errChan := make(chan error, 1)

	go func() {
		err := c.grpc.BlockStore(ctx, syncKey, guidChan)
		if err != nil {
			c.log.Err(err).Msg("error block sync")
			errChan <- err

			return
		}
	}()

	select {
	case newGUID := <-guidChan:
		if syncKey != newGUID {
			c.log.Warn().Msg("sync key not equal server request key")
		}

		break
	case err := <-errChan:
		c.log.Err(err).Msg("error block sync")

		return err
	}

	return nil
}

//nolint:cyclop
func findDifferences(local, server models.UserMeta) (needDownload, needUpload models.MetaData) {
	needDownload = make(models.MetaData)
	needUpload = make(models.MetaData)

	localData := make(map[string]models.DataInfo)
	serverData := make(map[string]models.DataInfo)

	for _, data := range local.Data {
		localData[data.DataID] = data
	}
	for _, data := range server.Data {
		serverData[data.DataID] = data
	}

	// проверяем данные сервера, если данные не обнаружены или устарели добавляем в список на скачивание
	for dataID, sData := range serverData {
		lData, found := localData[dataID]
		if !found || lData.UpdatedAt.Before(sData.UpdatedAt) {
			if _, deleted := local.DeletedData[dataID]; !deleted {
				needDownload[dataID] = sData
			}
		}
	}

	// проверяем локальные данные, если данные не обнаружены или устарели добавляем в список на выгрузку
	for dataID, lData := range localData {
		sData, found := serverData[dataID]
		if !found || sData.UpdatedAt.Before(lData.UpdatedAt) {
			if _, deleted := server.DeletedData[dataID]; !deleted {
				needUpload[dataID] = lData
			}
		}
	}

	return needDownload, needUpload
}
