package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type CouponAction uint8

const (
	CouponAction_Gain   CouponAction = 0 // 获得
	CouponAction_Used   CouponAction = 1 // 使用
	CouponAction_Expire CouponAction = 2 // 过期
)

type Coupon struct {
	DateTime   string       `gorm:"column:date_time; type:varchar(19); comment:时间" json:"date_time"`
	CouponID   int64        `gorm:"column:coupon_id; type:bigint; comment:优惠券id" json:"coupon_id"`
	CouponType int8         `gorm:"column:coupon_type; type:tinyint; comment:优惠券类型" json:"coupon_type"`
	CouponName string       `gorm:"column:coupon_name; type:varchar(64); comment:优惠券名称" json:"coupon_name"`
	UserID     int64        `gorm:"column:user_id; type:bigint; comment:用户id" json:"user_id"`
	UserName   string       `gorm:"column:user_name; type:varchar(64); comment:用户昵称" json:"user_name"`
	Action     CouponAction `gorm:"column:action; type:tinyint unsigned; comment:行为" json:"action"`
	Amount     int64        `gorm:"column:amount; type:bigint; comment:抵扣金额" json:"amount"`
}

var baseSelectField = []string{
	"uc.coupon_id",
	"c.type as coupon_type",
	"c.name as coupon_name",
	"uc.user_id",
	"u.nickname as user_name",
}

type Explain struct {
	Gain   bool
	Used   bool
	Expire bool
}

type CouponDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewCouponDao(center *gorm.DB, log *logger.Logger) *CouponDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".CouponDao")))
	return &CouponDao{
		center: center,
		logger: log,
	}
}

func (d *CouponDao) List(dateTimeRange [2]time.Time, explain Explain, queryParams database.QueryWhereGroup, pager *app.Pager) (data []*Coupon, count int64, err error) {
	if err = d.allDB(d.center, dateTimeRange, explain, queryParams).Count(&count).Scopes(database.Paginate(pager.Page, pager.PageSize)).Find(&data).Error; err != nil {
		d.logger.Errorf("List: %v", err)
		return
	}

	return
}

func (d *CouponDao) All(dateTimeRange [2]time.Time, explain Explain, queryParams database.QueryWhereGroup) (data []*Coupon, err error) {
	if err = d.allDB(d.center, dateTimeRange, explain, queryParams).Find(&data).Error; err != nil {
		d.logger.Errorf("All: %v", err)
		return
	}

	return
}

func (d *CouponDao) allDB(tx *gorm.DB, dateTimeRange [2]time.Time, explain Explain, queryParams database.QueryWhereGroup) *gorm.DB {
	var dbList []any
	if explain.Gain {
		dbList = append(dbList, d.allGain(tx, dateTimeRange, queryParams))
	}
	if explain.Used {
		dbList = append(dbList, d.allUsed(tx, dateTimeRange, queryParams))
	}
	if explain.Expire {
		dbList = append(dbList, d.allExpire(tx, dateTimeRange, queryParams))
	}

	sqlList := make([]string, len(dbList))
	for ind := range sqlList {
		sqlList[ind] = "?"
	}
	return tx.Table("("+strings.Join(sqlList, " union all ")+") as t", dbList...).Order("date_time desc, coupon_id, user_id")
}

func (d *CouponDao) allGain(tx *gorm.DB, dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	selectField := append(baseSelectField, []string{
		fmt.Sprintf("date_format(uc.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
		"0 as action",
		"0 as amount",
	}...)
	return d.baseDB(tx, queryParams).
		Select(selectField).
		Where("uc.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))
}

func (d *CouponDao) allUsed(tx *gorm.DB, dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	selectField := append(baseSelectField, []string{
		fmt.Sprintf("date_format(uc.used_time, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
		"1 as action",
		"uc.deduction_amount as amount",
	}...)
	return d.baseDB(tx, queryParams).
		Select(selectField).
		Where("uc.used_time between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))
}

func (d *CouponDao) allExpire(tx *gorm.DB, dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	selectField := append(baseSelectField, []string{
		fmt.Sprintf("date_format(uc.expired_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
		"2 as action",
		"0 as amount",
	}...)
	return d.baseDB(tx, queryParams).
		Select(selectField).
		Where("uc.state = 4").
		Where("uc.expired_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))
}

func (d *CouponDao) baseDB(tx *gorm.DB, queryParams database.QueryWhereGroup) *gorm.DB {
	return tx.
		Table("users u, user_coupon uc").
		Joins("left join coupon c on uc.coupon_id = c.id").
		Where("uc.user_id = u.id").
		Scopes(database.ScopeQuery(queryParams))
}
