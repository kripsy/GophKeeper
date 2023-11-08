//nolint:testpackage,goerr113
package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	mock_filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

func TestClientUsecase_createSecret(t *testing.T) {
	file := "filepath"
	testErr := fmt.Errorf("error")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defer os.RemoveAll(file)

	dataInfo := models.DataInfo{Name: "test", Description: "test"}
	log := zerolog.New(os.Stdout)
	type args struct {
		secretType int
		success    bool
	}
	tests := []struct {
		name    string
		usecase ClientUsecase
		args    args
	}{{
		name: "create basic auth",
		usecase: func() ClientUsecase {
			data := filemanager.BasicAuth{Login: "test", Password: "password"}
			cli := mock_ui.NewMockUserInterface(mockCtrl)
			cli.EXPECT().AddBasicAuth().Return(data, nil)
			cli.EXPECT().AddMetaInfo().Return(dataInfo, nil)
			fm := mock_filemanager.NewMockFileStorage(mockCtrl)

			testDataInfo := dataInfo
			testDataInfo.DataType = filemanager.BasicAuthType
			fm.EXPECT().AddToStorage(dataInfo.Name, data, testDataInfo).Return(nil)
			usecase := ClientUsecase{
				userData:    &models.UserData{},
				fileManager: fm,
				ui:          cli,
				log:         log,
			}

			return usecase
		}(),
		args: args{secretType: filemanager.BasicAuthType, success: true},
	},
		{
			name: "create note",
			usecase: func() ClientUsecase {
				data := filemanager.Note{Text: "test"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(data, nil)
				cli.EXPECT().AddMetaInfo().Return(dataInfo, nil)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.NoteType
				fm.EXPECT().AddToStorage(dataInfo.Name, data, testDataInfo).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.NoteType, success: true},
		},
		{
			name: "create card data",
			usecase: func() ClientUsecase {
				data := filemanager.CardData{Number: "test", Date: "02/22", CVV: "123"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddCard().Return(data, nil)
				cli.EXPECT().AddMetaInfo().Return(dataInfo, nil)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.CardDataType
				fm.EXPECT().AddToStorage(dataInfo.Name, data, testDataInfo).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.CardDataType, success: true},
		},
		{
			name: "create file",
			usecase: func() ClientUsecase {
				err := os.WriteFile(file, []byte("test"), permissions.FileMode)
				if err != nil {
					t.Fatalf("write file err: %v", err)
				}
				data := filemanager.File{Data: []byte("test")}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().GetFilePath().Return(file)
				cli.EXPECT().AddMetaInfo().Return(dataInfo, nil)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				testDataInfo.FileName = &file
				fm.EXPECT().AddToStorage(dataInfo.Name, data, testDataInfo).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.FileType, success: true},
		},
		{
			name: "create file wrong path",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().GetFilePath().Return("test")
				cli.EXPECT().PrintErr(ui.CreateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.FileType
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.FileType, success: true},
		},
		{
			name: "create note with meta info err",
			usecase: func() ClientUsecase {
				data := filemanager.Note{Text: "test"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(data, nil)
				cli.EXPECT().AddMetaInfo().Return(models.DataInfo{}, testErr)
				cli.EXPECT().PrintErr(ui.CreateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.NoteType, success: true},
		},
		{
			name: "create note with secret err",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(filemanager.Note{}, testErr)
				cli.EXPECT().PrintErr(ui.CreateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.NoteType, success: true},
		},
		{
			name: "create note filestorage err",
			usecase: func() ClientUsecase {
				data := filemanager.Note{Text: "test"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(data, nil)
				cli.EXPECT().AddMetaInfo().Return(dataInfo, nil)
				cli.EXPECT().PrintErr(ui.CreateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				testDataInfo := dataInfo
				testDataInfo.DataType = filemanager.NoteType
				fm.EXPECT().AddToStorage(dataInfo.Name, data, testDataInfo).Return(testErr)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.NoteType, success: true},
		},
		{
			name: "create exit",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretType: filemanager.NoteType, success: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.usecase.createSecret(tt.args.secretType, tt.args.success)
		})
	}
}
