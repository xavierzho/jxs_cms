package local

import (
	"time"

	"data_backend/pkg/logger"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Logger *logger.Logger
)

func SetupLogger() (err error) {
	fileName := logger.LoggerSetting.LogSavePath + "/v2/" + logger.LoggerSetting.LogFileName + time.Now().Format(logger.FILE_LAYOUT) + logger.LoggerSetting.LogFileExt
	jackLog := &lumberjack.Logger{
		Filename:  fileName,
		MaxSize:   500,
		LocalTime: true,
	}
	Logger = logger.NewLogger(Ctx, jackLog)
	return nil
}
