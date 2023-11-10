//nolint:staticcheck
package usecase

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kripsy/GophKeeper/internal/client/grpc"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
)

func TestClientUsecaseInMenu(t *testing.T) {
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
		name          string
		fields        fields
		menuInput     string
		expectedCalls int
	}{
		{
			name: "Test Add Secret",
			fields: fields{
				userData: &models.UserData{
					Meta: models.UserMeta{
						IsSyncStorage: false,
					},
				},
			},
			menuInput:     ui.SyncSecrets,
			expectedCalls: 1,
		},
		{
			name: "Test Sync Secrets",
			fields: fields{
				userData: &models.UserData{
					Meta: models.UserMeta{
						IsSyncStorage: false,
					},
				},
			},
			menuInput:     ui.AddSecretKey,
			expectedCalls: 1,
		},
		{
			name: "Test SecretsKey",
			fields: fields{
				userData: &models.UserData{
					Meta: models.UserMeta{
						IsSyncStorage: false,
					},
				},
			},
			menuInput:     ui.SecretsKey,
			expectedCalls: 1,
		},
		// {
		// 	name: "Test ExitKey",
		// 	fields: fields{
		// 		userData: &models.UserData{
		// 			Meta: models.UserMeta{
		// 				IsSyncStorage: false,
		// 			},
		// 		},
		// 	},
		// 	menuInput:     ui.ExitKey,
		// 	expectedCalls: 1,
		// },
		{
			name: "Test About",
			fields: fields{
				userData: &models.UserData{
					Meta: models.UserMeta{
						IsSyncStorage: false,
					},
				},
			},
			menuInput:     ui.About,
			expectedCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUI := mock_ui.NewMockUserInterface(ctrl)
			c := &ClientUsecase{
				dataPath:    tt.fields.dataPath,
				uploadPath:  tt.fields.uploadPath,
				aboutMsg:    tt.fields.aboutMsg,
				userData:    tt.fields.userData,
				grpc:        tt.fields.grpc,
				fileManager: tt.fields.fileManager,
				ui:          mockUI,
				log:         tt.fields.log,
			}
			ok, val := getPositionMenu(tt.menuInput)

			if ok {
				mockUI.EXPECT().Menu(gomock.Any()).Return(val).AnyTimes()
				mockUI.EXPECT().ChooseSecretType().Return(1, false).AnyTimes()
				mockUI.EXPECT().GetSecret(gomock.Any()).Return("a", false).AnyTimes()
				mockUI.EXPECT().UpdateSecret(gomock.Any()).Return("a", 1, true).AnyTimes()
				mockUI.EXPECT().AddMetaInfo().Return(models.DataInfo{
					Name:        "asd",
					DataID:      "aaa",
					DataType:    1,
					Description: "asd",
					Hash:        "asd",
					UpdatedAt:   time.Now(),
				}, nil).AnyTimes()
				//nc (c *CLI) AddMetaInfo() (models.DataInfo, error)
				// mockUI.EXPECT().ChooseSecretType().Return(gomock.Any(), gomock.Any()).AnyTimes()

				go c.InMenu()
			}
			time.Sleep(2 * time.Second)
		})
	}
}

func getPositionMenu(str string) (bool, int) {
	for k, v := range ui.MenuTable {
		if v == str {
			return true, k
		}
	}
	return false, 0
}
