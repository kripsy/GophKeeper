package models

import (
	"errors"
	"fmt"
)

// SyncError represents a synchronization error with a custom message.
type SyncError struct {
	Err error
}

// NewSyncError creates a new synchronization error with the provided text.
func NewSyncError(errorNum SyncErrorEnum) error {
	err, ok := SyncErrorMessages[errorNum]
	if !ok {
		err = SyncErrorMessages[ErrSyncUnexpectedEnum] // Use a default error if the error number is not found
	}

	return fmt.Errorf("sync error: %w", err) // Wrap the static error
}

// Error returns the error message of the SyncError.
func (ue *SyncError) Error() string {
	return ue.Err.Error()
}

// SyncErrorEnum contains the enumeration of synchronization errors.
type SyncErrorEnum int

const (
	ErrUserSyncExistsEnum SyncErrorEnum = iota // This sync for this user already exists
	ErrSyncExistsEnum                          // Sync for this user already exists
	ErrSyncNotFoundEnum                        // Sync not found
	ErrSyncUnexpectedEnum                      // Unexpected error
)

var (
	ErrUserSyncExists = errors.New("this sync for this user already exists")
	ErrSyncExists     = errors.New("sync for this user already exists")
	ErrSyncNotFound   = errors.New("sync not found")
	ErrSyncUnexpected = errors.New("unexpected sync error")
)

// SyncErrorMessages contains error messages corresponding to SyncErrorEnum.
//
//nolint:gochecknoglobals
var SyncErrorMessages = map[SyncErrorEnum]error{
	ErrUserSyncExistsEnum: ErrUserSyncExists,
	ErrSyncExistsEnum:     ErrSyncExists,
	ErrSyncNotFoundEnum:   ErrSyncNotFound,
	ErrSyncUnexpectedEnum: ErrSyncUnexpected,
}
