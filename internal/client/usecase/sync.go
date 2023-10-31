package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"sync"
)

func (c *ClientUsecase) sync() {
	defer c.InMenu()
	if c.grpc.IsNotAvailable() || !c.grpc.TryToConnect() {
		fmt.Println("Failed connect to server")
		return
	}

	//	oldHashMeta := user.Meta.HashData

	done := make(chan struct{})
	errSync := make(chan struct{})
	defer close(errSync)
	defer close(done)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = cancel
	syncKey := uuid.New().String()

	c.blockSync(ctx, syncKey, done)

	serverMeta, err := c.downloadServerMeta(ctx, syncKey)
	if err != nil {
		c.log.Err(err).Msg("failed download server meta data")

		return
	}
	toDownload, toUpload := findDifferences(c.userData.Meta.Data, serverMeta)
	_, _ = toUpload, toDownload

	//downloadSecrets

	if err = c.uploadSecrets(ctx, syncKey, toUpload); err != nil {
		c.log.Err(err).Msg("error upload secrets")
		return
	}

	c.uploadMeta(ctx, syncKey)

	if err := c.grpc.ApplyChanges(ctx, syncKey); err != nil {
		c.log.Err(err).Msg("failed apply changes")

		return
	}
	//time.Sleep(time.Second * 4)
	//errSync <- struct{}{}
	//c.finishSync(done)
}

func (c *ClientUsecase) uploadMeta(ctx context.Context, syncKey string) {
	data := make(chan []byte)
	done := make(chan struct{})

	go func() {
		//	defer wg.Done()
		err := c.fileManager.ReadEncryptedByName(c.userData.User.GetMetaFileName(), data, done)
		if err != nil {
			c.log.Err(err).Msg("")

			//errors <- err
			return
		}
	}()
	err := c.grpc.UploadFile(ctx, c.userData.User.GetMetaFileName(), c.userData.Meta.HashData, syncKey, data, done)
	if err != nil {
		//errors <- err
		c.log.Err(err).Msg("")
		return
	}

}

func (c *ClientUsecase) uploadSecrets(ctx context.Context, syncKey string, toUpload models.MetaData) error {
	var wg sync.WaitGroup
	data := make(chan []byte)
	done := make(chan struct{})
	errors := make(chan error, 1)

	for dataID, info := range toUpload {
		wg.Add(1)

		go func() {
			err := c.fileManager.ReadEncryptedByName(dataID, data, done)
			if err != nil {
				errors <- err
				return
			}
		}()

		err := c.grpc.UploadFile(ctx, dataID, info.Hash, syncKey, data, done)
		if err != nil {
			errors <- err
			return err
		}
		wg.Done()
	}

	if len(errors) != 0 {
		return <-errors
	}

	wg.Wait()
	//for err := range errors {
	//	return err
	//}

	return nil
}

func (c *ClientUsecase) downloadServerMeta(ctx context.Context, syncKey string) (models.MetaData, error) {
	//err := c.grpc.DownloadFile(ctx, c.userData.User.GetMetaFileName(), "", syncKey, nil)
	//if err != nil {
	//	return nil, err
	//}

	//sd := make(models.MetaData)
	//	return sd, nil
	data := make(chan []byte, 1) // буферизированный канал, чтобы избежать блокирования
	done := make(chan struct{}, 1)
	//close(data)

	go func() {
		defer func(done chan struct{}) {
			done <- struct{}{}
		}(done)
		err := c.grpc.DownloadFile(ctx, c.userData.User.GetMetaFileName(), c.userData.Meta.HashData, syncKey, data, done)
		if err != nil {
			close(data)
			return
		}
	}()

	err := c.fileManager.AddEncryptedToStorage("test", data, models.DataInfo{DataID: "hz"})
	if err != nil {
		c.log.Err(err).Msg("AddEncryptedToStorage")
	}

	var concatenatedData []byte
	//loop:
	//	for {
	//		select {
	//		case chunk := <-data:
	//			concatenatedData = append(concatenatedData, chunk...)
	//		case <-done:
	//
	//			break loop
	//		}
	//	}

	//close(data) //

	key, err := c.userData.User.GetUserKey()
	if err != nil {

		return nil, err
	}
	metaData, err := utils.Decrypt(concatenatedData, key)
	if err != nil {
		return nil, err
	}
	var serverData models.UserData
	if err := json.Unmarshal(metaData, &serverData); err != nil {
		return nil, err
	}

	return serverData.Meta.Data, nil
	return nil, nil
}

func (c *ClientUsecase) blockSync(ctx context.Context, syncKey string, done chan struct{}) {
	go c.ui.Sync(done)
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
	case newGuid := <-guidChan:
		_ = newGuid
		//if syncKey != newGuid {
		//	errChan <- errors.New("uuid not equal")
		//}
		break
	case err := <-errChan:
		c.log.Err(err).Msg("error block sync")
	}
}

func (c *ClientUsecase) finishSync(done chan struct{}) {
	done <- struct{}{}
}

func findDifferences(local, server models.MetaData) (needDownload, needUpload models.MetaData) {
	needDownload = make(models.MetaData)
	needUpload = make(models.MetaData)

	localData := make(map[string]models.DataInfo)
	serverData := make(map[string]models.DataInfo)

	for _, data := range local {
		localData[data.DataID] = data
	}
	for _, data := range server {
		serverData[data.DataID] = data
	}

	// проверяем данные сервера, если данные не обнаружены или устарели добавляем в список на скачивание
	for dataID, sData := range serverData {
		lData, found := localData[dataID]
		if !found || lData.UpdatedAt.Before(sData.UpdatedAt) {
			needDownload[dataID] = sData
		}
	}

	// проверяем локальные данные, если данные не обнаружены или устарели добавляем в список на выгрузку
	for dataID, lData := range localData {
		sData, found := serverData[dataID]
		if !found || sData.UpdatedAt.Before(lData.UpdatedAt) {
			needUpload[dataID] = lData
		}
	}

	return needDownload, needUpload
}
