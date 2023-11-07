package filemanager_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/utils"
)

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		dataType int
		want     string
	}{
		{
			name:     "NoteType",
			dataType: filemanager.NoteType,
			want:     filemanager.NameNoteType,
		},
		{
			name:     "CardDataType",
			dataType: filemanager.CardDataType,
			want:     filemanager.NameCardDataType,
		},
		{
			name:     "BasicAuthType",
			dataType: filemanager.BasicAuthType,
			want:     filemanager.NameBasicAuthType,
		},
		{
			name:     "FileType",
			dataType: filemanager.FileType,
			want:     filemanager.NameFileType,
		},
		{
			name:     "unknown",
			dataType: 5,
			want:     "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filemanager.GetTypeName(tt.dataType); got != tt.want {
				t.Errorf("GetTypeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_String(t *testing.T) {
	tests := []struct {
		name string
		data filemanager.Data
		want string
	}{
		{
			name: "Note",
			data: filemanager.Note{Text: "test"},
			want: `Note : "test"`,
		},
		{
			name: "BasicAuth",
			data: filemanager.BasicAuth{Login: "test", Password: "test"},
			want: `Login: "test", Password: "test"`,
		},
		{
			name: "CardData",
			data: filemanager.CardData{Number: "123123123", Date: "02/22", CVV: "123"},
			want: `Number: "123123123", Date: "02/22", CVV: "123"`,
		},
		{
			name: "File",
			data: filemanager.File{Data: []byte("test")},
			want: `Successfully upload file`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_EncryptedData(t *testing.T) {
	key := []byte("superSuperTestSecretKeyWithSalt!")
	tests := []struct {
		name    string
		data    filemanager.Data
		key     []byte
		wantErr bool
	}{
		{
			name:    "Note",
			data:    filemanager.Note{Text: "test"},
			key:     key,
			wantErr: false,
		},
		{
			name:    "BasicAuth",
			data:    filemanager.BasicAuth{Login: "test", Password: "test"},
			key:     key,
			wantErr: false,
		},
		{
			name:    "CardData",
			data:    filemanager.CardData{Number: "123123123", Date: "02/22", CVV: "123"},
			key:     key,
			wantErr: false,
		},
		{
			name:    "File",
			data:    filemanager.File{Data: []byte("test")},
			key:     key,
			wantErr: false,
		},
		{
			name:    "Note key err",
			data:    filemanager.Note{Text: "test"},
			key:     []byte(""),
			wantErr: true,
		},
		{
			name:    "BasicAuth key err",
			data:    filemanager.BasicAuth{Login: "test", Password: "test"},
			key:     []byte(""),
			wantErr: true,
		},
		{
			name:    "CardData key err",
			data:    filemanager.CardData{Number: "123123123", Date: "02/22", CVV: "123"},
			key:     []byte(""),
			wantErr: true,
		},
		{
			name:    "File key err",
			data:    filemanager.File{Data: []byte("test")},
			key:     []byte(""),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.data.EncryptedData(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncryptedData() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !tt.wantErr {
				testData, err := json.Marshal(tt.data)
				if err != nil {
					t.Errorf("json.Marshal() error = %v", err)

					return
				}
				gotData, err := utils.Decrypt(got, tt.key)
				if err != nil {
					t.Errorf("utils.Decrypt() error = %v", err)

					return
				}
				if !reflect.DeepEqual(testData, gotData) {
					t.Errorf("EncryptedData() got = %v, want %v", gotData, testData)
				}
			}
		})
	}
}
