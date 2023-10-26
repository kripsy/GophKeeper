package log

import (
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

const fileName = "client.log"

func InitLogger(logPath string) zerolog.Logger {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		if err = os.MkdirAll(logPath, 0644); err != nil {
			panic(err)
		}
	}

	fileWriter, err := os.OpenFile(filepath.Join(logPath, fileName), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// Create a Zerolog logger that writes to the file writer.
	log := zerolog.New(fileWriter).With().Timestamp().Logger()
	log.Log().Msg("test")
	return log
}
