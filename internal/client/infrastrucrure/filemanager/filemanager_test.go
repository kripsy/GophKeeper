//nolint:gochecknoglobals, cyclop
package filemanager_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/filemanager"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
)

var (
	testKey  = []byte("superSuperTestScretKeyWithSalt!!")
	filename = "test.txt"
)

const (
	storageDir = "temp_storage/"
	userDir    = storageDir + "user"
)

func TestFileManager_AddToStorage(t *testing.T) {
	info := models.DataInfo{
		Name:        "testData",
		Description: "Test data description",
	}
	defer os.RemoveAll(storageDir)

	tests := []struct {
		name    string
		storage filemanager.FileStorage
		data    filemanager.Data
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok Note",
			storage: func() filemanager.FileStorage {
				fs, err := filemanager.NewFileManager(
					storageDir,
					userDir,
					userDir,
					models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    filemanager.Note{Text: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok Card",
			storage: func() filemanager.FileStorage {
				fs, err := filemanager.NewFileManager(
					storageDir,
					userDir,
					userDir,
					models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    filemanager.CardData{Number: "test", Date: "test", CVV: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok BasicAuth",
			storage: func() filemanager.FileStorage {
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    filemanager.BasicAuth{Login: "test", Password: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok File",
			storage: func() filemanager.FileStorage {
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data: filemanager.File{Data: []byte("test")},
			info: models.DataInfo{
				Name:     "file",
				FileName: &filename,
			},
			wantErr: false,
		},
		{
			name: "ok with duplicate name",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    filemanager.Note{Text: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "failed with short encryption key",
			storage: func() filemanager.FileStorage {
				fs, err := filemanager.NewFileManager(
					storageDir,
					userDir,
					userDir,
					models.UserMeta{Data: make(models.MetaData)}, []byte("testKey"))
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    filemanager.Note{Text: "test"},
			info:    info,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.storage.AddToStorage(tt.info.Name, tt.data, tt.info); (err != nil) != tt.wantErr {
				t.Errorf("AddToStorage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if _, err := os.Stat(filepath.Join(storageDir, info.DataID)); os.IsNotExist(err) {
					t.Errorf("file not exist AddToStorage() error = %v", err)
				}
			}
		})
	}
}

func TestFileManager_UpdateInfoByName(t *testing.T) {
	defer os.RemoveAll(storageDir)

	tests := []struct {
		name          string
		storage       *filemanager.FileManager
		data          filemanager.Data
		secretName    string
		newSecretName string
		newInfo       models.DataInfo
		wantInfo      models.DataInfo
		wantErr       bool
	}{
		{
			name: "ok update name",
			storage: func() *filemanager.FileManager {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = models.DataInfo{
					Name:        "testData",
					Description: "test data description",
				}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:          filemanager.Note{Text: "test"},
			secretName:    "testData",
			newSecretName: "NewTestData",
			newInfo: models.DataInfo{
				Name: "NewTestData",
			},
			wantInfo: models.DataInfo{
				Name:        "NewTestData",
				Description: "test data description",
			},
			wantErr: false,
		},
		{
			name: "ok update description",
			storage: func() *filemanager.FileManager {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = models.DataInfo{
					Name:        "testData",
					Description: "test data description",
				}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:          filemanager.Note{Text: "test"},
			secretName:    "testData",
			newSecretName: "testData",
			newInfo: models.DataInfo{
				Description: "NewTest data description",
			},
			wantInfo: models.DataInfo{
				Name:        "testData",
				Description: "NewTest data description",
			},
			wantErr: false,
		},
		{
			name: "ok update filename",
			storage: func() *filemanager.FileManager {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = models.DataInfo{
					Name:        "testData",
					Description: "test data description",
				}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:          filemanager.Note{Text: "test"},
			secretName:    "testData",
			newSecretName: "testData",
			newInfo: models.DataInfo{
				FileName: &filename,
			},
			wantInfo: models.DataInfo{
				Name:        "testData",
				Description: "test data description",
				FileName:    &filename,
			},
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() *filemanager.FileManager {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       filemanager.Note{Text: "test"},
			secretName: "testData",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.storage.UpdateInfoByName(tt.secretName, tt.newInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInfoByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			got := tt.storage.Meta.Data[tt.newSecretName]
			if !tt.wantErr {
				if got.Name != tt.wantInfo.Name ||
					got.Description != tt.wantInfo.Description ||
					got.FileName != tt.wantInfo.FileName {
					t.Errorf("UpdateInfoByName() error Meta not equal")
				}
			}
		})
	}
}

func TestFileManager_UpdateDataByName(t *testing.T) {
	defer os.RemoveAll(storageDir)
	info := models.DataInfo{
		Name:        "testData",
		Description: "Test data description",
	}

	tests := []struct {
		name    string
		storage filemanager.FileStorage
		info    models.DataInfo
		data    filemanager.Data
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), nil, permissions.FileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			data:    filemanager.Note{Text: "test"},
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: true,
		},
		{
			name: "failed with shortness key",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, []byte("testKey"))
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			data:    filemanager.Note{Text: "test"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.storage.UpdateDataByName(tt.info.Name, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileManager_GetByName(t *testing.T) {
	defer os.RemoveAll(storageDir)
	info := models.DataInfo{
		Name:        "testData",
		Description: "Test data description",
	}
	note := filemanager.Note{Text: "test"}
	marshaledData, err := json.Marshal(note)
	if err != nil {
		t.Fatalf("Failed to create testdata: %v", err)
	}
	data, err := note.EncryptedData(testKey)
	if err != nil {
		t.Fatalf("Failed to create testdata: %v", err)
	}
	tests := []struct {
		name    string
		storage filemanager.FileStorage
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), data, permissions.FileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: true,
		},
		{
			name: "failed with shortness key",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, []byte("testKey"))
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: true,
		},
		{
			name: "failed with not existing file",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = "filename"
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testData, _, err := tt.storage.GetByName(tt.info.Name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(testData, marshaledData) {
				t.Errorf("GetByName() error = recived data not equal with expected data")
			}
		})
	}
}

func TestFileManager_DeleteByName(t *testing.T) {
	info := models.DataInfo{
		Name:        "testData",
		Description: "Test data description",
	}
	defer os.RemoveAll(storageDir)

	tests := []struct {
		name    string
		storage filemanager.FileStorage
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}
				if err = os.WriteFile(filepath.Join(storageDir, filename), nil, permissions.FileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.storage.DeleteByName(tt.info.Name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileManager_AddEncryptedToStorage(t *testing.T) {
	info := models.DataInfo{
		Name:        "NewTestData",
		Description: "NewTest data description",
	}
	defer os.RemoveAll(storageDir)

	tests := []struct {
		name       string
		storage    filemanager.FileStorage
		data       filemanager.Data
		secretName string
		info       models.DataInfo
		wantErr    bool
	}{
		{
			name: "ok",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       filemanager.Note{Text: "test"},
			secretName: "testData",
			info: models.DataInfo{
				Name:        "NewTestData",
				Description: "NewTest data description",
				DataID:      "filename",
			},
			wantErr: false,
		},
		{
			name: "ok duplicate name",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       filemanager.Note{Text: "test"},
			secretName: "testData",
			info: models.DataInfo{
				Name:        "testData",
				Description: "NewTest data description",
				DataID:      "filename",
			},
			wantErr: false,
		},
		{
			name: "failed without filename",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       filemanager.Note{Text: "test"},
			secretName: "testData",
			info: models.DataInfo{
				Name:        "NewTestData",
				Description: "NewTest data description",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data chan []byte
			var wg sync.WaitGroup
			wg.Add(1)

			go func() {
				data = make(chan []byte, 1)
				data <- []byte("test")
				close(data)
				wg.Done()
			}()

			wg.Wait()

			if err := tt.storage.AddEncryptedToStorage(tt.secretName, data, tt.info); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInfoByName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileManager_ReadEncryptedByName1(t *testing.T) {
	defer os.RemoveAll(storageDir)

	note := filemanager.Note{Text: "test"}
	data, err := json.Marshal(note)
	if err != nil {
		t.Fatalf("Failed to create testdata: %v", err)
	}

	tests := []struct {
		name    string
		storage filemanager.FileStorage
		dataID  string
		want    chan []byte
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() filemanager.FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := filemanager.NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), data, permissions.FileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			dataID:  filename,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			dataChan, err := tt.storage.ReadEncryptedByName(tt.dataID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadEncryptedByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				var testData []byte
				for chunk := range dataChan {
					testData = append(testData, chunk...)
				}

				if !bytes.Equal(data, testData) {
					t.Errorf("ReadEncryptedByName() recived data not equal with expected data")
				}
			}
		})
	}
}
