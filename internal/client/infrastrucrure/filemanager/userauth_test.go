package filemanager_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
)

const (
	Login    = "testUser"
	Password = "testPassword"
)

func Test_userAuth_GetUser(t *testing.T) {
	defer os.RemoveAll(storageDir)
	auth, err := filemanager.NewUserAuth(storageDir)
	if err != nil {
		t.Fatalf("Failed to create NewUserAuth: %v", err)
	}
	tests := []struct {
		name        string
		auth        filemanager.Auth
		user        *models.User
		prepareFunc func()
		want        models.UserMeta
		wantErr     bool
	}{
		{
			name: "ok",
			auth: auth,
			prepareFunc: func() {
				if _, err = auth.CreateUser(&models.User{Username: Login, Password: Password}, true); err != nil {
					t.Fatalf("Failed prepare user: %v", err)
				}
			},
			user: &models.User{Username: Login, Password: Password},
			want: models.UserMeta{Username: Login, IsSyncStorage: true,
				Data: make(models.MetaData), DeletedData: make(models.Deleted)},
			wantErr: false,
		},
		{
			name:        "failed not created user",
			auth:        auth,
			prepareFunc: func() {},
			user:        &models.User{Username: "notCreatedUser", Password: Password},
			wantErr:     true,
		},
		{
			name: "failed invalid password",
			auth: auth,
			prepareFunc: func() {
				if _, err = auth.CreateUser(&models.User{Username: Login, Password: Password}, true); err != nil {
					t.Fatalf("Failed prepare user: %v", err)
				}
			},
			user:    &models.User{Username: Login, Password: "invalidPassword"},
			wantErr: true,
		},
		{
			name: "failed invalid password",
			auth: auth,
			prepareFunc: func() {
				if _, err = auth.CreateUser(&models.User{Username: Login, Password: Password}, true); err != nil {
					t.Fatalf("Failed prepare user: %v", err)
				}
			},
			user:    &models.User{Username: Login, Password: "invalidPassword"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepareFunc()
			got, err := tt.auth.GetUser(tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userAuth_CreateUser(t *testing.T) {
	defer os.RemoveAll(storageDir)
	auth, err := filemanager.NewUserAuth(storageDir)
	if err != nil {
		t.Fatalf("Failed to create NewUserAuth: %v", err)
	}
	tests := []struct {
		name          string
		auth          filemanager.Auth
		user          *models.User
		want          models.UserMeta
		IsSyncStorage bool
		wantErr       bool
	}{
		{
			name: "ok",
			auth: auth,
			user: &models.User{Username: Login, Password: Password},
			want: models.UserMeta{Username: Login, IsSyncStorage: true,
				Data: make(models.MetaData), DeletedData: make(models.Deleted)},
			IsSyncStorage: true,
			wantErr:       false,
		},
		{
			name: "not ok",
			auth: auth,
			user: &models.User{
				Username: "../../../../../../../../../..//../../../../../..test/../../..//",
				Password: Password},
			want:          models.UserMeta{},
			IsSyncStorage: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.auth.CreateUser(tt.user, tt.IsSyncStorage)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userAuth_IsUserExisting(t *testing.T) {
	user := "user"
	defer os.RemoveAll(storageDir)
	auth, err := filemanager.NewUserAuth(storageDir)
	if err != nil {
		t.Fatalf("Failed to create NewUserAuth: %v", err)
	}
	if err = os.WriteFile(filepath.Join(storageDir, user), nil, permissions.FileMode); err != nil {
		t.Fatalf("Failed to create testFile: %v", err)
	}
	tests := []struct {
		name    string
		auth    filemanager.Auth
		userDit string
		want    bool
	}{
		{
			name:    "user existing",
			auth:    auth,
			userDit: storageDir + user,
			want:    true,
		},
		{
			name:    "user not existing",
			auth:    auth,
			userDit: storageDir + "test",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.auth.IsUseExisting(tt.userDit); got != tt.want {
				t.Errorf("IsUserNotExisting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserAuth(t *testing.T) {
	tests := []struct {
		name     string
		userPath string
		wantErr  bool
	}{
		{
			name:     "ok",
			userPath: "./..",
			wantErr:  false,
		},
		{
			name:     "not ok",
			userPath: "../../../../../../../../../..//../../../../../..test/../../..//",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := filemanager.NewUserAuth(tt.userPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUserAuth() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
		})
	}
}
