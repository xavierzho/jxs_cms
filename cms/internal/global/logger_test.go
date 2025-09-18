package global

import (
	"testing"

	"data_backend/pkg/logger"
)

func TestLogger(t *testing.T) {
	logger.LoggerSetting = &logger.LoggerConfig{
		LogSavePath: ".",
		LogFileName: "test",
		LogFileExt:  ".log",
	}
	SetupLogger()

	Logger.Info("test print")
}
