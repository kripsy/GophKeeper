package usecase

import (
	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"time"
)

func (c *ClientUsecase) sync() {
	done := make(chan struct{})
	defer close(done)
	defer c.InMenu()

	newUUID, err := uuid.NewUUID()
	if err != nil {
		c.log.Error().Err(err).Msg("failed get uuid")
	}
	syncKey := newUUID.String()
	go c.startSync(syncKey, done)

	time.Sleep(time.Second * 4)
	c.finishSync(syncKey, done)
}

func (c *ClientUsecase) startSync(syncKey string, done <-chan struct{}) {
	go c.ui.Sync(done)
	for {
		select {
		case <-done:
			return
		default:
			_ = c.userData.User.Token
			//c.grpc.BlockSync(secretKey)
			time.Sleep(time.Millisecond * 700)
		}
	}
}

func (c *ClientUsecase) finishSync(syncKey string, done chan struct{}) {
	defer close(done)
	done <- struct{}{}

}

//func (c *ClientUsecase) SyncSecrets(user models.UserData) {
//	oldHashMeta := user.Meta.HashData
//
//	serverMetaData, err := g.SecretComparison(oldHashMeta)
//	if errors.Is(err, errors.New("CONFLICT")) {
//		fmt.Print("Another client is synchronizing, try later")
//
//		return
//	}
//	if serverMetaData == nil {
//		fmt.Print("Secret synchronized")
//
//		return
//	}
//
//	key, err := user.User.GetUserKey()
//	if err != nil {
//		//
//		return
//	}
//	data, err := utils.Decrypt(serverMetaData, key)
//	if err != nil {
//		//
//		return
//	}
//	var serverMeta models.UserMeta
//	if err := json.Unmarshal(data, &serverMeta); err != nil {
//		//
//		return
//	}
//
//	go func() { //отвельная функция
//		for {
//			g.SyncProcess(user.Meta.HashData)
//			time.Sleep(time.Second * 5)
//		}
//	}()
//
//	download, upload := FindDifferences(user.Meta.Data, serverMeta.Data)
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
//
//	g.SyncCompleted(oldHashMeta, user.Meta.HashData) // на каждое добавление скачанного секрета обновляется хэш мета
//	// internal/client/infrastrucrure/filemanager/filemanager.go:166
//}

func FindDifferences(local, server models.MetaData) (needDownload, needUpload models.MetaData) {
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
