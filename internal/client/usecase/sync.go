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

// sync initiates the synchronization process with the server.
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

// synchronizeWithServer synchronizes local and server data.
func (c *ClientUsecase) synchronizeWithServer(ctx context.Context, syncKey string, serverMeta *models.UserMeta) error {
	if nil == serverMeta {
		return nil
	}

	defer c.syncDeletedSecret(serverMeta.DeletedData)

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

// uploadMeta uploads the user's meta information to the server.
func (c *ClientUsecase) uploadMeta(ctx context.Context, syncKey string) error {
	data, err := c.fileManager.ReadEncryptedByID(c.userData.User.GetMetaFileName())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = c.grpc.UploadFile(ctx, c.userData.User.GetMetaFileName(), c.userData.Meta.HashData, syncKey, data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// syncDeletedSecret handles synchronization of deleted secrets.
func (c *ClientUsecase) syncDeletedSecret(deleted models.Deleted) {
	var needDeleted []string
	for dataID := range deleted {
		if _, ok := c.userData.Meta.DeletedData[dataID]; !ok {
			c.userData.Meta.DeletedData[dataID] = struct{}{}
			needDeleted = append(needDeleted, dataID)
		}
	}

	if len(needDeleted) == 0 {
		return
	}

	if err := c.fileManager.DeleteByIDs(needDeleted); err != nil {
		c.log.Err(err).Msg("failed sync deleted secret")
	}
}

// uploadSecrets uploads secrets to the server.
func (c *ClientUsecase) uploadSecrets(ctx context.Context, syncKey string, toUpload models.MetaData) error {
	var wg sync.WaitGroup

	for dataID, info := range toUpload {
		wg.Add(1)

		data, err := c.fileManager.ReadEncryptedByID(dataID)
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

// downloadSecrets downloads secrets from the server.
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

// downloadServerMeta downloads the server's metadata.
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

// blockSync blocks synchronization until a signal is received from the server.
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

// findDifferences identifies differences between local and server data.
//
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

	// Check the server data, skip the remote ones locally.
	// If the data is not detected or outdated, add it to the list for download.
	for dataID, sData := range serverData {
		if local.DeletedData.IsDeleted(dataID) {
			continue
		}
		lData, found := localData[dataID]
		if !found || lData.UpdatedAt.Before(sData.UpdatedAt) {
			needDownload[dataID] = sData
		}
	}

	// 	We check the local data, skip the remote data on the server.
	//	If the data is not detected or outdated, add it to the list for uploading
	for dataID, lData := range localData {
		if server.DeletedData.IsDeleted(dataID) {
			continue
		}

		sData, found := serverData[dataID]
		if !found || sData.UpdatedAt.Before(lData.UpdatedAt) {
			needUpload[dataID] = lData
		}
	}

	return needDownload, needUpload
}
