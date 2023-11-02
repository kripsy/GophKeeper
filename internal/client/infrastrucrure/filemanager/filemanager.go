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
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
)

const chunkSize = 3 * 1024 * 1024 // 3 МБ

type FileManager struct {
	storageDir string
	uploadDir  string
	userDir    string
	key        []byte
	meta       models.UserMeta
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
		if err = os.MkdirAll(storageDir, dirMode); err != nil {
			return nil, err
		}
	}

	err := meta.GetHash()
	if err != nil {
		return nil, err
	}

	fm := &FileManager{
		storageDir: storageDir,
		uploadDir:  uploadDir,
		userDir:    userDir,
		meta:       meta,
		key:        key,
	}

	return fm, nil

}

func (fm *FileManager) AddToStorage(name string, data Data, info models.DataInfo) error {
	_, ok := fm.meta.Data[name]
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

	if err = os.WriteFile(filepath.Join(fm.storageDir, info.DataID), encryptedData, fileMode); err != nil {
		return err
	}

	fm.meta.Data[name] = info

	return fm.saveMetaData(info.UpdatedAt)
}

func (fm *FileManager) AddEncryptedToStorage(name string, data chan []byte, info models.DataInfo) error {
	_, ok := fm.meta.Data[name]
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

	fm.meta.Data[name] = info

	return fm.saveMetaData(info.UpdatedAt)
}

func (fm *FileManager) GetByName(name string) ([]byte, models.DataInfo, error) {
	info, ok := fm.meta.Data[name]
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

			data <- buffer[:n]
		}
	}(data)

	return data, nil
}

func (fm *FileManager) UpdateDataByName(name string, data Data) error {
	savedInfo, ok := fm.meta.Data[name]
	if !ok {
		return errors.New("not found secret")
	}

	encryptedData, err := data.EncryptedData(fm.key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(fm.storageDir, savedInfo.DataID), encryptedData, fileMode); err != nil {
		return err
	}

	return fm.saveMetaData(time.Now())
}

func (fm *FileManager) UpdateInfoByName(name string, info models.DataInfo) error {
	savedInfo, ok := fm.meta.Data[name]
	if !ok {
		return errors.New("not found secret")
	}

	info.DataID = savedInfo.DataID
	info.UpdatedAt = time.Now()

	if info.Name == "" {
		info.Name = savedInfo.Name
	}

	if info.Description == "" {
		info.Description = savedInfo.Description
	}

	delete(fm.meta.Data, name)

	fm.meta.Data[info.Name] = info

	return fm.saveMetaData(time.Now())
}

func (fm *FileManager) DeleteByName(name string) error {
	info, ok := fm.meta.Data[name]
	if !ok {
		return errors.New("not found secret")
	}

	err := os.Remove(filepath.Join(fm.storageDir, info.DataID))
	if err != nil {
		return err
	}

	delete(fm.meta.Data, name)

	return fm.saveMetaData(time.Now())
}

func (fm *FileManager) saveMetaData(updatedAt time.Time) error {
	//fm.meta.UpdatedAt = updatedAt
	data, err := json.Marshal(fm.meta)
	if err != nil {
		return err // todo удалить данные в случае ошибки, загрузите повторно
	}

	encrypt, err := utils.Encrypt(data, fm.key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(fm.userDir), encrypt, fileMode); err != nil {
		return err
	}

	err = fm.meta.GetHash()
	if err != nil {
		return err
	}

	return nil
}

func (fm *FileManager) getUniqueName(name string) string {
	counter := 1
	baseName := name
	for {
		_, ok := fm.meta.Data[baseName]
		if !ok {
			return baseName
		}
		baseName = fmt.Sprintf("%s(%d)", name, counter)
		counter++
	}
}
