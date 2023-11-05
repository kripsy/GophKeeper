package filemanager

import (
	"bytes"
	"encoding/json"
	"github.com/kripsy/GophKeeper/internal/models"
	"os"
	"path/filepath"
	"sync"
	"testing"
	//"time"
	//	"github.com/google/uuid"
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
		storage FileStorage
		data    Data
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok Note",
			storage: func() FileStorage {
				fs, err := NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    Note{Text: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok Card",
			storage: func() FileStorage {
				fs, err := NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    CardData{Number: "test", Date: "test", CVV: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok BasicAuth",
			storage: func() FileStorage {
				fs, err := NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    BasicAuth{Login: "test", Password: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "ok File",
			storage: func() FileStorage {
				fs, err := NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data: File{Data: []byte("test")},
			info: models.DataInfo{
				Name:     "file",
				FileName: &filename,
			},
			wantErr: false,
		},
		{
			name: "ok with duplicate name",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    Note{Text: "test"},
			info:    info,
			wantErr: false,
		},
		{
			name: "failed with short encryption key",
			storage: func() FileStorage {
				fs, err := NewFileManager(storageDir, userDir, userDir, models.UserMeta{Data: make(models.MetaData)}, []byte("testKey"))
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:    Note{Text: "test"},
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
	info := models.DataInfo{
		Name:        "NewTestData",
		Description: "NewTest data description",
	}
	defer os.RemoveAll(storageDir)

	tests := []struct {
		name       string
		storage    FileStorage
		data       Data
		secretName string
		newInfo    models.DataInfo
		wantErr    bool
	}{
		{
			name: "ok update name",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
			secretName: "testData",
			newInfo: models.DataInfo{
				Name: "NewTestData",
			},
			wantErr: false,
		},
		{
			name: "ok update description",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
			secretName: "testData",
			newInfo: models.DataInfo{
				Description: "NewTest data description",
			},
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
			secretName: "testData",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.storage.UpdateInfoByName(tt.secretName, tt.newInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateInfoByName() error = %v, wantErr %v", err, tt.wantErr)
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
	note := Note{Text: "test"}
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
		storage FileStorage
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), data, fileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
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
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, []byte("testKey"))
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
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = "filename"
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
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
				t.Errorf("GetByName() error =  recived data not equal with expected data")
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
		storage FileStorage
		info    models.DataInfo
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}
				if err = os.WriteFile(filepath.Join(storageDir, filename), nil, fileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
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
		storage    FileStorage
		data       Data
		secretName string
		info       models.DataInfo
		wantErr    bool
	}{
		{
			name: "ok",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
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
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
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
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				meta.Data["testData"] = info
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			data:       Note{Text: "test"},
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

	note := Note{Text: "test"}
	data, err := json.Marshal(note)
	if err != nil {
		t.Fatalf("Failed to create testdata: %v", err)
	}

	tests := []struct {
		name    string
		storage FileStorage
		dataID  string
		want    chan []byte
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), data, fileMode); err != nil {
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

func TestFileManager_UpdateDataByName(t *testing.T) {
	defer os.RemoveAll(storageDir)
	info := models.DataInfo{
		Name:        "testData",
		Description: "Test data description",
	}

	tests := []struct {
		name    string
		storage FileStorage
		info    models.DataInfo
		data    Data
		wantErr bool
	}{
		{
			name: "ok",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				if err = os.WriteFile(filepath.Join(storageDir, filename), nil, fileMode); err != nil {
					t.Fatalf("Failed to create testFile: %v", err)
				}

				return fs
			}(),
			info:    info,
			data:    Note{Text: "test"},
			wantErr: false,
		},
		{
			name: "failed if secret not exist",
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, testKey)
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
			storage: func() FileStorage {
				meta := models.UserMeta{Data: make(models.MetaData)}
				testInfo := info
				testInfo.DataID = filename
				meta.Data["testData"] = testInfo
				fs, err := NewFileManager(storageDir, userDir, userDir, meta, []byte("testKey"))
				if err != nil {
					t.Fatalf("Failed to create FileManager: %v", err)
				}

				return fs
			}(),
			info:    info,
			data:    Note{Text: "test"},
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
