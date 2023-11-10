// Package filemanager provides the ability to work with secret files regardless of their type
package filemanager

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
)

// chunkSize defines the size of chunks when reading encrypted files.
const chunkSize = 4 * 1000 * 1000 // 4 МБ

// FileManager struct represents a manager for secret files, providing methods for file operations.
type FileManager struct {
	storageDir string
	uploadDir  string
	userDir    string
	key        []byte
	Meta       models.UserMeta
}

// FileStorage interface defines methods for file operations.
type FileStorage interface {
	// AddToStorage adds data to storage and updates metadata.
	AddToStorage(name string, data Data, info models.DataInfo) error
	// AddEncryptedToStorage adds encrypted data to storage and updates metadata.
	AddEncryptedToStorage(name string, data chan []byte, info models.DataInfo) error
	// GetByName retrieves decrypted data and metadata by name.
	GetByName(name string) ([]byte, models.DataInfo, error)
	// ReadEncryptedByName retrieves encrypted data by data ID.
	ReadEncryptedByName(name string) (chan []byte, error)
	// UpdateDataByName updates encrypted data in storage.
	UpdateDataByName(name string, data Data) error
	// UpdateInfoByName updates metadata by name.
	UpdateInfoByName(name string, info models.DataInfo) error
	// DeleteByIDs deletes data and metadata by ids.
	DeleteByIDs(ids []string) error
	// DeleteByName deletes data and metadata by name.
	DeleteByName(name string) error
}

// NewFileManager creates a new FileManager instance with the provided parameters.
func NewFileManager(storageDir, uploadDir, userDir string, meta models.UserMeta, key []byte) (*FileManager, error) {
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err = os.MkdirAll(storageDir, permissions.DirMode); err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	}

	return &FileManager{
		storageDir: storageDir,
		uploadDir:  uploadDir,
		userDir:    userDir,
		Meta:       meta,
		key:        key,
	}, nil
}

func (fm *FileManager) AddToStorage(name string, data Data, info models.DataInfo) error {
	_, ok := fm.Meta.Data[name]
	if ok {
		return fm.AddToStorage(fm.getUniqueName(name), data, info)
	}

	info.DataID = uuid.New().String()
	info.UpdatedAt = time.Now()
	encryptedData, err := data.EncryptedData(fm.key)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	hash, err := data.GetHash()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	info.Hash = hash

	if err = os.WriteFile(filepath.Join(fm.storageDir, info.DataID), encryptedData, permissions.FileMode); err != nil {
		return fmt.Errorf("%w", err)
	}

	fm.Meta.Data[name] = info

	return fm.saveMetaData()
}

func (fm *FileManager) AddEncryptedToStorage(name string, data chan []byte, info models.DataInfo) error {
	_, ok := fm.Meta.Data[name]
	if ok {
		return fm.AddEncryptedToStorage(fm.getUniqueName(name), data, info)
	}

	outFile, err := os.Create(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer outFile.Close()

	// Write encrypted data to the file in chunks.
	for chunk := range data {
		if _, writeErr := outFile.Write(chunk); writeErr != nil {
			return fmt.Errorf("%w", err)
		}
	}

	fm.Meta.Data[name] = info

	return fm.saveMetaData()
}

func (fm *FileManager) GetByName(name string) ([]byte, models.DataInfo, error) {
	info, ok := fm.Meta.Data[name]
	if !ok {
		return nil, models.DataInfo{}, fmt.Errorf("%w", errNotFoundSecret)
	}

	data, err := os.ReadFile(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return nil, models.DataInfo{}, fmt.Errorf("%w", err)
	}

	decryptedData, err := utils.Decrypt(data, fm.key)
	if err != nil {
		return nil, models.DataInfo{}, fmt.Errorf("%w", err)
	}

	return decryptedData, info, nil
}

func (fm *FileManager) ReadEncryptedByName(dataID string) (chan []byte, error) {
	file, err := os.Open(filepath.Join(fm.storageDir, dataID))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	data := make(chan []byte, 1)
	buffer := make([]byte, chunkSize)

	// Start a goroutine to read the file in chunks.
	go func(data chan []byte) {
		defer file.Close()
		defer close(data)
		for {
			n, err := file.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
			}

			chunk := make([]byte, n)
			copy(chunk, buffer[:n])
			data <- chunk
		}
	}(data)

	return data, nil
}

func (fm *FileManager) UpdateDataByName(name string, data Data) error {
	savedInfo, ok := fm.Meta.Data[name]
	if !ok {
		return fmt.Errorf("%w", errNotFoundSecret)
	}

	encryptedData, err := data.EncryptedData(fm.key)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.WriteFile(
		filepath.Join(fm.storageDir, savedInfo.DataID),
		encryptedData,
		permissions.FileMode,
	); err != nil {
		return fmt.Errorf("%w", err)
	}

	return fm.saveMetaData()
}

func (fm *FileManager) UpdateInfoByName(name string, info models.DataInfo) error {
	savedInfo, ok := fm.Meta.Data[name]
	if !ok {
		return fmt.Errorf("%w", errNotFoundSecret)
	}

	if info.Name != "" && savedInfo.Name != info.Name {
		savedInfo.Name = info.Name
		delete(fm.Meta.Data, name)
	}

	if info.Description != "" {
		savedInfo.Description = info.Description
	}

	if info.FileName != nil {
		savedInfo.FileName = info.FileName
	}

	savedInfo.UpdatedAt = time.Now()
	fm.Meta.Data[savedInfo.Name] = savedInfo

	return fm.saveMetaData()
}

func (fm *FileManager) DeleteByName(name string) error {
	info, ok := fm.Meta.Data[name]
	if !ok {
		return fmt.Errorf("%w", errNotFoundSecret)
	}

	err := os.Remove(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	delete(fm.Meta.Data, name)
	fm.Meta.DeletedData[info.DataID] = struct{}{}

	return fm.saveMetaData()
}

func (fm *FileManager) DeleteByIDs(ids []string) error {
	var delErr error
	for _, id := range ids {
		if _, err := os.Stat(filepath.Join(fm.storageDir, id)); !os.IsNotExist(err) {
			err := os.Remove(filepath.Join(fm.storageDir, id))
			if err != nil {
				delErr = fmt.Errorf("%w, %s: %w", delErr, id, err)
			}
		}

	finder:
		for name, info := range fm.Meta.Data {
			if info.DataID == id {
				delete(fm.Meta.Data, name)

				break finder
			}
		}
	}
	if err := fm.saveMetaData(); err != nil {
		delErr = fmt.Errorf("%w, save meta err: %w", delErr, err)
	}
	if delErr != nil {
		return fmt.Errorf("%w", delErr)
	}

	return nil
}

func (fm *FileManager) saveMetaData() error {
	data, err := json.Marshal(fm.Meta)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	encrypt, err := utils.Encrypt(data, fm.key)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if err = os.WriteFile(fm.userDir, encrypt, permissions.FileMode); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func (fm *FileManager) getUniqueName(name string) string {
	counter := 1
	baseName := name
	for {
		_, ok := fm.Meta.Data[baseName]
		if !ok {
			return baseName
		}
		baseName = fmt.Sprintf("%s(%d)", name, counter)
		counter++
	}
}
