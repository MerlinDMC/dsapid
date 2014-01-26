package storage

import (
	"errors"
)

var (
	ErrStorageFileNotFound    error = errors.New("File not found")
	ErrStorageFileNotReadable error = errors.New("File not readable")
	ErrStorageFileNotWritable error = errors.New("File not writable")
	ErrStorageFileInvalid     error = errors.New("Syntax error in storage file")
	ErrStorageItemNotFound    error = errors.New("Item not available in storage")
)
