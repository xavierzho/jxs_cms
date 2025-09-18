package internal

import (
	"fmt"
	"log"
	"time"

	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/setting"
)

func InitSetting(config *setting.Config) (err error) {
	if err = config.ReadSection("Server", &global.ServerSetting); err != nil {
		return fmt.Errorf("ReadSection Server: %w", err)
	}

	if err = config.ReadSection("APP", &global.APPSetting); err != nil {
		return fmt.Errorf("ReadSection APP: %w", err)
	}
	if err = config.ReadSection("Tel", &global.TelSetting); err != nil {
		return fmt.Errorf("ReadSection Tel: %w", err)
	}

	global.APPSetting.DefaultContextTimeout *= time.Second
	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if err = pkg.SetTimeZone(global.ServerSetting.TimeZone); err != nil {
		return fmt.Errorf("global.SetTimeZone: %w", err)
	}
	if err = global.SetLanguage(global.ServerSetting.Language); err != nil {
		return fmt.Errorf("global.SetLanguage: %w", err)
	}

	return nil
}

func InitObject() (err error) {
	// 创建日志
	if err = global.SetupLogger(); err != nil {
		log.Fatalf("SetupLogger: %v", err)
	}
	if err = global.SetupAlarm(); err != nil {
		return fmt.Errorf("SetupAlarm: %w", err)
	}
	// 迁移模式仅初始化配置和日志
	if global.ServerSetting.RunMode == global.RUN_MODE_MIGRATE {
		return
	}
	// 载入i18n
	if err = global.SetupI18n(); err != nil {
		log.Fatalf("SetupI18n: %v", err)
	}
	// 验证器
	if err = global.SetupValidator(); err != nil {
		log.Fatalf("SetupValidator: %v", err)
	}
	// IP
	// if err = global.SetupIpDB(); err != nil {
	// 	log.Fatalf("SetupIpDB: %v", err)
	// }

	return
}
