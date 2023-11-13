//nolint:testpackage,goerr113,maintidx
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

func TestClientUsecase_updateSecret(t *testing.T) {
	file := "pathf"
	secretName := "test"
	testErr := fmt.Errorf("error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defer os.RemoveAll(file)

	log := zerolog.New(os.Stdout)

	type args struct {
		secretName string
		updateType int
		success    bool
	}
	tests := []struct {
		name    string
		usecase ClientUsecase
		args    args
	}{
		{
			name: "update basic auth",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.BasicAuthType}

				data := filemanager.BasicAuth{Login: "test", Password: "password"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddBasicAuth().Return(data, nil)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().UpdateDataByName(secretName, data).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update note",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				data := filemanager.Note{Text: "test"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(data, nil)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().UpdateDataByName(secretName, data).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update card data",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.CardDataType}

				data := filemanager.CardData{Number: "test", Date: "02/22", CVV: "123"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddCard().Return(data, nil)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().UpdateDataByName(secretName, data).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update file",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.FileType}
				err := os.WriteFile(file, []byte("test"), permissions.FileMode)
				if err != nil {
					t.Fatalf("write file err: %v", err)
				}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().GetFilePath().Return(file)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().AddFileToStorage(false, secretName, gomock.Any(), gomock.Any()).Return(nil)
				fm.EXPECT().UpdateInfoByName(secretName, gomock.Any()).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update not existing file",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.FileType}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().GetFilePath().Return("test")
				cli.EXPECT().PrintErr(ui.UpdateErr)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().AddFileToStorage(false, gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update file name err",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.FileType}
				err := os.WriteFile(file, []byte("test"), permissions.FileMode)
				if err != nil {
					t.Fatalf("write file err: %v", err)
				}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().GetFilePath().Return(file)
				cli.EXPECT().PrintErr(ui.UpdateErr)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().AddFileToStorage(false, gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update secret note err",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				data := filemanager.Note{Text: "test"}
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(data, nil)
				cli.EXPECT().PrintErr(ui.UpdateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().UpdateDataByName(secretName, data).Return(testErr)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update meta info",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddMetaInfo().Return(models.DataInfo{Description: "test"}, nil)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				fm.EXPECT().UpdateInfoByName(secretName, gomock.Any()).Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 1, success: true},
		},
		{
			name: "update meta info filemanager err",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddMetaInfo().Return(models.DataInfo{Name: "test"}, nil)
				cli.EXPECT().PrintErr(ui.UpdateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().UpdateInfoByName(secretName, gomock.Any()).Return(testErr)
				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 1, success: true},
		},
		{
			name: "update meta info ui err",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddMetaInfo().Return(models.DataInfo{}, testErr)
				cli.EXPECT().PrintErr(ui.UpdateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 1, success: true},
		},
		{
			name: "update note data ui err",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				meta[secretName] = models.DataInfo{Name: secretName, DataType: filemanager.NoteType}

				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().AddNote().Return(filemanager.Note{}, testErr)
				cli.EXPECT().PrintErr(ui.UpdateErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 0, success: true},
		},
		{
			name: "update exit",
			usecase: func() ClientUsecase {
				meta := make(models.MetaData)
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				usecase := ClientUsecase{
					userData:    &models.UserData{Meta: models.UserMeta{Data: meta}},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: secretName, updateType: 1, success: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.usecase.updateSecret(tt.args.secretName, tt.args.updateType, tt.args.success)
		})
	}
}
