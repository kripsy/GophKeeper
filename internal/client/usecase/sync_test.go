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
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
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
