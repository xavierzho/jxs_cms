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
	Wallet   uint `gorm:"column:wallet; default:0" json:"wallet"`     // 吉祥币余额 (0)
	Merchant uint `gorm:"column:merchant; default:0" json:"merchant"` // 邀请返佣余额 (1)
	Gold     uint `gorm:"column:gold; default:0" json:"gold"`         // 余额余额 (2)
	Jidou    uint `gorm:"-" json:"jidou"`                             // 吉豆余额 (3)
	Point    uint `gorm:"column:point; default:0" json:"point"`       // 积分余额 (10)
	Loyalty  uint `gorm:"column:loyalty; default:0" json:"loyalty"`   // 吉祥值余额 (11)
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
		Select(
			"ifnull(sum(case `type` when 0 then balance else 0 end), 0) as wallet",
			"ifnull(sum(case `type` when 1 then balance else 0 end), 0) as merchant",
			"ifnull(sum(case `type` when 2 then balance else 0 end), 0) as gold",
			"ifnull(sum(case `type` when 3 then balance else 0 end), 0) as jidou",
			"ifnull(sum(case `type` when 10 then balance else 0 end), 0) as point",
			"ifnull(sum(case `type` when 11 then balance else 0 end), 0) as loyalty",
		).
		Table("wallet w, users u").
		Where("w.user_id = u.id").
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
		return nil, err
	}

	return data, nil
}
