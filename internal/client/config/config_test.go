package config_test

import (
	"github.com/kripsy/GophKeeper/internal/client/config"
	"os"
	"reflect"
	"testing"
)

func TestGetConfig1(t *testing.T) {
	oldArgs := os.Args
	tests := []struct {
		name    string
		args    []string
		want    config.Config
		wantErr bool
	}{
		{
			name: "ok with flag",
			args: append(oldArgs,
				"",
				"-cfg-path", "/config.yaml",
				"-storage-path", "/path/to/storage",
				"-upload-path", "/path/to/uploads",
				"-server-addr", "127.0.0.1:8080"),
			want: config.Config{
				StoragePath:   "./1keeper/Data",
				UploadPath:    "./1keeper/Upload",
				ServerAddress: "127.0.0.1:50051",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			got, err := config.GetConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfig() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
