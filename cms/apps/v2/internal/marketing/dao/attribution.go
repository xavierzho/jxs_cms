package dao

import (
	"context"
	"time"

	"data_backend/pkg/logger"
	"gorm.io/gorm"
)

type UserAttribution struct {
	ID           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserID       int64     `gorm:"column:user_id;index" json:"user_id"`
	Channel      string    `gorm:"column:channel;type:varchar(32);index" json:"channel"` // e.g., "xiaomi"
	OAID         string    `gorm:"column:oaid;type:varchar(64);index" json:"oaid"`
	IMEI         string    `gorm:"column:imei;type:varchar(64)" json:"imei"`
	CallbackURL  string    `gorm:"column:callback_url;type:varchar(512)" json:"callback_url"`
	CampaignInfo string    `gorm:"column:campaign_info;type:text" json:"campaign_info"` // JSON or string to store extra params
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;DEFAULT:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}

func (UserAttribution) TableName() string {
	return "user_attribution"
}

type AttributionDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewAttributionDao(engine *gorm.DB, log *logger.Logger) *AttributionDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".AttributionDao")))
	return &AttributionDao{
		engine: engine,
		logger: log,
	}
}

func (d *AttributionDao) Create(attr *UserAttribution) error {
	return d.engine.Create(attr).Error
}

func (d *AttributionDao) GetByUserID(userID int64) (*UserAttribution, error) {
	var attr UserAttribution
	err := d.engine.Where("user_id = ?", userID).First(&attr).Error
	if err != nil {
		return nil, err
	}
	return &attr, nil
}

func (d *AttributionDao) GetByOAID(oaid string) (*UserAttribution, error) {
	var attr UserAttribution
	err := d.engine.Where("oaid = ?", oaid).Order("created_at desc").First(&attr).Error
	if err != nil {
		return nil, err
	}
	return &attr, nil
}
