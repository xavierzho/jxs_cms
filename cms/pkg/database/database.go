package database

import (
	"log"
	"os"
	"time"

	"data_backend/pkg"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func NewDBEngine(RunMode string, connStr string, ConnMaxLifetime int, writer gormLogger.Writer) (*gorm.DB, error) {
	config := &gorm.Config{}
	var logLevel gormLogger.LogLevel
	var Colorful bool
	switch RunMode {
	case gin.ReleaseMode:
		logLevel = gormLogger.Warn
		Colorful = false
	default:
		logLevel = gormLogger.Info
		Colorful = true
		writer = log.New(os.Stdout, "\r\n", log.LstdFlags) // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
	}
	config.Logger = New(
		writer,
		gormLogger.Config{
			SlowThreshold:             500 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  Colorful,
		},
	)
	config.CreateBatchSize = 1000
	config.DisableForeignKeyConstraintWhenMigrating = true
	config.NowFunc = func() time.Time {
		return time.Now().In(pkg.Location)
	}

	masterDB := mysql.Open(connStr)
	db, err := gorm.Open(masterDB, config)
	if err != nil {
		return nil, err
	}

	// err = db.Use(dbresolver.Register(dbresolver.Config{
	// 	Sources: []gorm.Dialector{masterDB},
	// }))
	// if err != nil {
	// 	return nil, err
	// }

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetime) * time.Second)

	return db, nil
}

func NewGameDBEngine(RunMode string, connStr string, poolSize int) (*gorm.DB, error) {
	config := &gorm.Config{}
	var logLevel gormLogger.LogLevel
	if RunMode == gin.ReleaseMode {
		logLevel = gormLogger.Warn
	} else {
		logLevel = gormLogger.Info
	}
	config.Logger = New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		gormLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)
	config.CreateBatchSize = 1000
	config.DisableForeignKeyConstraintWhenMigrating = true
	config.NowFunc = func() time.Time {
		return time.Now().In(pkg.Location)
	}

	masterDB := mysql.Open(connStr)
	db, err := gorm.Open(masterDB, config)
	if err != nil {
		return nil, err
	}

	// err = db.Use(dbresolver.Register(dbresolver.Config{
	// 	Sources: []gorm.Dialector{masterDB},
	// }))

	return db, err
}
