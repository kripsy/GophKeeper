package filemanager

import "errors"

var (
	errNotEqualData   = errors.New("error compared user Data")
	errNotFoundSecret = errors.New("not found secret")
)
