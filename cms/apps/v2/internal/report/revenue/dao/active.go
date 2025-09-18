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

// 活跃用户
// * 登录表一天会有多条
// comment: 每次运行 传入 当前日期
type Active struct {
	iDao.DailyModel
	ActivateCnt    uint `gorm:"column:activate_cnt; default:0" json:"activate_cnt"`         // 打开设备数
	ActivateCntNew uint `gorm:"column:activate_cnt_new; default:0" json:"activate_cnt_new"` // 非注册用户打开app设备数
	RegisterCnt    uint `gorm:"column:register_cnt; default:0" json:"register_cnt"`         // 新注册用户数
	ActiveCnt      uint `gorm:"column:active_cnt; default:0" json:"active_cnt"`             // 日活;活跃用户数;登录用户数
	ActiveCntNew   uint `gorm:"column:active_cnt_new; default:0" json:"active_cnt_new"`
	ActiveCntOld   uint `gorm:"column:active_cnt_old; default:0" json:"active_cnt_old"`
	MaxOnlineCnt   uint `gorm:"column:max_online_cnt; default:0" json:"max_online_cnt"`   // 最大在线用户数
	ValidatedCnt7  uint `gorm:"column:validated_cnt_7; default:0" json:"validated_cnt_7"` // 有效用户数, 注册日起7天内登录2天及以上且消费过的用户数
	ValidatedCnt15 uint `gorm:"column:Validated_cnt_15; default:0" json:"Validated_cnt_15"`
}

func (Active) TableName() string {
	return "revenue_active"
}

type ActiveDao struct {
	*iDao.DailyModelDao[*Active]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewActiveDao(engine, center *gorm.DB, log *logger.Logger) *ActiveDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".ActiveDao")))
	return &ActiveDao{
		DailyModelDao: iDao.NewDailyModelDao[*Active](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *ActiveDao) Generate(cDate time.Time) (dataRegister, dataActive, dataValidated *Active, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataRegister, err = d.generateRegister(cDate)
		return err
	})

	eg.Go(func() (err error) {
		dataActive, err = d.generateActive(cDate)
		return err
	})

	eg.Go(func() (err error) {
		dataValidated, err = d.generateValidated(cDate)
		return err
	})

	err = eg.Wait()

	return
}

// 日期为注册日期
func (d *ActiveDao) generateRegister(cDate time.Time) (data *Active, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(u.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"count(distinct id) as register_cnt",
		).
		Table("users u").
		Where("u.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Group(fmt.Sprintf("date_format(u.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRegister: %v", err)
		return nil, err
	}

	return data, nil
}

// 日期为登录日期
func (d *ActiveDao) generateActive(cDate time.Time) (data *Active, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(l.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"count(distinct u.id) as active_cnt",
			"count(distinct (case when datediff(l.created_at, u.created_at) = 0 then u.id else null end)) as active_cnt_new",
			"count(distinct (case when datediff(l.created_at, u.created_at) <> 0 then u.id else null end)) as active_cnt_old",
		).
		Table("logon_logs l").
		Joins("join users u on l.user_id = u.id and u.is_admin = 0").
		Where("l.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Group(fmt.Sprintf("date_format(l.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateActive: %v", err)
		return nil, err
	}

	return data, nil
}

// 日期为登录日期
// 当天的前7/15天内(含当天) 登录大于2次且消费过 的用户数
func (d *ActiveDao) generateValidated(cDate time.Time) (data *Active, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	count(distinct case when cnt_7 > 2 then t.user_id else null end) as validated_cnt_7,
	count(distinct case when cnt_15 > 2 then t.user_id else null end) as validated_cnt_15
from
	(
	select
		l.user_id,
		l.date,
		count(case when datediff(l.date, ll.login_date) <=6 and datediff(l.date, l.pay_date) <=6 then ll.login_date else null end) as cnt_7,
		count(case when datediff(l.date, ll.login_date) <=14 and datediff(l.date, l.pay_date) <=14 then ll.login_date else null end) as cnt_15
	from
		(
		select
			l.user_id,
			date_format(l.created_at, '%s') as date,
			date_format(max(bl.created_at), '%s') as pay_date
		from
			logon_logs l,
			balance_log bl,
			users u
		where
			l.created_at between '%s' and '%s'
			and l.user_id = bl.user_id
			and bl.created_at between '%s' and '%s'
			and (bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601))
			and bl.update_amount <= 0
			and l.user_id = u.id
			and u.is_admin = 0
		group by
			l.user_id, date_format(l.created_at, '%s')
		) l,
		(
		select distinct
			l.user_id,	
			date_format(l.created_at, '%s') as login_date
		from
			logon_logs l,
			users u
		where
			l.created_at between '%s' and '%s'
			and l.user_id = u.id
			and u.is_admin = 0
		) ll
	where
		l.user_id = ll.user_id
		and datediff(l.date, ll.login_date) between 0 and 14
	group by
		l.user_id, l.date
	) t
group by
	t.date
`,
		pkg.SQL_DATE_FORMAT, pkg.SQL_DATE_FORMAT,
		cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT), // 第 x 天
		cDate.Add(-14*24*time.Hour).Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT), // 第 x-15 ~ x 天
		pkg.SQL_DATE_FORMAT,
		pkg.SQL_DATE_FORMAT,
		cDate.Add(-14*24*time.Hour).Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT), // 第 x-15 ~ x 天
	)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateActive: %v", err)
		return nil, err
	}

	return data, nil
}
