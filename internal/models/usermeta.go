package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"
)

type MetaData map[string]DataInfo
type Deleted map[string]struct{}

type UserMeta struct {
	Username      string   `json:"user_name"`
	IsSyncStorage bool     `json:"is_local_storage"`
	Data          MetaData `json:"data"`
	DeletedData   Deleted  `json:"deleted_data"`
	HashData      string   `json:"-"`
}

type DataInfo struct {
	Name        string    `json:"name,omitempty"`
	DataID      string    `json:"data_id"`
	DataType    int       `json:"data_type"`
	Description string    `json:"description"`
	FileName    *string   `json:"file_name,omitempty"`
	Hash        string    `json:"hash"`
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

//nolint:godox
func (md *UserMeta) GetHash() error { // todo delete
	meta, err := json.Marshal(md.Data)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	md.HashData = fmt.Sprintf("%x", sha256.Sum256(meta))

	return nil
}
