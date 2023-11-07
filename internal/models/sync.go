package models

import (
	"errors"
)

// SyncError represents a synchronization error with a custom message.
type SyncError struct {
	Err error
}

// NewSyncError creates a new synchronization error with the provided text.
func NewSyncError(errorNum SyncErrorEnum) error {
	return &SyncError{
		Err: errors.New(SyncErrorMessages[errorNum]),
	}
}

// Error returns the error message of the SyncError.
func (ue *SyncError) Error() string {
	return ue.Err.Error()
}

// SyncErrorEnum contains the enumeration of synchronization errors.
type SyncErrorEnum int

const (
	ErrUserSyncExists SyncErrorEnum = iota // This sync for this user already exists
	ErrSyncExists                          // Sync for this user already exists
	ErrSyncNotFound                        // Sync not found
)

// SyncErrorMessages contains error messages corresponding to SyncErrorEnum.
//
//nolint:gochecknoglobals
var SyncErrorMessages = map[SyncErrorEnum]string{
	ErrUserSyncExists: "This sync for this user already exists",
	ErrSyncExists:     "Sync for this user already exists",
	ErrSyncNotFound:   "Sync not found",
}
