package filemanager

import (
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
)

const chunkSize = 4 * 1000 * 1000 // 4 МБ

type FileManager struct {
	storageDir string
	uploadDir  string
	userDir    string
	key        []byte
	Meta       models.UserMeta
}

type FileStorage interface {
	AddToStorage(name string, data Data, info models.DataInfo) error
	AddEncryptedToStorage(name string, data chan []byte, info models.DataInfo) error
	GetByName(name string) ([]byte, models.DataInfo, error)
	ReadEncryptedByName(name string) (chan []byte, error)
	UpdateDataByName(name string, data Data) error
	UpdateInfoByName(name string, info models.DataInfo) error
	DeleteByName(name string) error
}

func NewFileManager(storageDir, uploadDir, userDir string, meta models.UserMeta, key []byte) (*FileManager, error) {
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err = os.MkdirAll(storageDir, permissions.DirMode); err != nil {
			return nil, err
		}
	}

	err := meta.GetHash()
	if err != nil {
		return nil, err
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
		return err
	}
	hash, err := data.GetHash()
	if err != nil {
		return err
	}
	info.Hash = hash

	if err = os.WriteFile(filepath.Join(fm.storageDir, info.DataID), encryptedData, permissions.FileMode); err != nil {
		return err
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
		return err
	}

	defer outFile.Close()

	for chunk := range data {
		if _, writeErr := outFile.Write(chunk); writeErr != nil {
			return err
		}
	}

	fm.Meta.Data[name] = info

	return fm.saveMetaData()
}

func (fm *FileManager) GetByName(name string) ([]byte, models.DataInfo, error) {
	info, ok := fm.Meta.Data[name]
	if !ok {
		return nil, models.DataInfo{}, errors.New("not found secret")
	}

	data, err := os.ReadFile(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return nil, models.DataInfo{}, err
	}

	decryptedData, err := utils.Decrypt(data, fm.key)
	if err != nil {
		return nil, models.DataInfo{}, err
	}

	return decryptedData, info, nil
}

func (fm *FileManager) ReadEncryptedByName(dataID string) (chan []byte, error) {
	file, err := os.Open(filepath.Join(fm.storageDir, dataID))
	if err != nil {
		return nil, err
	}

	data := make(chan []byte, 1)
	buffer := make([]byte, chunkSize)

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
		return errors.New("not found secret")
	}

	encryptedData, err := data.EncryptedData(fm.key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(
		filepath.Join(fm.storageDir, savedInfo.DataID),
		encryptedData,
		permissions.FileMode,
	); err != nil {
		return err
	}

	return fm.saveMetaData()
}

func (fm *FileManager) UpdateInfoByName(name string, info models.DataInfo) error {
	savedInfo, ok := fm.Meta.Data[name]
	if !ok {
		return errors.New("not found secret")
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
		return errors.New("not found secret")
	}

	err := os.Remove(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return err
	}

	delete(fm.Meta.Data, name)

	return fm.saveMetaData()
}

func (fm *FileManager) saveMetaData() error {
	data, err := json.Marshal(fm.Meta)
	if err != nil {
		return err // todo удалить данные в случае ошибки, загрузите повторно
	}

	encrypt, err := utils.Encrypt(data, fm.key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Dir(fm.userDir), encrypt, permissions.FileMode); err != nil {
		return err
	}

	err = fm.Meta.GetHash()
	if err != nil {
		return err
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
