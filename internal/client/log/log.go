// Package log provides logging functionalities for the GophKeeper application.
// It includes methods for initializing and configuring the logger.
package log

import (
	"os"
	"path/filepath"

	"github.com/kripsy/GophKeeper/internal/client/permissions"
	"github.com/rs/zerolog"
)

const fileName = "client.log"

// InitLogger initializes the zerolog logger with the provided log path.
// It creates the log file if it doesn't exist and sets up the logger to write to this file.
// The logger includes caller and timestamp information in each log entry.
func InitLogger(logPath string) zerolog.Logger {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err = os.MkdirAll(logPath, permissions.DirMode); err != nil {
			panic(err)
		}
	}

	fileWriter, err := os.OpenFile(
		filepath.Join(logPath, fileName),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		permissions.FileMode)
	if err != nil {
		panic(err)
	}

	return zerolog.New(fileWriter).With().Caller().Timestamp().Logger()
}
