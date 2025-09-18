package local

import (
	"testing"

	"data_backend/pkg/logger"
)

func SetupTest() {
	logger.LoggerSetting = &logger.LoggerConfig{
		LogSavePath: ".",
		LogFileName: "test",
		LogFileExt:  ".log",
	}
	if err := SetupLogger(); err != nil {
		panic(err)
	}

	DatabaseSetting = &DatabaseConfig{
		CenterDB:        "root:123456@tcp(192.168.8.92:3306)/blind_box?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&multiStatements=True",
		CMSDB:           "root:123456@tcp(127.0.0.1:3306)/chaoshe_data?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&multiStatements=True",
		Pool:            5,
		ConnMaxLifetime: 100,
	}
	if err := SetupDBEngine(); err != nil {
		panic(err)
	}

}

func TestMain(m *testing.M) {
	SetupTest()
	m.Run()
}
