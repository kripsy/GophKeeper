//nolint:testpackage
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	mock_grpc "github.com/kripsy/GophKeeper/internal/client/grpc/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	mock_filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_findDifferences(t *testing.T) {
	testTime := time.Now()
	tests := []struct {
		name             string
		local            models.UserMeta
		server           models.UserMeta
		wantNeedDownload models.MetaData
		wantNeedUpload   models.MetaData
	}{
		{
			name: "download & upload",
			local: models.UserMeta{Data: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime},
			},
			wantNeedUpload: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime},
			},
		},
		{
			name: "upload",
			local: models.UserMeta{Data: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}}},
			server:           models.UserMeta{Data: models.MetaData{}},
			wantNeedDownload: models.MetaData{},
			wantNeedUpload: models.MetaData{
				"data1": models.DataInfo{DataID: "data1", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
		},
		{
			name:  "download",
			local: models.UserMeta{Data: models.MetaData{}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}},
			wantNeedUpload: models.MetaData{},
		},
		{
			name: "upload with deleted data",
			local: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-2 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
			},
			server:           models.UserMeta{DeletedData: models.Deleted{"data2": struct{}{}}},
			wantNeedDownload: models.MetaData{},
			wantNeedUpload: models.MetaData{
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime}},
		},
		{
			name: "download with deleted",
			local: models.UserMeta{Data: models.MetaData{},
				DeletedData: models.Deleted{"data2": struct{}{}}},
			server: models.UserMeta{Data: models.MetaData{
				"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}}},
			wantNeedDownload: models.MetaData{
				"data3": models.DataInfo{DataID: "data3", UpdatedAt: testTime.Add(-1 * time.Hour)},
				"data4": models.DataInfo{DataID: "data4", UpdatedAt: testTime}},
			wantNeedUpload: models.MetaData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNeedDownload, gotNeedUpload := findDifferences(tt.local, tt.server)
			if !reflect.DeepEqual(gotNeedDownload, tt.wantNeedDownload) {
				t.Errorf("findDifferences() gotNeedDownload = %v, want %v", gotNeedDownload, tt.wantNeedDownload)
			}
			if !reflect.DeepEqual(gotNeedUpload, tt.wantNeedUpload) {
				t.Errorf("findDifferences() gotNeedUpload = %v, want %v", gotNeedUpload, tt.wantNeedUpload)
			}
		})
	}
}

var ErrEmpty = errors.New("")

func TestClientUsecaseDownloadServerMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_grpc.NewMockClient(ctrl)
	userKey := "111111111111111111111111"
	testMeta := &models.UserMeta{
		Username: "testuser",
	}
	//nolint:goconst
	syncKey := "syncKey"
	userData :=
		&models.UserData{
			User: models.User{
				Username: "testuser",
				Password: "testpassword",
				Key:      []byte(userKey),
			},
			Meta: models.UserMeta{
				Username: "testuser",
			},
		}
	key, err := userData.User.GetUserKey()
	assert.NoError(t, err)

	userData2 :=
		&models.UserData{
			User: models.User{
				Username: "testuser",
				Password: "testpassword2",
				Key:      []byte("fake"),
			},
			Meta: models.UserMeta{
				Username: "testuser",
			},
		}
	key2, err := userData2.User.GetUserKey()
	assert.NoError(t, err)
	type fields struct {
		dataPath    string
		uploadPath  string
		aboutMsg    string
		userData    *models.UserData
		grpc        grpc.Client
		fileManager filemanager.FileStorage
		ui          ui.UserInterface
		log         zerolog.Logger
	}
	type args struct {
		//nolint:containedctx
		ctx     context.Context
		syncKey string
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		want    *models.UserMeta
		wantErr bool
	}{
		{
			name: "failed in c.grpc.DownloadFile",
			args: args{
				ctx:     context.Background(),
				syncKey: syncKey,
			},
			setup: func() {
				mockClient.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(), gomock.Any(),
					syncKey).Return(make(chan []byte), ErrEmpty).Times(1)
			},
			wantErr: true,
			fields: fields{
				grpc: mockClient,
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
						Password: "testpassword",
						Key:      []byte("testkey"),
					},
					Meta: models.UserMeta{},
				},
			},
			want: nil,
		},
		{
			name: "len concatenatedData is zero",
			args: args{
				ctx:     context.Background(),
				syncKey: syncKey,
			},
			setup: func() {
				emptyChan := make(chan []byte)
				close(emptyChan)
				mockClient.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(), gomock.Any(),
					syncKey).Return(emptyChan, nil).Times(1)
			},
			wantErr: false,
			fields: fields{
				grpc: mockClient,
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
						Password: "testpassword",
						Key:      []byte("testkey"),
					},
					Meta: models.UserMeta{},
				},
			},
			want: &models.UserMeta{},
		},
		{
			name: "successful data retrieval and decoding",
			args: args{
				ctx:     context.Background(),
				syncKey: syncKey,
			},
			setup: func() {
				dataChan := make(chan []byte, 1)
				mashaldata, err := json.Marshal(testMeta)
				assert.NoError(t, err)
				encryptedData, err := utils.Encrypt(mashaldata, key)
				assert.NoError(t, err)
				dataChan <- encryptedData
				close(dataChan)
				mockClient.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(), gomock.Any(),
					syncKey).Return(dataChan, nil).Times(1)
			},
			wantErr: false,
			fields: fields{
				grpc:     mockClient,
				userData: userData,
			},
			want: testMeta,
		},
		{
			name: "failed decode data",
			args: args{
				ctx:     context.Background(),
				syncKey: syncKey,
			},
			setup: func() {
				dataChan := make(chan []byte, 1)
				mashaldata, err := json.Marshal(testMeta)
				assert.NoError(t, err)
				encryptedData, err := utils.Encrypt(mashaldata, key2)
				assert.NoError(t, err)
				dataChan <- encryptedData
				close(dataChan)
				mockClient.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(), gomock.Any(),
					syncKey).Return(dataChan, nil).Times(1)
			},
			wantErr: true,
			fields: fields{
				grpc:     mockClient,
				userData: userData,
			},
			want: nil,
		},

		{
			name: "failed unmarshall data",
			args: args{
				ctx:     context.Background(),
				syncKey: syncKey,
			},
			setup: func() {
				dataChan := make(chan []byte, 1)

				encryptedData, err := utils.Encrypt([]byte("bar"), key)
				assert.NoError(t, err)
				dataChan <- encryptedData
				close(dataChan)
				mockClient.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(), gomock.Any(),
					syncKey).Return(dataChan, nil).Times(1)
			},
			wantErr: true,
			fields: fields{
				grpc:     mockClient,
				userData: userData,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				dataPath:    tt.fields.dataPath,
				uploadPath:  tt.fields.uploadPath,
				aboutMsg:    tt.fields.aboutMsg,
				userData:    tt.fields.userData,
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				ui:          tt.fields.ui,
				log:         tt.fields.log,
			}
			tt.setup()

			got, err := c.downloadServerMeta(tt.args.ctx, tt.args.syncKey)

			if (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.downloadServerMeta() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ClientUsecase.downloadServerMeta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClientUsecaseDownloadSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_grpc.NewMockClient(ctrl)
	mockFileManager := mock_filemanager.NewMockFileStorage(ctrl)

	ctx := context.Background()
	syncKey := "syncKey"
	testTime := time.Now()
	testMetaData := models.MetaData{
		"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
	}
	type fields struct {
		grpc        grpc.Client
		fileManager filemanager.FileStorage
		// other fields
	}
	type args struct {
		//nolint:containedctx
		ctx        context.Context
		syncKey    string
		toDownload models.MetaData
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		wantErr bool
	}{
		{
			name: "successful download and storage",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func() {
				dataChan := make(chan []byte, 1)
				data, err := json.Marshal(testMetaData)
				assert.NoError(t, err)
				dataChan <- data
				close(dataChan)
				mockClient.EXPECT().DownloadFile(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any()).Return(dataChan, nil).Times(1)
				mockFileManager.EXPECT().AddEncryptedToStorage(gomock.Any(),
					gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			fields: fields{
				grpc:        mockClient,
				fileManager: mockFileManager,
			},
			wantErr: false,
		},
		{
			name: "unsuccessful download and storage",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func() {
				for _, info := range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mockClient.EXPECT().DownloadFile(gomock.Any(),
						info.DataID, gomock.Any(),
						syncKey).Return(dataChan, ErrEmpty).Times(1)
				}
			},
			fields: fields{
				grpc:        mockClient,
				fileManager: mockFileManager,
			},
			wantErr: true,
		},
		{
			name: "unsuccessful storage",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func() {
				for _, info := range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mockClient.EXPECT().DownloadFile(gomock.Any(), info.DataID, gomock.Any(), syncKey).Return(dataChan, nil).Times(1)
					mockFileManager.EXPECT().AddEncryptedToStorage(info.Name, gomock.Any(), info).Return(ErrEmpty).Times(1)
				}
			},
			fields: fields{
				grpc:        mockClient,
				fileManager: mockFileManager,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				userData: &models.UserData{
					Meta: models.UserMeta{
						HashData: "testhash",
					},
				},
			}
			tt.setup()
			err := c.downloadSecrets(tt.args.ctx, tt.args.syncKey, tt.args.toDownload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.downloadSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUsecaseBlockSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mock_grpc.NewMockClient(ctrl)
	mockFileManager := mock_filemanager.NewMockFileStorage(ctrl)

	ctx := context.Background()
	syncKey := "syncKey"
	testTime := time.Now()
	testMetaData := models.MetaData{
		"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
	}
	type fields struct {
		grpc        grpc.Client
		fileManager filemanager.FileStorage
		// other fields
	}
	type args struct {
		//nolint:containedctx
		ctx        context.Context
		syncKey    string
		toDownload models.MetaData
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func()
		args    args
		wantErr bool
	}{
		{
			name: "GUID not match",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func() {
				guidChan := make(chan interface{}, 1)
				mockClient.EXPECT().BlockStore(gomock.Any(),
					gomock.Any(), gomock.Any()).Do(func(ctx context.Context,
					key string, ch chan string) {
					ch <- "differentSyncKey"
					close(guidChan)
				})
			},
			fields: fields{
				grpc:        mockClient,
				fileManager: mockFileManager,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				userData: &models.UserData{
					Meta: models.UserMeta{
						HashData: "testhash",
					},
				},
			}
			tt.setup()
			err := c.blockSync(tt.args.ctx, tt.args.syncKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.downloadSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUsecaseUploadSecrets(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	syncKey := "syncKey"
	testTime := time.Now()
	testMetaData := models.MetaData{
		"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
	}
	type fields struct {
		grpc        grpc.Client
		fileManager filemanager.FileStorage
		// other fields
	}
	type args struct {
		//nolint:containedctx
		ctx        context.Context
		syncKey    string
		toDownload models.MetaData
	}
	//nolint:dupl
	tests := []struct {
		name    string
		fields  fields
		setup   func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient)
		args    args
		wantErr bool
	}{
		{
			name: "error in c.fileManager.ReadEncryptedByID",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(nil, ErrEmpty).AnyTimes()
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: true,
		},
		{
			name: "error in c.grpc.UploadFile",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan, nil).AnyTimes()
					mc.EXPECT().UploadFile(gomock.Any(),
						gomock.Any(), gomock.Any(),
						gomock.Any(), dataChan).Return(ErrEmpty).MinTimes(1)
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan, nil).AnyTimes()
					mc.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dataChan).Return(nil).MinTimes(1)
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				userData: &models.UserData{
					Meta: models.UserMeta{
						HashData: "testhash",
					},
				},
			}
			mockClient := mock_grpc.NewMockClient(ctrl)
			mockFileManager := mock_filemanager.NewMockFileStorage(ctrl)
			c.grpc = mockClient
			c.fileManager = mockFileManager
			tt.setup(mockFileManager, mockClient)
			err := c.uploadSecrets(tt.args.ctx, tt.args.syncKey, tt.args.toDownload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.downloadSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUsecaseSyncDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userDeletedTest := make(map[string]struct{})
	userDeletedTest["1"] = struct{}{}
	userDeletedTest["2"] = struct{}{}

	type fields struct {
		grpc        grpc.Client
		fileManager filemanager.FileStorage
	}
	type args struct {
		//nolint:containedctx,unused
		ctx context.Context
		//nolint:unused
		syncKey string
		//nolint:unused
		toDownload models.MetaData
		deleted    models.Deleted
	}
	tests := []struct {
		name    string
		fields  fields
		setup   func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient)
		args    args
		wantErr bool
	}{
		{
			name: "len deleted 0",
			args: args{
				deleted: models.Deleted{},
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {

			},
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "len deleted not 0",
			args: args{
				//nolint:gofmt
				deleted: map[string]struct{}{
					"1": struct{}{},
					"2": struct{}{},
				},
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {

			},
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "len deleted not 0 and more than in user and success delete",
			args: args{
				deleted: map[string]struct{}{
					//nolint:gofmt
					"1": struct{}{},
					"2": struct{}{},
					"3": struct{}{},
				},
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				mfm.EXPECT().DeleteByIDs(gomock.Any()).Return(nil).MinTimes(1)
			},
			fields:  fields{},
			wantErr: false,
		},
		{
			name: "len deleted not 0 and more than in user and failed delete",
			args: args{
				deleted: map[string]struct{}{
					//nolint:gofmt
					"1": struct{}{},
					"2": struct{}{},
					"3": struct{}{},
				},
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				mfm.EXPECT().DeleteByIDs(gomock.Any()).Return(ErrEmpty).MinTimes(1)
			},
			fields:  fields{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				userData: &models.UserData{
					Meta: models.UserMeta{
						HashData: "testhash",
						DeletedData: map[string]struct{}{
							//nolint:gofmt
							"1": struct{}{},
							"2": struct{}{},
						},
					},
				},
			}
			mockClient := mock_grpc.NewMockClient(ctrl)
			mockFileManager := mock_filemanager.NewMockFileStorage(ctrl)
			c.grpc = mockClient
			c.fileManager = mockFileManager
			tt.setup(mockFileManager, mockClient)
			c.syncDeletedSecret(tt.args.deleted)
		})
	}
}

func TestClientUsecaseUploadMeta(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	syncKey := "syncKey"
	testTime := time.Now()
	testMetaData := models.MetaData{
		"data2": models.DataInfo{DataID: "data2", UpdatedAt: testTime.Add(-1 * time.Hour)},
	}
	type fields struct {
		grpc        grpc.Client
		fileManager filemanager.FileStorage
		// other fields
	}
	type args struct {
		//nolint:containedctx
		ctx        context.Context
		syncKey    string
		toDownload models.MetaData
	}
	//nolint:dupl
	tests := []struct {
		name    string
		fields  fields
		setup   func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient)
		args    args
		wantErr bool
	}{
		{
			name: "error in c.fileManager.ReadEncryptedByID",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(nil, ErrEmpty).AnyTimes()
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: true,
		},
		{
			name: "error in c.grpc.UploadFile",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan, nil).AnyTimes()
					mc.EXPECT().UploadFile(gomock.Any(),
						gomock.Any(), gomock.Any(),
						gomock.Any(), dataChan).Return(ErrEmpty).MinTimes(1)
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:        ctx,
				syncKey:    syncKey,
				toDownload: testMetaData,
			},
			setup: func(mfm *mock_filemanager.MockFileStorage, mc *mock_grpc.MockClient) {
				//nolint:gofmt,revive,gosimple
				for _, _ = range testMetaData {
					dataChan := make(chan []byte, 1)
					data, err := json.Marshal(testMetaData)
					assert.NoError(t, err)
					dataChan <- data
					close(dataChan)
					mfm.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan, nil).AnyTimes()
					mc.EXPECT().UploadFile(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), dataChan).Return(nil).MinTimes(1)
				}
			},
			fields: fields{
				// grpc:        mockClient,
				// fileManager: mockFileManager,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientUsecase{
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				userData: &models.UserData{
					Meta: models.UserMeta{
						HashData: "testhash",
					},
				},
			}
			mockClient := mock_grpc.NewMockClient(ctrl)
			mockFileManager := mock_filemanager.NewMockFileStorage(ctrl)
			c.grpc = mockClient
			c.fileManager = mockFileManager
			tt.setup(mockFileManager, mockClient)
			err := c.uploadMeta(tt.args.ctx, tt.args.syncKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.downloadSecrets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUsecaseSynchronizeWithServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	type fields struct {
		dataPath   string
		uploadPath string
		aboutMsg   string
		userData   *models.UserData
		//nolint:unused
		grpc grpc.Client
		//nolint:unused
		fileManager filemanager.FileStorage
		//nolint:unused
		ui  ui.UserInterface
		log zerolog.Logger
	}
	type args struct {
		//nolint:containedctx
		ctx        context.Context
		syncKey    string
		serverMeta *models.UserMeta
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		setup   func(mui *mock_ui.MockUserInterface,
			grpc *mock_grpc.MockClient,
			mockFileManager *mock_filemanager.MockFileStorage)
	}{
		{
			name: "Empty server Meta",
			args: args{
				serverMeta: nil,
			},
			wantErr: false,
			setup: func(mui *mock_ui.MockUserInterface,
				grpc *mock_grpc.MockClient,
				mockFileManager *mock_filemanager.MockFileStorage) {
			},
		},
		{
			name: "Zero different",
			args: args{
				serverMeta: &models.UserMeta{
					DeletedData: models.Deleted{},
					Data:        models.MetaData{},
				},
			},
			fields: fields{
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
					},
					Meta: models.UserMeta{
						DeletedData: models.Deleted{},
						Data:        models.MetaData{},
					},
				},
			},
			setup: func(mui *mock_ui.MockUserInterface,
				grpc *mock_grpc.MockClient,
				mockFileManager *mock_filemanager.MockFileStorage) {
			},
			wantErr: false,
		},
		{
			name: "Error downloadSecrets from server",
			args: args{
				serverMeta: &models.UserMeta{
					DeletedData: models.Deleted{},
					Data: models.MetaData{
						"b": models.DataInfo{
							Name:   "file2",
							DataID: "fileID2",
							Hash:   "hash2",
						},
					},
				},
			},
			fields: fields{
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
					},
					Meta: models.UserMeta{
						DeletedData: models.Deleted{},
						Data: models.MetaData{
							"a": models.DataInfo{
								Name:   "file",
								DataID: "fileID",
								Hash:   "hash",
							},
						},
					},
				},
			},
			setup: func(mui *mock_ui.MockUserInterface,
				grpc *mock_grpc.MockClient,
				mockFileManager *mock_filemanager.MockFileStorage) {
				grpc.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(nil, ErrEmpty).Times(1)
				mui.EXPECT().PrintErr(gomock.Any()).Times(1)
			},
			wantErr: true,
		},
		{
			name: "Error uploadSecrets from server",
			args: args{
				serverMeta: &models.UserMeta{
					DeletedData: models.Deleted{},
					Data:        models.MetaData{
						// "b": models.DataInfo{
						// 	Name:   "file2",
						// 	DataID: "fileID2",
						// 	Hash:   "hash2",
						// },
					},
				},
			},
			fields: fields{
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
					},
					Meta: models.UserMeta{
						DeletedData: models.Deleted{},
						Data: models.MetaData{
							"a": models.DataInfo{
								Name:   "file",
								DataID: "fileID",
								Hash:   "hash",
							},
							// "b": models.DataInfo{
							// 	Name:   "file2",
							// 	DataID: "fileID2",
							// 	Hash:   "hash2",
							// },
						},
					},
				},
			},
			setup: func(mui *mock_ui.MockUserInterface,
				grpc *mock_grpc.MockClient,
				mockFileManager *mock_filemanager.MockFileStorage) {
				dataChan := make(chan []byte, 1)
				close(dataChan)
				grpc.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(dataChan, nil).AnyTimes()
				mui.EXPECT().PrintErr(gomock.Any()).AnyTimes()
				mockFileManager.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan,
					ErrEmpty).AnyTimes()
			},
			wantErr: true,
		},
		{
			name: "Success",
			args: args{
				serverMeta: &models.UserMeta{
					DeletedData: models.Deleted{},
					Data:        models.MetaData{},
				},
			},
			fields: fields{
				userData: &models.UserData{
					User: models.User{
						Username: "testuser",
					},
					Meta: models.UserMeta{
						DeletedData: models.Deleted{},
						Data: models.MetaData{
							"a": models.DataInfo{
								Name:   "file",
								DataID: "fileID",
								Hash:   "hash",
							},
						},
					},
				},
			},
			setup: func(mui *mock_ui.MockUserInterface,
				grpc *mock_grpc.MockClient,
				mockFileManager *mock_filemanager.MockFileStorage) {
				dataChan := make(chan []byte, 1)
				close(dataChan)
				grpc.EXPECT().DownloadFile(gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any()).Return(dataChan, nil).AnyTimes()
				mui.EXPECT().PrintErr(gomock.Any()).AnyTimes()
				mockFileManager.EXPECT().ReadEncryptedByID(gomock.Any()).Return(dataChan,
					nil).AnyTimes()
				grpc.EXPECT().UploadFile(gomock.Any(), gomock.Any(),
					gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUI := mock_ui.NewMockUserInterface(mockCtrl)
			grpc := mock_grpc.NewMockClient(mockCtrl)
			mockFileManager := mock_filemanager.NewMockFileStorage(mockCtrl)
			tt.setup(mockUI, grpc, mockFileManager)
			c := &ClientUsecase{
				dataPath:    tt.fields.dataPath,
				uploadPath:  tt.fields.uploadPath,
				aboutMsg:    tt.fields.aboutMsg,
				userData:    tt.fields.userData,
				grpc:        grpc,
				fileManager: mockFileManager,
				ui:          mockUI,
				log:         tt.fields.log,
			}

			if err := c.synchronizeWithServer(tt.args.ctx,
				tt.args.syncKey,
				tt.args.serverMeta); (err != nil) != tt.wantErr {
				t.Errorf("ClientUsecase.synchronizeWithServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
