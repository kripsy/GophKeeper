package usecase

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"time"
)

func (c *ClientUsecase) sync() {
	not := c.grpc.IsNotAvailable()
	try := c.grpc.TryToConnect()
	if not && try {
		fmt.Println("Failed connect to server")
		return
	}

	//	oldHashMeta := user.Meta.HashData

	done := make(chan struct{})
	errSync := make(chan struct{})
	defer close(done)
	defer close(errSync)
	defer c.InMenu()

	newUUID, err := uuid.NewUUID()
	if err != nil {
		c.log.Error().Err(err).Msg("failed get uuid")
	}

	syncKey := newUUID.String()

	go c.blockSync(syncKey, done, errSync)
	serverMeta, err := c.downloadServerMeta()
	toDownload, toUpload := findDifferences(c.userData.Meta.Data, serverMeta)
	_, _ = toUpload, toDownload
	//	for id, meta := range download {
	//		data, err = g.DownloadSecret(id, user.Meta.HashData)
	//		if err != nil {
	//
	//		}
	//		//err:=filemanager.AddEncrypedData(meta.Name, data, meta)
	//		_, _ = data, meta // чтобы не краснел
	//	}
	//	for id, _ := range upload {
	//		//data,err:= filemanager.GetDataByID(id)
	//		g.UploadSecret(id, []byte("data"), [32]byte{}, user.Meta.HashData)
	//	}

	time.Sleep(time.Second * 4)
	errSync <- struct{}{}
	//c.finishSync(done)
}

func (c *ClientUsecase) downloadServerMeta() (models.MetaData, error) {
	//c.grpc.Download(c.userData.User.GetMetaFileName(),)
	key, err := c.userData.User.GetUserKey()
	if err != nil {

		return nil, err
	}
	data, err := utils.Decrypt(nil, key)
	if err != nil {

		return nil, err
	}
	var serverData models.UserData
	if err := json.Unmarshal(data, &serverData); err != nil {

		return nil, err
	}

	return serverData.Meta.Data, nil
}

func (c *ClientUsecase) blockSync(syncKey string, done chan struct{}, err <-chan struct{}) {
	go c.ui.Sync(done)
	_ = c.userData.User.Token

	for {
		select {
		case <-done:
			//c.grpc.BlockSync(secretKey)
			fmt.Println("Success sync")
			return
		case <-err:
			done <- struct{}{}
			fmt.Println("Failed sync")
			return
		case <-time.Tick(time.Second):
			//c.grpc.BlockSync(secretKey)
		}
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
