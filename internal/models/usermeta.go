package models

import (
	"crypto/sha256"
	"encoding/json"
	"path/filepath"
	"time"
)

type MetaData map[string]DataInfo

type UserMeta struct {
	Username       string    `json:"user_name"`
	IsLocalStorage bool      `json:"is_local_storage"`
	Data           MetaData  `json:"data"`
	HashData       [32]byte  `json:"-"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DataInfo struct {
	Name        string    `json:"name,omitempty"` // todo  подумать
	DataID      string    `json:"data_id"`
	DataType    int       `json:"data_type"`
	Description string    `json:"description""`
	FileName    *string   `json:"file_name,omitempty"`
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

func (md UserMeta) GetHash() error {
	meta, err := json.Marshal(md.Data)
	if err != nil {
		return err
	}

	md.HashData = sha256.Sum256(meta)

	return nil
}
