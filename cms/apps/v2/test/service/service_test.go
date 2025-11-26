package service_test

import (
	"fmt"
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/internal/global"
	"data_backend/pkg/logger"
)

func SetupTest() {
	global.ServerSetting = &global.ServerConfig{
		RunMode: global.RUN_MODE_DEBUG,
	}

	global.APPSetting = &global.APPConfig{
		DefaultPageSize:       50,
		MaxPageSize:           100,
		DefaultContextTimeout: time.Minute,
	}

	logger.LoggerSetting = &logger.LoggerConfig{
		LogSavePath: ".",
		LogFileName: "test",
		LogFileExt:  ".log",
	}
	if err := local.SetupLogger(); err != nil {
		panic(err)
	}
	if err := global.SetupLogger(); err != nil {
		panic(err)
	}

	local.DatabaseSetting = &local.DatabaseConfig{
		CMSDB:           "root:123456@tcp(127.0.0.1:3306)/chaoshe_data?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&multiStatements=True",
		CenterDB:        "root:123456@tcp(192.168.8.31:3306)/blind_box?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai&multiStatements=True",
		Pool:            5,
		ConnMaxLifetime: 100,
	}
	if err := local.SetupDBEngine(); err != nil {
		panic(err)
	}

	local.RedisSetting = &local.RedisConfig{
		Host:  "127.0.0.1:6379",
		CMSDB: 6,
	}
	if err := local.SetupRedis(); err != nil {
		panic(err)
	}

	global.StoragePath = "../../../../cmd/cms_backend/storage"
	if err := global.SetupI18n(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	SetupTest()

	fmt.Println("Test Run")
	m.Run()
}
