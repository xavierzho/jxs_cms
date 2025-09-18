package dao_test

import (
	"context"
	"fmt"
	"testing"

	"data_backend/internal/dao"
	"data_backend/internal/global"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"
	"data_backend/pkg/redisdb"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

var log *logger.Logger
var db *gorm.DB
var rdb *redisdb.RedisClient

func TestMain(m *testing.M) {
	if err := SetupLogger(); err != nil {
		panic(err)
	}
	if err := SetupDBEngine(); err != nil {
		panic(err)
	}
	if err := SetupRedis(); err != nil {
		panic(err)
	}

	if !db.Migrator().HasTable(dao.User{}.TableName()) {
		if err := dao.InitFirstUser(db, log); err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(dao.OperationLog{}.TableName()) {
		if err := db.AutoMigrate(dao.OperationLog{}); err != nil {
			panic(err)
		}
	}

	fmt.Println("Test Run")
	m.Run()
}

func SetupLogger() (err error) {
	jackLog := &lumberjack.Logger{
		Filename:  "test-0000-00-00.log",
		MaxSize:   500,
		LocalTime: true,
	}
	log = logger.NewLogger(context.Background(), jackLog)
	return nil
}

func SetupDBEngine() (err error) {
	dbStr := "root:123456@tcp(127.0.0.1:3306)/chaoshe_data?charset=utf8mb4&parseTime=True&multiStatements=True"
	db, err = database.NewDBEngine(global.RUN_MODE_DEBUG.String(), dbStr, 100, log)
	if err != nil {
		return err
	}

	return nil
}

func SetupRedis() (err error) {
	rdb, err = redisdb.NewRedisClient("127.0.0.1:6379", "", 6)
	if err != nil {
		return err
	}

	return nil
}
