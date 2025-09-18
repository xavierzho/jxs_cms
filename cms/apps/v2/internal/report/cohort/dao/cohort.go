package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cohort struct {
	iDao.DailyTypeModel
	Total             uint            `gorm:"column:total; default:0; comment:当日总数" json:"total"`
	FirstDay          uint            `gorm:"column:first_day; default:0;" json:"first_day"`
	SecondDay         uint            `gorm:"column:second_day; default:0;" json:"second_day"`
	ThirdDay          uint            `gorm:"column:third_day; default:0;" json:"third_day"`
	FourthDay         uint            `gorm:"column:fourth_day; default:0;" json:"fourth_day"`
	FifthDay          uint            `gorm:"column:fifth_day; default:0;" json:"fifth_day"`
	SixthDay          uint            `gorm:"column:sixth_day; default:0;" json:"sixth_day"`
	SeventhDay        uint            `gorm:"column:seventh_day; default:0;" json:"seventh_day"`
	FourteenthDay     uint            `gorm:"column:fourteenth_day; default:0;" json:"fourteenth_day"`
	ThirtiethDay      uint            `gorm:"column:thirtieth_day; default:0;" json:"thirtieth_day"`
	SixtiethDay       uint            `gorm:"column:sixtieth_day; default:0;" json:"sixtieth_day"`
	NinetyDay         uint            `gorm:"column:ninety_day; default:0;" json:"ninety_day"`
	No180Day          uint            `gorm:"column:no_180_day; default:0;" json:"no_180_day"`
	FirstDayRate      decimal.Decimal `gorm:"-" json:"first_day_rate"`
	SecondDayRate     decimal.Decimal `gorm:"-" json:"second_day_rate"`
	ThirdDayRate      decimal.Decimal `gorm:"-" json:"third_day_rate"`
	FourthDayRate     decimal.Decimal `gorm:"-" json:"fourth_day_rate"`
	FifthDayRate      decimal.Decimal `gorm:"-" json:"fifth_day_rate"`
	SixthDayRate      decimal.Decimal `gorm:"-" json:"sixth_day_rate"`
	SeventhDayRate    decimal.Decimal `gorm:"-" json:"seventh_day_rate"`
	FourteenthDayRate decimal.Decimal `gorm:"-" json:"fourteenth_day_rate"`
	ThirtiethDayRate  decimal.Decimal `gorm:"-" json:"thirtieth_day_rate"`
	SixtiethDayRate   decimal.Decimal `gorm:"-" json:"sixtieth_day_rate"`
	NinetyDayRate     decimal.Decimal `gorm:"-" json:"ninety_day_rate"`
	No180DayRate      decimal.Decimal `gorm:"-" json:"no_180_day_rate"`
}

func (Cohort) TableName() string {
	return "cohort"
}

// Save 时必须字段
func (Cohort) SaveBaseField() []string {
	return []string{"date", "data_type", "created_at", "updated_at", "total"}
}

type CohortDao struct {
	*iDao.DailyTypeModelDao[*Cohort]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewCohortDao(engine, center *gorm.DB, log *logger.Logger) *CohortDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".CohortDao")))
	return &CohortDao{
		DailyTypeModelDao: iDao.NewDailyTypeModelDao[*Cohort](engine, log),
		engine:            engine,
		center:            center,
		logger:            log,
	}
}

// 新增留存
// 维度: 用户注册日期 统计项: 第X天登录的用户数
func (d *CohortDao) GenerateNewUserActive(regDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	COUNT(DISTINCT(user_id))				AS total,
	COUNT(if(active_day = 0, 1, Null))		AS first_day,
	COUNT(if(active_day = 1, 1, Null))		AS second_day,
	COUNT(if(active_day = 2, 1, Null))		AS third_day,
	COUNT(if(active_day = 3, 1, Null))		AS fourth_day,
	COUNT(if(active_day = 4, 1, Null))		AS fifth_day,
	COUNT(if(active_day = 5, 1, Null))		AS sixth_day,
	COUNT(if(active_day = 6, 1, Null))		AS seventh_day,
	COUNT(if(active_day = 13, 1, Null))		AS fourteenth_day,
	COUNT(if(active_day = 29, 1, Null))		AS thirtieth_day,
	COUNT(if(active_day = 59, 1, Null))		AS sixtieth_day,
	COUNT(if(active_day = 89, 1, Null))		AS ninety_day,
	COUNT(if(active_day = 179, 1, Null))	AS no_180_day
from
	(
	select distinct
		u.id as user_id,
		date_format(u.created_at, '%s') 		as date,
		DATEDIFF(l.created_at, u.created_at)	as active_day
	from
		users u
		left join logon_logs l on l.user_id = u.id and l.created_at between '%s' and '%s'
	where
		u.created_at between '%s' and '%s'
		and u.is_admin = 0
	) t
group by
	date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		regDate.Format(pkg.DATE_TIME_MIL_FORMAT), regDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateNewUserActive: %v", err)
		return nil, err
	}

	return data, nil
}

// 有效用户
// 维度: 用户注册日期 统计项: 第X天的前7天内(含X)登录次数(一天只算一次)大于2 且 消费次数大于 0 的用户数
func (d *CohortDao) GenerateNewUserValidated(regDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	l.date,
	count(distinct l.user_id)															AS total,
	count(case when logon_count_7 > 2 and pay_count_7 > 0 then 1 else Null end)			AS first_day,
	count(case when logon_count_8 > 2 and pay_count_8 > 0 then 1 else Null end)			AS second_day,
	count(case when logon_count_9 > 2 and pay_count_9 > 0 then 1 else Null end)			AS third_day,
	count(case when logon_count_10 > 2 and pay_count_10 > 0 then 1 else Null end)		AS fourth_day,
	count(case when logon_count_11 > 2 and pay_count_11 > 0 then 1 else Null end)		AS fifth_day,
	count(case when logon_count_12 > 2 and pay_count_12 > 0 then 1 else Null end)		AS sixth_day,
	count(case when logon_count_13 > 2 and pay_count_13 > 0 then 1 else Null end)		AS seventh_day,
	count(case when logon_count_14 > 2 and pay_count_14 > 0 then 1 else Null end)		AS fourteenth_day,
	count(case when logon_count_30 > 2 and pay_count_30 > 0 then 1 else Null end)		AS thirtieth_day,
	count(case when logon_count_60 > 2 and pay_count_60 > 0 then 1 else Null end)		AS sixtieth_day,
	count(case when logon_count_90 > 2 and pay_count_90 > 0 then 1 else Null end)		AS ninety_day,
	count(case when logon_count_180 > 2 and pay_count_180 > 0 then 1 else Null end)		AS one_hundred_and_eighty_day
from
	(
	select
		user_id,
		date,
		COUNT(if(active_day between 0 and 6, active_day, Null))			AS logon_count_7,
		COUNT(if(active_day between 1 and 7, active_day, Null))			AS logon_count_8,
		COUNT(if(active_day between 2 and 8, active_day, Null))			AS logon_count_9,
		COUNT(if(active_day between 3 and 9, active_day, Null))			AS logon_count_10,
		COUNT(if(active_day between 4 and 10, active_day, Null))		AS logon_count_11,
		COUNT(if(active_day between 5 and 11, active_day, Null))		AS logon_count_12,
		COUNT(if(active_day between 6 and 12, active_day, Null))		AS logon_count_13,
		COUNT(if(active_day between 7 and 13, active_day, Null))		AS logon_count_14,
		COUNT(if(active_day between 23 and 29, active_day, Null))		AS logon_count_30,
		COUNT(if(active_day between 53 and 59, active_day, Null))		AS logon_count_60,
		COUNT(if(active_day between 83 and 89, active_day, Null))		AS logon_count_90,
		COUNT(if(active_day between 173 and 179, active_day, Null))		AS logon_count_180
	from
		(
		select distinct
			u.id as user_id,
			date_format(u.created_at, '%[1]s') as date,
			DATEDIFF(l.created_at, u.created_at) as active_day
		from
			users u
			left join logon_logs l on l.user_id = u.id and l.created_at between '%s' and '%s'
		where
			u.created_at between '%s' and '%s'
			and u.is_admin = 0
		) l
	group by
		user_id, date
	) l,
	(
	select distinct
		u.id as user_id,
		date_format(u.created_at, '%[1]s') as date,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 0 and 6, bl.created_at, Null))		AS pay_count_7,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 1 and 7, bl.created_at, Null))		AS pay_count_8,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 2 and 8, bl.created_at, Null))		AS pay_count_9,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 3 and 9, bl.created_at, Null))		AS pay_count_10,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 4 and 10, bl.created_at, Null))		AS pay_count_11,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 5 and 11, bl.created_at, Null))		AS pay_count_12,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 6 and 12, bl.created_at, Null))		AS pay_count_13,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 7 and 13, bl.created_at, Null))		AS pay_count_14,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 23 and 29, bl.created_at, Null))		AS pay_count_30,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 53 and 59, bl.created_at, Null))		AS pay_count_60,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 83 and 89, bl.created_at, Null))		AS pay_count_90,
		COUNT(if(DATEDIFF(bl.created_at, u.created_at) between 173 and 179, bl.created_at, Null))	AS pay_count_180
	from
		users u
		left join balance_log bl on bl.user_id = u.id and bl.created_at between '%s' and '%s'
			and (bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0
	where
		u.created_at between '%s' and '%s'
		and u.is_admin = 0
	group by
		u.id,
		date_format(u.created_at, '%[1]s')
	) p
Where
	l.user_id = p.user_id
group by
	l.date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		regDate.Format(pkg.DATE_TIME_MIL_FORMAT), regDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateNewUserValidated: %v", err)
		return nil, err
	}

	return data, nil
}

// 新增消费
// 维度: 用户注册日期 统计项: 第X天消费的用户数(消费:潮玩+集市+充值)
func (d *CohortDao) GenerateNewUserConsume(regDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	count(distinct user_id)				AS total,
	COUNT(IF(pay_day = 0, 1, Null))		AS first_day,
	COUNT(IF(pay_day = 1, 1, Null))		AS second_day,
	COUNT(IF(pay_day = 2, 1, Null))		AS third_day,
	COUNT(IF(pay_day = 3, 1, Null))		AS fourth_day,
	COUNT(IF(pay_day = 4, 1, Null))		AS fifth_day,
	COUNT(IF(pay_day = 5, 1, Null))		AS sixth_day,
	COUNT(IF(pay_day = 6, 1, Null))		AS seventh_day,
	COUNT(IF(pay_day = 13, 1, Null))	AS fourteenth_day,
	COUNT(IF(pay_day = 29, 1, Null))	AS thirtieth_day,
	COUNT(IF(pay_day = 59, 1, Null))	AS sixtieth_day,
	COUNT(IF(pay_day = 89, 1, Null))	AS ninety_day,
	COUNT(IF(pay_day = 179, 1, Null))	AS no_180_day
from
	(
	select distinct
		u.id as user_id,
		date_format(u.created_at, '%s') 		as date,
		DATEDIFF(bl.finish_at, u.created_at) 	as pay_day
	from
		users u
		left join balance_log bl on bl.user_id = u.id and bl.finish_at between '%s' and '%s'
			and (bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0
	where
		u.created_at between '%s' and '%s'
		and u.is_admin = 0
	) t
group by 
	date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		regDate.Format(pkg.DATE_TIME_MIL_FORMAT), regDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateNewUserConsume: %v", err)
		return nil, err
	}

	return data, nil
}

// 参与留存
// 维度: 参与潮玩/集市/发货日期 统计项: 参与后第X天的登录用户数
func (d *CohortDao) GeneratePatingUserActive(patingDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
	select
		t.date,
		count(distinct t.user_id)													AS total,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 0, t.user_id, Null))		AS first_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 1, t.user_id, Null))		AS second_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 2, t.user_id, Null))		AS third_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 3, t.user_id, Null))		AS fourth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 4, t.user_id, Null))		AS fifth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 5, t.user_id, Null))		AS sixth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 6, t.user_id, Null))		AS seventh_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 13, t.user_id, Null))	AS fourteenth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 29, t.user_id, Null))	AS thirtieth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 59, t.user_id, Null))	AS sixtieth_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 89, t.user_id, Null))	AS ninety_day,
		COUNT(distinct IF(datediff(l.created_at, t.date) = 179, t.user_id, Null))	AS no_180_day
	from
		(
		select distinct date_format(tv.created_at, '%[1]s') as date, tv.user_id
		from market_order tv, users u -- 集市 创建者
		where tv.created_at between '%s' and '%s' and tv.user_id = u.id and u.is_admin = 0
		union
		select distinct date_format(muo.created_at, '%[1]s') as date, muo.user_id
		from market_user_offer muo, users u -- 集市 交易者
		where muo.created_at between '%s' and '%s' and muo.user_id = u.id and u.is_admin = 0
		union
		SELECT distinct FROM_UNIXTIME(left(o.pay_time, 10), '%[1]s') as date, o.user_id
		FROM %[6]s o, users u -- 发货用户
		Where o.pay_time between %[7]d and %d and o.user_id = u.id and u.is_admin = 0
		union
		select distinct date_format(bl.created_at, '%[1]s') as date, bl.user_id
		from balance_log bl, users u -- bet 用户
		where
			bl.created_at between '%s' and '%s'
			and (bl.source_type between 100 and 199 or bl.source_type in (601)) and bl.update_amount <= 0
			and bl.user_id = u.id and u.is_admin = 0
		) t
		left join logon_logs l on l.user_id = t.user_id and l.created_at between '%s' and '%s'
	group by
		t.date
		`,
		pkg.SQL_DATE_FORMAT,
		patingDate.Format(pkg.DATE_TIME_MIL_FORMAT), patingDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		"`order`",
		patingDate.UnixMilli(), patingDate.Add(24*time.Hour-time.Millisecond).UnixMilli(),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GeneratePatingUserActive: %v", err)
		return nil, err
	}

	return data, nil
}

// 消费留存
// 维度: 消费日期 统计项: 消费后第X天的登录用户数
func (d *CohortDao) GenerateConsumeUserActive(payDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	COUNT(DISTINCT(user_id))				AS total,
	COUNT(if(active_day = 0, 1, Null))		AS first_day,
	COUNT(if(active_day = 1, 1, Null))		AS second_day,
	COUNT(if(active_day = 2, 1, Null))		AS third_day,
	COUNT(if(active_day = 3, 1, Null))		AS fourth_day,
	COUNT(if(active_day = 4, 1, Null))		AS fifth_day,
	COUNT(if(active_day = 5, 1, Null))		AS sixth_day,
	COUNT(if(active_day = 6, 1, Null))		AS seventh_day,
	COUNT(if(active_day = 13, 1, Null))		AS fourteenth_day,
	COUNT(if(active_day = 29, 1, Null))		AS thirtieth_day,
	COUNT(if(active_day = 59, 1, Null))		AS sixtieth_day,
	COUNT(if(active_day = 89, 1, Null))		AS ninety_day,
	COUNT(if(active_day = 179, 1, Null))	AS no_180_day
from
	(
	select distinct
		bl.user_id,
		date_format(bl.finish_at, '%s') 		as date,
		DATEDIFF(l.created_at, bl.finish_at) as active_day
	from
		balance_log bl, users u
		left join logon_logs l on l.user_id = u.id and l.created_at between '%s' and '%s'
	where
		bl.finish_at between '%s' and '%s'
		and (bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0
		and bl.user_id = u.id and u.is_admin = 0
	) t
group by
	date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		payDate.Format(pkg.DATE_TIME_MIL_FORMAT), payDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateConsumeUserActive: %v", err)
		return nil, err
	}

	return data, nil
}

// TODO 受邀留存
// 维度: 受邀用户注册日期 统计项: 第X天登录的用户数
func (d *CohortDao) GenerateInvitedNewUserActive(regDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	COUNT(DISTINCT(user_id))				AS total,
	COUNT(if(active_day = 0, 1, Null))		AS first_day,
	COUNT(if(active_day = 1, 1, Null))		AS second_day,
	COUNT(if(active_day = 2, 1, Null))		AS third_day,
	COUNT(if(active_day = 3, 1, Null))		AS fourth_day,
	COUNT(if(active_day = 4, 1, Null))		AS fifth_day,
	COUNT(if(active_day = 5, 1, Null))		AS sixth_day,
	COUNT(if(active_day = 6, 1, Null))		AS seventh_day,
	COUNT(if(active_day = 13, 1, Null))		AS fourteenth_day,
	COUNT(if(active_day = 29, 1, Null))		AS thirtieth_day,
	COUNT(if(active_day = 59, 1, Null))		AS sixtieth_day,
	COUNT(if(active_day = 89, 1, Null))		AS ninety_day,
	COUNT(if(active_day = 179, 1, Null))	AS no_180_day
from
	(
	select distinct
		u.id as user_id,
		date_format(u.created_at, '%s') 		as date,
		DATEDIFF(l.created_at, u.created_at)	as active_day
	from
		task_record_view ru,
		users u
		left join logon_logs l on l.user_id = u.id and l.created_at between '%s' and '%s'
	where
		u.created_at between '%s' and '%s' and u.id = ru.task_obj_id and ru.task_type = 1
		and u.is_admin = 0
	) t
group by
	date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		regDate.Format(pkg.DATE_TIME_MIL_FORMAT), regDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateInvitedNewUserActive: %v", err)
		return nil, err
	}

	return data, nil
}

// TODO 受邀新增消费
// 维度: 受邀用户注册日期 统计项: 第X天消费的用户数(消费:潮玩+集市+充值)
func (d *CohortDao) GenerateInvitedNewUserConsume(regDate time.Time, startUpdateDate, lastUpdateDate time.Time) (data *Cohort, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	date,
	count(distinct user_id)				AS total,
	COUNT(IF(pay_day = 0, 1, Null))		AS first_day,
	COUNT(IF(pay_day = 1, 1, Null))		AS second_day,
	COUNT(IF(pay_day = 2, 1, Null))		AS third_day,
	COUNT(IF(pay_day = 3, 1, Null))		AS fourth_day,
	COUNT(IF(pay_day = 4, 1, Null))		AS fifth_day,
	COUNT(IF(pay_day = 5, 1, Null))		AS sixth_day,
	COUNT(IF(pay_day = 6, 1, Null))		AS seventh_day,
	COUNT(IF(pay_day = 13, 1, Null))	AS fourteenth_day,
	COUNT(IF(pay_day = 29, 1, Null))	AS thirtieth_day,
	COUNT(IF(pay_day = 59, 1, Null))	AS sixtieth_day,
	COUNT(IF(pay_day = 89, 1, Null))	AS ninety_day,
	COUNT(IF(pay_day = 179, 1, Null))	AS no_180_day
from
	(
	select distinct
		u.id as user_id,
		date_format(u.created_at, '%s') 		as date,
		DATEDIFF(bl.created_at, u.created_at) 	as pay_day
	from
		task_record_view ru,
		users u
		left join balance_log bl on bl.user_id = u.id and bl.created_at between '%s' and '%s' and bl.payed = 1 and bl.after_sale_status = 0 and bl.pay_price > 0 and bl.biz_type in (1,2,3,4,6)
		-- 1-一番赏 2-盲盒机 3-集市-交易 4-潮玩赏 5-发货订单  6-集市-换购
	where
		u.created_at between '%s' and '%s' and u.id = ru.task_obj_id and ru.task_type = 1
		and u.is_admin = 0
	) t
group by 
	date
	`,
		pkg.SQL_DATE_FORMAT,
		startUpdateDate.Format(pkg.DATE_TIME_MIL_FORMAT), lastUpdateDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		regDate.Format(pkg.DATE_TIME_MIL_FORMAT), regDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("GenerateInvitedNewUserConsume: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *CohortDao) Save(updateField []string, data ...*Cohort) (err error) {
	if len(data) == 0 {
		return nil
	}
	err = d.engine.
		Model(data).
		Select(append(data[0].SaveBaseField(), updateField...)).
		Save(data).Error
	if err != nil {
		d.logger.Errorf("Save: %v", err)
		return err
	}

	return nil
}

// 总注册用户数
func (d *CohortDao) GetNewUserCnt(dateRange [2]time.Time, queryParams []database.QueryWhere) (count int64, err error) {
	err = d.center.
		Select("count(distinct u.id)").
		Table("users u").
		Where("u.created_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Scopes(database.ScopeQuery(queryParams)).
		Scan(&count).Error // 新增注册
	if err != nil {
		d.logger.Errorf("GetNewUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}

// 注册用户数
func (d *CohortDao) GetNewUserCntList(dateRange [2]time.Time, queryParams []database.QueryWhere) (data []map[string]interface{}, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(u.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"count(distinct u.id) as user_cnt", // 新增注册
		).
		Table("users u").
		Where("u.created_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Scopes(database.ScopeQuery(queryParams)).
		Group(fmt.Sprintf("date_format(u.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetNewUserCntList: %v", err)
		return nil, err
	}

	return data, nil
}

// 参与用户数
func (d *CohortDao) GetPatingUserCnt(dateRange [2]time.Time, queryParams database.QueryWhereGroup) (count int64, err error) {
	whereSql, sqlParam := queryParams.GetQuerySqlParams()
	if whereSql != "" {
		whereSql = " and " + whereSql
	}
	sqlParam = append(sqlParam, sqlParam...)

	err = d.center.Raw(fmt.Sprintf(`
select
	count(distinct t.user_id) as user_cnt
from
	(
	select distinct tv.user_id
	from market_order tv, users u
	where tv.created_at between '%[1]s' and '%s' %s and tv.user_id = u.id and u.is_admin = 0
	union
	select distinct muo.user_id
	from market_user_offer muo, users u
	where muo.created_at between '%[1]s' and '%s' %s and muo.user_id = u.id and u.is_admin = 0
	union
	SELECT distinct o.user_id
	FROM %s o, users u -- 发货用户
	Where o.pay_time between %[5]d and %d %[3]s and o.user_id = u.id and u.is_admin = 0
	union
	select distinct bl.user_id
	from balance_log bl, users u
	where
		bl.created_at between '%[1]s' and '%s' %s
		and (bl.source_type between 100 and 199 or bl.source_type in (601)) and bl.update_amount <= 0
		and bl.user_id = u.id and u.is_admin = 0
	) t
	`,
		dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		whereSql,
		"`order`",
		dateRange[0].UnixMilli(), dateRange[1].Add(24*time.Hour-time.Millisecond).UnixMilli(),
	), sqlParam...).
		Scan(&count).Error
	if err != nil {
		d.logger.Errorf("GetPatingUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}

// 付费用户数
func (d *CohortDao) GetPayUserCnt(dateRange [2]time.Time, queryParams []database.QueryWhere) (count int64, err error) {
	err = d.center.
		Select("count(distinct bl.user_id) as user_cnt").
		Table("balance_log bl, users u").
		Where("bl.finish_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("(bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0").
		Where("bl.user_id = u.id").
		Where("u.is_admin = 0").
		Scopes(database.ScopeQuery(queryParams)).
		Scan(&count).Error
	if err != nil {
		d.logger.Errorf("GetPayUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}

// TODO 受邀用户数
func (d *CohortDao) GetInvitedUserCnt(dateRange [2]time.Time, queryParams []database.QueryWhere) (count int64, err error) {
	err = d.center.
		Select("count(distinct ru.task_obj_id)").
		Table("task_record_view ru, users u").
		Where("u.id = ru.task_obj_id").
		Where("ru.task_type = 1").
		Where("ru.created_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Scopes(database.ScopeQuery(queryParams)).
		Scan(&count).Error // 新增注册
	if err != nil {
		d.logger.Errorf("GetInvitedUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}
