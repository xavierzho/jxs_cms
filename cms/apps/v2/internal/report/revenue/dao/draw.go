package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

// 退款(￥)数据 单位: 分
// comment: 每次运行 传入 当前日期
type Draw struct {
	iDao.DailyModel
	Amount  uint `gorm:"column:amount; default:0" json:"amount"`
	UserCnt uint `gorm:"column:user_cnt; default:0" json:"user_cnt"`
	Tax     uint `gorm:"column:tax; default:0" json:"tax"` // 总抽水
	TaxNew  uint `gorm:"column:tax_new; default:0" json:"tax_new"`
	TaxOld  uint `gorm:"column:tax_old; default:0" json:"tax_old"`
}

func (Draw) TableName() string {
	return "revenue_draw"
}

type DrawDao struct {
	*iDao.DailyModelDao[*Draw]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewDrawDao(engine, center *gorm.DB, log *logger.Logger) *DrawDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".DrawDao")))
	return &DrawDao{
		DailyModelDao: iDao.NewDailyModelDao[*Draw](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *DrawDao) Generate(cDate time.Time) (data *Draw, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(pdo.finish_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"cast(sum(pdo.amount) as UNSIGNED) as amount",
			"count(distinct pdo.user_id) as user_cnt",
		).
		Table("pay_payout_order pdo, users u").
		Where("pdo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("pdo.state in (6, 12)").
		Where("pdo.user_id = u.id").
		Where("u.role = 0").
		Group(fmt.Sprintf("date_format(pdo.finish_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
		return nil, err
	}

	return data, nil
}

// draw union all savingRefund
func (d *DrawDao) Generate2(cDate time.Time) (data *Draw, err error) {
	pdoDB := d.center.
		Select(
			fmt.Sprintf("date_format(pdo.finish_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"pdo.amount",
			"pdo.user_id",
		).
		Table("pay_payout_order pdo, users u").
		Where("pdo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("pdo.state in (6, 12)").
		Where("pdo.user_id = u.id").
		Where("u.role = 0")

	rodDB := d.center.
		Select(
			fmt.Sprintf("date_format(rod.refund_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"rod.amount",
			"rod.user_id",
		).
		Table("refund_order_detail rod, users u").
		Where("rod.refund_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("rod.status = 3").
		Where("rod.user_id = u.id").
		Where("u.role = 0")

	err = d.center.
		Select(
			"date",
			"cast(sum(t.amount) as UNSIGNED) as amount",
			"count(distinct t.user_id) as user_cnt",
		).
		Table("(? union all ?) t", pdoDB, rodDB).
		Group("date").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDrawData: %v", err)
		return nil, err
	}

	return data, nil
}
