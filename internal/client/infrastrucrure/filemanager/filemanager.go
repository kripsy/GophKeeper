// Package filemanager provides the functionality to manage secret files in the GophKeeper application,
// regardless of their type. It includes capabilities for adding, reading, updating, and deleting
// both regular and encrypted files, as well as managing metadata associated with these files.
package filemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"golang.org/x/sync/errgroup"
)

// chunkSize defines the size of chunks when reading encrypted files.
const (
	chunkSize       = 4 * 1000 * 1000 // 4 МБ
	cipherAmendment = 28
)

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
	// AddFileToStorage adds file data to storage and updates metadata.
	AddFileToStorage(newFile bool, name string, filePath string, info models.DataInfo) error
	// AddEncryptedToStorage adds encrypted data to storage and updates metadata.
	AddEncryptedToStorage(name string, data chan []byte, info models.DataInfo) error
	// GetByInfo retrieves decrypted data and metadata by info.
	GetByInfo(info models.DataInfo) ([]byte, models.DataInfo, error)
	// ReadEncryptedByID retrieves encrypted data by data ID.
	ReadEncryptedByID(name string) (chan []byte, error)
	// ReadFileFromStorage decrypts and moves the file to the specified path.
	ReadFileFromStorage(filePath string, info models.DataInfo) error
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

// AddToStorage adds data to storage and updates metadata.
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

// AddFileToStorage adds file data to storage and updates metadata.
func (fm *FileManager) AddFileToStorage(newFile bool, name string, filePath string, info models.DataInfo) error {
	_, ok := fm.Meta.Data[name]
	if ok && newFile {
		return fm.AddFileToStorage(newFile, fm.getUniqueName(name), filePath, info)
	}

	info.UpdatedAt = time.Now()

	dataChan := make(chan []byte)

	g, gctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		//nolint:revive
		if err := fm.readFile(gctx, filePath, chunkSize, utils.Encrypt, dataChan); err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		//nolint:revive
		if err := fm.writeFile(gctx, filepath.Join(fm.storageDir, info.DataID), dataChan); err != nil {
			return err
		}

		return nil
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("%w", err)
	}

	fm.Meta.Data[name] = info

	return fm.saveMetaData()
}

// ReadFileFromStorage decrypts and moves the file to the specified path.
func (fm *FileManager) ReadFileFromStorage(uploadDir string, info models.DataInfo) error {
	dataChan := make(chan []byte)
	g, gctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		//nolint:revive
		if err := fm.readFile(
			gctx,
			filepath.Join(fm.storageDir, info.DataID),
			chunkSize+cipherAmendment,
			utils.Decrypt,
			dataChan,
		); err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		//nolint:revive
		if err := fm.writeFile(gctx, filepath.Join(uploadDir, *info.FileName), dataChan); err != nil {
			return err
		}

		return nil
	})

	err := g.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("%w", err)
	}

	return nil
}

// AddEncryptedToStorage adds encrypted data to storage and updates metadata.
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

// GetByInfo retrieves decrypted data and metadata by info.
func (fm *FileManager) GetByInfo(info models.DataInfo) ([]byte, models.DataInfo, error) {
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

// ReadEncryptedByID retrieves encrypted data by data ID.
func (fm *FileManager) ReadEncryptedByID(dataID string) (chan []byte, error) {
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

// UpdateDataByName updates encrypted data in storage.
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

// UpdateInfoByName updates metadata by name.
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

// DeleteByName deletes data and metadata by name.
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

// DeleteByIDs deletes data and metadata by ids.
func (fm *FileManager) DeleteByIDs(ids []string) error {
	var delErr error
	for _, id := range ids {
		if _, err := os.Stat(filepath.Join(fm.storageDir, id)); !os.IsNotExist(err) {
			err = os.Remove(filepath.Join(fm.storageDir, id))
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

// Additional private methods like saveMetaData, getUniqueName, readFile, writeFile
// are also implemented in FileManager for internal file and metadata handling.
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
func (fm *FileManager) readFile(
	ctx context.Context,
	path string,
	chunkSize int,
	cipher func(data []byte, cipherKey []byte) ([]byte, error),
	dataChan chan<- []byte,
) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()

	buffer := make([]byte, chunkSize)

reader:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", ctx.Err())
		default:
			n, err := file.Read(buffer)
			if err != nil && !errors.Is(err, io.EOF) {
				return fmt.Errorf("%w", err)
			}
			if n == 0 {
				break reader
			}
			data, err := cipher(buffer[:n], fm.key)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
			dataChan <- data
		}
	}

	close(dataChan)

	return nil
}

func (fm *FileManager) writeFile(ctx context.Context, path string, dataChan <-chan []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	defer file.Close()

writer:
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", ctx.Err())
		case data, ok := <-dataChan:
			if !ok {
				break writer
			}
			_, err = file.Write(data)
			if err != nil {
				return fmt.Errorf("%w", err)
			}
		}
	}

	return nil
}
