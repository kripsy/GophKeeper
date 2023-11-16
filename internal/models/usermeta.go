// Package models defines data structures used in the application, including metadata and user-related information.
package models

import (
	"path/filepath"
	"time"
)

// MetaData represents meta-information about secrets DataInfo by keys and secret name.
type MetaData map[string]DataInfo

// Deleted represents a set of deleted data IDs.
type Deleted map[string]struct{}

// UserMeta represents metadata associated with a user and secrets.
type UserMeta struct {
	Username      string   `json:"user_name"`
	IsSyncStorage bool     `json:"is_local_storage"`
	Data          MetaData `json:"data"`
	DeletedData   Deleted  `json:"deleted_data"`
	HashData      string   `json:"-"`
}

// DataInfo represents information about user data.
type DataInfo struct {
	Name        string    `json:"name,omitempty"`
	DataID      string    `json:"data_id"`
	DataType    int       `json:"data_type"`
	Description string    `json:"description"`
	FileName    *string   `json:"file_name,omitempty"`
	Hash        string    `json:"hash"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SetFileName sets the FileName field based on the provided path.
func (di *DataInfo) SetFileName(path string) {
	fileName := filepath.Base(path)
	if fileName == "" {
		di.FileName = &path

		return
	}
	di.FileName = &fileName
}

// IsDeleted checks if a data ID is marked as deleted in the Deleted set.
func (d Deleted) IsDeleted(dataID string) bool {
	_, deleted := d[dataID]

	return deleted
}
