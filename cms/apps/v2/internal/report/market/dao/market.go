package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Market struct {
	iDao.DailyModel
	UserCnt  uint `gorm:"column:user_cnt; default:0" json:"user_cnt"`
	OrderCnt uint `gorm:"column:order_cnt; default:0" json:"order_cnt"`
	Amount0  uint `gorm:"column:amount_0; default:0" json:"amount_0"`
	Amount1  uint `gorm:"column:amount_1; default:0" json:"amount_1"`
	Amount2  uint `gorm:"column:amount_2; default:0" json:"amount_2"`
}

func (Market) TableName() string {
	return "market"
}

type MarketDao struct {
	*iDao.DailyModelDao[*Market]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewMarketDao(engine, center *gorm.DB, log *logger.Logger) *MarketDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".MarketDao")))
	return &MarketDao{
		DailyModelDao: iDao.NewDailyModelDao[*Market](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *MarketDao) Generate(cDate time.Time) (dataCnt, dataAmount *Market, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataCnt, err = d.generateCnt(cDate)
		return err
	})
	eg.Go(func() (err error) {
		dataAmount, err = d.generateAmount(cDate)
		return err
	})

	err = eg.Wait()

	return
}

func (d *MarketDao) generateCnt(cDate time.Time) (data *Market, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	t.date,
	count(distinct t.user_id) as user_cnt,
	sum(order_cnt) as order_cnt
from
	(	
	select date_format(t.created_at, '%[1]s') as date, user_id, count(distinct t.id) as order_cnt
	from market_order t, users u
	where t.created_at between '%s' and '%s' and t.user_id = u.id and u.role = 0
	group by date_format(t.created_at, '%[1]s'), user_id
	union all
	select date_format(t.created_at, '%[1]s') as date, user_id, 0 as order_cnt
	from market_user_offer t, users u
	where t.created_at between '%s' and '%s' and t.user_id = u.id and u.role = 0
	group by date_format(t.created_at, '%[1]s'), user_id
	) t
group by
	t.date
	`,
		pkg.SQL_DATE_FORMAT,
		cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateCnt: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *MarketDao) generateAmount(cDate time.Time) (data *Market, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(case muou.role when 0 then bl.update_amount else 0 end) as amount_0",
			"sum(case muou.role when 1 then bl.update_amount else 0 end) as amount_1",
			"sum(case muou.role when 2 then bl.update_amount else 0 end) as amount_2",
		).
		Table("balance_log bl").
		Joins("join users u on bl.user_id = u.id and u.role = 0").
		Joins("join market_user_offer muo on bl.source_id = muo.id").
		Joins("join users muou on muo.user_id = muou.id").
		Where("bl.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("bl.source_type in (203, 204) and bl.update_amount > 0").
		Group(fmt.Sprintf("date_format(bl.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateAmount: %v", err)
		return nil, err
	}

	return data, nil
}
