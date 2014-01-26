package sync

import (
	"errors"
)

var (
	ErrSyncAlreadyRunning  error = errors.New("sync already running")
	ErrChecksumNotMatching error = errors.New("checksum mismatch")
)
