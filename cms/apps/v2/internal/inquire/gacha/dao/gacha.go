package dao

import (
	"context"

	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type GachaDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewGachaDao(center *gorm.DB, log *logger.Logger) *GachaDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".GachaDao")))
	return &GachaDao{
		center: center,
		logger: log,
	}
}
