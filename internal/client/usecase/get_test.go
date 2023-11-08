//nolint:testpackage,goerr113,maintidx
package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	mock_filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

func TestClientUsecase_getSecrets(t *testing.T) {
	file := "filePath"
	testErr := fmt.Errorf("error")
	secretName := "test"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defer os.RemoveAll(file)

	dataInfo := models.DataInfo{Name: "test", Description: "test"}
	log := zerolog.New(os.Stdout)

	type args struct {
		secretName string
		success    bool
	}
	tests := []struct {
		name    string
		usecase ClientUsecase
		args    args
	}{
		{
			name: "get basic auth",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.BasicAuthType
				data := filemanager.BasicAuth{Login: "test", Password: "password"}
				body, err := json.Marshal(data)
				if err != nil {
					t.Fatalf("marshal data err: %v", err)
				}

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(body, testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get note",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.NoteType
				data := filemanager.Note{Text: "test"}
				body, err := json.Marshal(data)
				if err != nil {
					t.Fatalf("marshal data err: %v", err)
				}

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(body, testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get card data",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.CardDataType
				data := filemanager.CardData{Number: "test", Date: "02/22", CVV: "123"}
				body, err := json.Marshal(data)
				if err != nil {
					t.Fatalf("marshal data err: %v", err)
				}

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(body, testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get file",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				testDataInfo.FileName = &file
				data := filemanager.File{Data: []byte("test")}
				body, err := json.Marshal(data)
				if err != nil {
					t.Fatalf("marshal data err: %v", err)
				}

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(body, testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().UploadFileTo("").Return(file, true)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get file path err",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				wrongPath := file + "/wrong"
				testDataInfo.FileName = &wrongPath
				data := filemanager.File{Data: []byte("test")}
				body, err := json.Marshal(data)
				if err != nil {
					t.Fatalf("marshal data err: %v", err)
				}

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(body, testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().UploadFileTo("").Return(file, true)
				cli.EXPECT().PrintErr(ui.GetErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get file invalid body err",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				testDataInfo.FileName = &file

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return([]byte("invalid body"), testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().UploadFileTo("").Return(file, true)
				cli.EXPECT().PrintErr(ui.GetErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get file invalid body err",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				testDataInfo.FileName = &file

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return([]byte(""), testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().UploadFileTo("").Return(file, false)
				cli.EXPECT().PrintErr(ui.GetErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get note with invalid body err",
			usecase: func() ClientUsecase {
				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.NoteType

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return([]byte("invalid body"), testDataInfo, nil)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().PrintErr(ui.GetErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get note with filemanager err",
			usecase: func() ClientUsecase {
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().GetByName(secretName).Return(nil, models.DataInfo{}, testErr)

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().PrintErr(ui.GetErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: true},
		},
		{
			name: "get exit",
			usecase: func() ClientUsecase {
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				cli := mock_ui.NewMockUserInterface(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, success: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.usecase.getSecrets(tt.args.secretName, tt.args.success)
		})
	}
}
