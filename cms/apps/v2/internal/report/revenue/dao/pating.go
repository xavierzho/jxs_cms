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

// 参与用户
// comment: 每次运行 传入当前日期
type Pating struct {
	iDao.DailyModel
	UserCnt    uint `gorm:"column:user_cnt; default:0" json:"user_cnt"` // 参与用户数: 参与各类潮玩+集市(创建, 下单)
	UserCntNew uint `gorm:"column:user_cnt_new; default:0" json:"user_cnt_new"`
}

func (Pating) TableName() string {
	return "revenue_pating"
}

type PatingDao struct {
	*iDao.DailyModelDao[*Pating]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewPatingDao(engine, center *gorm.DB, log *logger.Logger) *PatingDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".PatingDao")))
	return &PatingDao{
		DailyModelDao: iDao.NewDailyModelDao[*Pating](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *PatingDao) Generate(cDate time.Time) (data *Pating, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	t.date,
	count(distinct t.user_id) as user_cnt,
	count(distinct case t.is_new when 1 then t.user_id else null end) as user_cnt_new
from
	(
	select distinct
		date_format(tv.created_at, '%[1]s') as date,
		(case when datediff(tv.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		tv.user_id
	from market_order tv, users u -- 集市 创建者
	where tv.created_at between '%s' and '%s' and tv.user_id = u.id and u.is_admin = 0
	union
	select distinct
		date_format(muo.created_at, '%[1]s') as date,
		(case when datediff(muo.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		muo.user_id
	from market_user_offer muo, users u -- 集市 交易者
	where muo.created_at between '%s' and '%s' and muo.user_id = u.id and u.is_admin = 0
	union
	SELECT distinct
		FROM_UNIXTIME(left(o.pay_time, 10), '%[1]s') as date,
		(case when datediff(FROM_UNIXTIME(left(o.pay_time, 10)), u.created_at) = 0 then 1 else 0 end) is_new,
		o.user_id
	FROM %[4]s o, users u -- 发货用户
	Where o.pay_time between %[5]d and %d and o.user_id = u.id and u.is_admin = 0
	union
	select distinct
		date_format(bl.created_at, '%[1]s') as date, 
		(case when datediff(bl.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		bl.user_id
	from balance_log bl, users u
	where
		bl.created_at between '%s' and '%s'
		and (bl.source_type between 100 and 199 or bl.source_type in (601)) -- 抽赏 + 商城
		and bl.update_amount <= 0
		and bl.user_id = u.id and u.is_admin = 0
	) t
group by
	t.date
	`,
		pkg.SQL_DATE_FORMAT,
		cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		"`order`",
		cDate.UnixMilli(), cDate.Add(24*time.Hour-time.Millisecond).UnixMilli(),
	)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
		return nil, err
	}

	return data, nil
}
