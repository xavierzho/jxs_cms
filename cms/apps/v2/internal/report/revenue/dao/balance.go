package dao

import (
	"context"

	iDao "data_backend/internal/dao"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

// 营收汇总数据 单位: 分
// comment: 每次运行 不传入日期
type Balance struct {
	iDao.DailyModel
	Wallet   uint `gorm:"column:wallet; default:0" json:"wallet"`     // 用户钱包余额
	Merchant uint `gorm:"column:merchant; default:0" json:"merchant"` // 用户商户余额
}

func (Balance) TableName() string {
	return "revenue_balance"
}

type BalanceDao struct {
	*iDao.DailyModelDao[*Balance]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewBalanceDao(engine, center *gorm.DB, log *logger.Logger) *BalanceDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".BalanceDao")))
	return &BalanceDao{
		DailyModelDao: iDao.NewDailyModelDao[*Balance](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *BalanceDao) Generate() (data *Balance, err error) {
	err = d.center.
		Select("sum(balance) as wallet").
		Table("wallet w, users u").
		Where("w.user_id = u.id").
		Where("u.is_admin = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
		return nil, err
	}

	return data, nil
}
