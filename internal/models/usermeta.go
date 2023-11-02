package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"
)

type MetaData map[string]DataInfo

type UserMeta struct {
	Username       string   `json:"user_name"`
	IsLocalStorage bool     `json:"is_local_storage"`
	Data           MetaData `json:"data"`
	HashData       string   `json:"-"`
}

type DataInfo struct {
	Name        string    `json:"name,omitempty"`
	DataID      string    `json:"data_id"`
	DataType    int       `json:"data_type"`
	Description string    `json:"description"`
	FileName    *string   `json:"file_name,omitempty"`
	Hash        string    `json:"hash"`
	IsDeleted   bool      `json:"is_deleted"` // todo вынести в отдельное поле
	UpdatedAt   time.Time `json:"updated_at"`
}

func (di *DataInfo) SetFileName(path string) {
	fileName := filepath.Base(path)
	if fileName == "" {
		di.FileName = &path

		return
	}
	di.FileName = &fileName
}

func (md *UserMeta) GetHash() error { //todo delete
	meta, err := json.Marshal(md.Data)
	if err != nil {
		return err
	}

	md.HashData = fmt.Sprintf("%x", sha256.Sum256(meta))

	return nil
}
