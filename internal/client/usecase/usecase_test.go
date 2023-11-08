//nolint:testpackage,goerr113,maintidx
package usecase

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	mock_grpc "github.com/kripsy/GophKeeper/internal/client/grpc/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	mock_filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

func TestClientUsecase_SetUser(t *testing.T) {
	file := "filepath"
	testErr := fmt.Errorf("error")

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	defer os.RemoveAll(file)
	log := zerolog.New(os.Stdout)

	tests := []struct {
		name    string
		usecase ClientUsecase
		wantErr bool
	}{
		{
			name: "register online",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)
				cli.EXPECT().TryAgain().Return(false)
				grpc.EXPECT().IsNotAvailable().Return(false)
				cli.EXPECT().GetRepeatedPassword().Return("password", nil)
				grpc.EXPECT().IsAvailable().Return(true)
				cli.EXPECT().IsLocalStorage().Return(true)
				grpc.EXPECT().Register("username", gomock.Any()).Return(nil)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "register online with err",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username0", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)
				cli.EXPECT().TryAgain().Return(false)
				grpc.EXPECT().IsNotAvailable().Return(false)
				cli.EXPECT().GetRepeatedPassword().Return("password", nil)
				grpc.EXPECT().IsAvailable().Return(true)
				cli.EXPECT().IsLocalStorage().Return(true)
				grpc.EXPECT().Register("username0", gomock.Any()).Return(testErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			wantErr: true,
		},
		{
			name: "register local",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username1", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)
				cli.EXPECT().TryAgain().Return(false)
				grpc.EXPECT().IsNotAvailable().Return(true)
				cli.EXPECT().PrintErr("Could not connect to the server, only local registration is available")
				cli.EXPECT().GetRepeatedPassword().Return("password", nil)
				grpc.EXPECT().IsAvailable().Return(false)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "register local with wrong repeated pass",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username2", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)
				cli.EXPECT().TryAgain().Return(false)
				grpc.EXPECT().IsNotAvailable().Return(false)
				cli.EXPECT().GetRepeatedPassword().Return("passwod", nil)
				cli.EXPECT().PrintErr(errPasswordMismatch.Error())

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			wantErr: true,
		},
		{
			name: "register local with repeated pass err",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username3", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)
				cli.EXPECT().TryAgain().Return(false)
				grpc.EXPECT().IsNotAvailable().Return(false)
				cli.EXPECT().GetRepeatedPassword().Return("", testErr)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			wantErr: true,
		},
		{
			name: "login online",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(true)
				grpc.EXPECT().Login("username", gomock.Any()).Return(nil)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "login online with err",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil).AnyTimes()
				grpc.EXPECT().TryToConnect().Return(true).AnyTimes()
				grpc.EXPECT().Login("username", gomock.Any()).Return(testErr)
				grpc.EXPECT().Login("username", gomock.Any()).Return(nil)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "login local",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "check user online",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username_check", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(true)
				grpc.EXPECT().Login("username_check", gomock.Any()).Return(testErr)
				cli.EXPECT().TryAgain().Return(true)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
		{
			name: "check user online",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)
				grpc := mock_grpc.NewMockClient(mockCtrl)
				cli.EXPECT().GetUser().Return(models.User{Username: "username_check", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(true)
				grpc.EXPECT().Login("username_check", gomock.Any()).Return(nil)
				grpc.EXPECT().BlockStore(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

				cli.EXPECT().TryAgain().Return(true)
				cli.EXPECT().GetUser().Return(models.User{Username: "username", Password: "password"}, nil)
				grpc.EXPECT().TryToConnect().Return(false)

				usecase := ClientUsecase{
					userData:    &models.UserData{},
					dataPath:    file,
					fileManager: fm,
					grpc:        grpc,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.usecase.SetUser(); (err != nil) != tt.wantErr {
				t.Errorf("SetUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClientUsecase_SetFileManager(t *testing.T) {
	log := zerolog.New(os.Stdout)
	path := "path"
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
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				userData:   &models.UserData{},
				dataPath:   path,
				uploadPath: path,
				aboutMsg:   path,
				ui:         nil,
				log:        log,
			},
			wantErr: false,
		},
		{
			name: "!ok",
			fields: fields{
				userData:   &models.UserData{},
				dataPath:   "../../../../../../../../../..//../../../../../..test/../../..//",
				uploadPath: path,
				aboutMsg:   path,
				ui:         nil,
				log:        log,
			},
			wantErr: true,
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
			if err := c.SetFileManager(); (err != nil) != tt.wantErr {
				t.Errorf("SetFileManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewUsecase(t *testing.T) {
	log := zerolog.New(os.Stdout)
	path := "path"
	type args struct {
		dataPath      string
		uploadPath    string
		aboutMsg      string
		serverAddress string
		ui            ui.UserInterface
		log           zerolog.Logger
	}
	tests := []struct {
		name string
		args args
		want ClientUsecase
	}{
		{
			name: "ok",
			args: args{
				dataPath:      path,
				uploadPath:    path,
				aboutMsg:      path,
				serverAddress: "test",
				ui:            nil,
				log:           log,
			},
			want: ClientUsecase{
				dataPath:   path,
				uploadPath: path,
				aboutMsg:   path,
				grpc:       grpc.NewClient("test", log),
				ui:         nil,
				log:        log,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUsecase(
				tt.args.dataPath,
				tt.args.uploadPath,
				tt.args.aboutMsg,
				tt.args.serverAddress,
				tt.args.ui,
				tt.args.log,
			); got == nil {
				t.Errorf("NewUsecase() = %v, want %v", got, tt.want)
			}
		})
	}
}
