package local

import (
	"data_backend/internal/global"
	"data_backend/pkg/database"
	"data_backend/pkg/redisdb"

	"gorm.io/gorm"
)

type DatabaseConfig struct {
	CenterDB        string
	CMSDB           string
	Pool            int
	ConnMaxLifetime int
}

type RedisConfig struct {
	Host     string
	Password string
	CMSDB    int
}

var (
	DatabaseSetting *DatabaseConfig
	CenterDB        *gorm.DB
	CMSDB           *gorm.DB

	RedisClient      *redisdb.RedisClient
	RedisSetting     *RedisConfig
	MigrateModelList []any
)

func SetupDBEngine() (err error) {
	CenterDB, err = database.NewDBEngine(global.ServerSetting.RunMode.String(), DatabaseSetting.CenterDB, DatabaseSetting.ConnMaxLifetime, Logger)
	if err != nil {
		return err
	}
	CMSDB, err = database.NewDBEngine(global.ServerSetting.RunMode.String(), DatabaseSetting.CMSDB, DatabaseSetting.ConnMaxLifetime, Logger)
	if err != nil {
		return err
	}

	return nil
}

func SetupRedis() (err error) {
	RedisClient, err = redisdb.NewRedisClient(RedisSetting.Host, RedisSetting.Password, RedisSetting.CMSDB)
	if err != nil {
		return err
	}

	return nil
}

func MigrateModel() error {
	return CMSDB.AutoMigrate(MigrateModelList...)
}
