package logger

import (
	"fmt"

	"go.uber.org/zap"
)

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
