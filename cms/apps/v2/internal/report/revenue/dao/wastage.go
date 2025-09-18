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

// 流失数据
// 更新历史数据
// comment: 每次运行 传入第前7天, 日期作为注册日期
type Wastage struct {
	iDao.DailyModel
	WastageCnt1 uint `gorm:"column:wastage_cnt_1; default:0" json:"wastage_cnt_1"` // 流失数 注册日起2-7天均未登录
	WastageCnt3 uint `gorm:"column:wastage_cnt_3; default:0" json:"wastage_cnt_3"` // 流失数 注册日起2-3天登录过, 4-7天均未登录
}

func (Wastage) TableName() string {
	return "revenue_wastage"
}

type WastageDao struct {
	*iDao.DailyModelDao[*Wastage]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewWastageDao(engine, center *gorm.DB, log *logger.Logger) *WastageDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".WastageDao")))
	return &WastageDao{
		DailyModelDao: iDao.NewDailyModelDao[*Wastage](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *WastageDao) Generate(cDate time.Time) (data *Wastage, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	t.date,
	count(case when t.login_date_2_7 = 0 then t.user_id else null end) as wastage_cnt_1,
	count(case when t.login_date_2_3 > 0 and t.login_date_4_7 = 0 then t.user_id else null end) as wastage_cnt_3
from
	(
	select
		u.id as user_id,
		date_format(u.created_at, '%s') as date,
		count(distinct case when datediff(l.created_at, u.created_at) between 1 and 6 then date_format(l.created_at, '%[1]s') else null end) as login_date_2_7,
		count(distinct case when datediff(l.created_at, u.created_at) between 1 and 2 then date_format(l.created_at, '%[1]s') else null end) as login_date_2_3,
		count(distinct case when datediff(l.created_at, u.created_at) between 3 and 6 then date_format(l.created_at, '%[1]s') else null end) as login_date_4_7
	from
		users u
	left join logon_logs l on l.user_id = u.id and l.created_at between '%s' and '%s'
	where
		u.created_at between '%s' and '%s'
		and u.is_admin = 0
	group by
		u.id, date_format(u.created_at, '%s')
	) t
group by
	t.date
	`,
		pkg.SQL_DATE_FORMAT,
		cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(7*24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		pkg.SQL_DATE_FORMAT,
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
		return nil, err
	}

	return data, nil
}
