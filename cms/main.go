package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"data_backend/apps"
	"data_backend/internal"
	"data_backend/internal/global"
	"data_backend/pkg/encrypt/aes"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"
	"data_backend/pkg/setting"
	"data_backend/pkg/token"
	"data_backend/pkg/util"
)

var (
	port        string
	runMode     string
	configPath  string
	storagePath string
	config      *setting.Config
)

func init() {
	var err error
	// 设置 flag
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式: release, debug, migrate")
	flag.StringVar(&configPath, "configPath", "configs", "指定要使用的配置文件路径")
	flag.StringVar(&storagePath, "storagePath", "storage", "指定要使用的本地文件路径")
	flag.Parse()

	// 载入配置文件
	if err = SetupGlobalSetting(port, runMode, configPath); err != nil {
		log.Fatal(err)
	}

	// InitObject
	if err = InitObject(); err != nil {
		log.Fatal(err)
	}
}

func initSetting() (err error) {
	// 设置金额精度
	util.SetPrecision(util.DECIMAL_HUNDRED)

	if err = internal.InitSetting(config); err != nil {
		return fmt.Errorf("internal.InitSetting: %w", err)
	}
	if err = config.ReadSection("JWT", &token.JWTSetting); err != nil {
		return fmt.Errorf("ReadSection JWT: %w", err)
	}
	token.JWTSetting.Expire *= time.Hour
	if err = config.ReadSection("SecretKey", &aes.SecretKeySetting); err != nil {
		return fmt.Errorf("ReadSection SecretKey: %w", err)
	}
	if err = config.ReadSection("Logger", &logger.LoggerSetting); err != nil {
		return fmt.Errorf("ReadSection Logger: %w", err)
	}
	if err = config.ReadSection("Telegram", &message.TelegramSetting); err != nil {
		return fmt.Errorf("ReadSection Telegram: %w", err)
	}
	if err = config.ReadSection("WeChat", &message.WeChatSetting); err != nil {
		return fmt.Errorf("ReadSection WeChat: %w", err)
	}
	if err = config.ReadSection("SMS", &message.SMSSetting); err != nil {
		return fmt.Errorf("ReadSection SMS: %w", err)
	}

	if err = apps.InitSetting(config); err != nil {
		return fmt.Errorf("apps.InitSetting: %w", err)
	}

	return nil
}

func SetupGlobalSetting(port, runMode, configPath string) (err error) {
	if config, err = setting.NewSetting(strings.Split(configPath, ",")...); err != nil {
		return err
	}
	if err = initSetting(); err != nil {
		return err
	}

	if port != "" {
		global.ServerSetting.HttpPort = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = global.RunMode(runMode)
	}
	global.StoragePath = storagePath

	config.WatchSettingChange(initSetting, func(str string) {
		if global.Alarm != nil {
			global.Alarm.SendMsg(str, message.CmsId)
		} else {
			log.Fatal(str)
		}
	})

	return nil
}

func InitObject() (err error) {
	if err = internal.InitObject(); err != nil {
		return fmt.Errorf("internal.InitObject: %w", err)
	}
	if err = apps.InitObject(); err != nil {
		return fmt.Errorf("apps.InitObject: %w", err)
	}

	return nil
}

func main() {
	// 仅迁移
	if global.ServerSetting.RunMode == global.RUN_MODE_MIGRATE {
		err := apps.MigrateModel()
		if err != nil {
			log.Fatalf("main.AutoMigrate: %s", err.Error())
		}
		return
	}

	// 初始化apps 并 获取路由
	router, err := apps.InitAPPsRouter()
	if err != nil {
		log.Fatalf("apps.Initial: %v", err)
	}
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.Alarm.SendMsg(fmt.Sprintf("CMS ListenAndServe error: %v", err), int32(message.CmsId))
			log.Fatalf("s.ListenAndServe: %v", err)
		}
	}()
	if global.ServerSetting.RunMode == global.RUN_MODE_RELEASE {
		global.Alarm.SendMsg("Start CMS", message.CmsId)
	}

	// 重定向
	f, err := os.OpenFile("panic.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	redirectStderr(f)

	// 退出处理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
