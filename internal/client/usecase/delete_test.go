//nolint:testpackage,goerr113
package usecase

import (
	"fmt"
	"os"
	"testing"

	mock_filemanager "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager/mocks"
	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui"
	mock_ui "github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/mocks"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/rs/zerolog"
	"go.uber.org/mock/gomock"
)

func TestClientUsecase_deleteSecret(t *testing.T) {
	testErr := fmt.Errorf("error")
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
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
			name: "delete",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().DeleteByName("test").Return(nil)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: "test", success: true},
		},
		{
			name: "delete with err",
			usecase: func() ClientUsecase {
				cli := mock_ui.NewMockUserInterface(mockCtrl)
				cli.EXPECT().PrintErr(ui.DeleteErr)

				fm := mock_filemanager.NewMockFileStorage(mockCtrl)

				fm.EXPECT().DeleteByName("test").Return(testErr)
				usecase := ClientUsecase{
					userData:    &models.UserData{},
					fileManager: fm,
					ui:          cli,
					log:         log,
				}

				return usecase
			}(),
			args: args{secretName: "test", success: true},
		},
		{
			name: "delete exit",
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
			args: args{secretName: "test", success: false},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.usecase.deleteSecret(tt.args.secretName, tt.args.success)
		})
	}
}
