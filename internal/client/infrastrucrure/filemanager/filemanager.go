package filemanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/utils"
	"os"
	"path/filepath"
	"time"
)

type FileManager struct {
	storageDir string
	uploadDir  string
	userDir    string
	key        []byte
	meta       models.UserMeta
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

	newUUID, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	info.DataID = newUUID.String()
	info.UpdatedAt = time.Now()
	encryptedData, err := data.EncryptedData(fm.key)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(fm.storageDir, info.DataID), encryptedData, fileMode); err != nil {
		return err
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
	fm.meta.UpdatedAt = updatedAt
	data, err := json.Marshal(fm.meta)
	if err != nil {
		return err // удалить данные в случае ошибки, загрузите повторно
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
