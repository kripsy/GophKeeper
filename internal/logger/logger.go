// Package logger provides logging functionalities for the GophKeeper application using zap logger.
// It includes methods to initialize the logger with various configurations.
package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// InitLog initializes a zap logger with the specified log level.
// It parses the log level, sets the configuration, and returns the configured logger.
// Returns an error if the log level parsing or logger configuration fails.
func InitLog(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)

	// lvl, err = zap.ParseAtomicLevel("Debug").

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	cfg := zap.NewProductionConfig()

	cfg.Level = lvl
	zl, err := cfg.Build()

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	log := zl
	log.Info(`Logger level`, zap.String("logLevel", level))

	return log, nil
}

// InitLogWithFilePath initializes a zap logger with the specified log level and output file path.
// It parses the log level, sets the output path for logging, and returns the configured logger.
// Returns an error if the log level parsing or logger configuration fails.
func InitLogWithFilePath(level string, logFilePath string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)

	// lvl, err = zap.ParseAtomicLevel("Debug").

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	cfg.OutputPaths = []string{
		logFilePath + "log.log",
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return logger, nil
}
