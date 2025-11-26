package dao

import (
	"context"
	"fmt"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type Turntable struct {
	Date         string `gorm:"column:date; type:datetime " json:"date"`
	UserID       int64  `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName     string `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	Period       int64  `gorm:"column:period; type:bigint" json:"period"`
	Name         string `gorm:"column:name; type:varchar(64)" json:"name"`
	Type         int64  `gorm:"column:type; type:bigint" json:"type"`
	ItemId       int64  `gorm:"column:item_id; type:bigint" json:"item_id"`
	ItemName     string `gorm:"column:item_name; type:varchar(64)" json:"item_name"`
	PointType    string `gorm:"column:point_type; type:varchar(64)" json:"point_type"` //抽奖支付类型  10 现金点  11 欧气值
	Point        int64  `gorm:"column:point; type:bigint" json:"point"`                //抽奖消耗
	PrizeValue   int64  `gorm:"column:prize_value; type:bigint" json:"prize_value"`    //奖品价值
	OldPointType int64  `gorm:"column:old_point_type; type:bigint" json:"old_point_type"`
	OldPoint     int64  `gorm:"column:old_point; type:bigint" json:"old_point"` //抽奖消耗
}

type TurntableDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewTurntableDao(center *gorm.DB, log *logger.Logger) *TurntableDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".TurntableDao")))
	return &TurntableDao{
		center: center,
		logger: log,
	}
}

func (d *TurntableDao) ListAndSummary(summaryField []string, dateRange [2]time.Time, queryParams database.QueryWhereGroup, pager app.Pager) (summary map[string]any, data []*Turntable, err error) {
	summary = make(map[string]any)
	err = d.all(dateRange, queryParams).Order("`date` desc").Scopes(database.Paginate(pager.Page, pager.PageSize)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("ListAndSummary list: %v", err)
		return
	}

	err = d.all(dateRange, queryParams).Select(summaryField).Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("ListAndSummary summary: %v", err)
		return
	}

	return
}

func (d *TurntableDao) All(dateRange [2]time.Time, queryParams database.QueryWhereGroup) (data []*Turntable, err error) {
	err = d.all(dateRange, queryParams).Order("`date` desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return
	}

	return
}

// PrizeWheel
func (d *TurntableDao) all(dateRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("date_format(apwh.created_at, '%s') as date", pkg.SQL_DATE_TIME_FORMAT),
			"apwh.user_id AS user_id", //用户ID
			"u.nickname AS user_name", //用户昵称
			"apwc.period AS period",   //期数
			"apwc.name AS name",       //转盘名称
			"apwac.value as  item_id", //商品ID
			"apwac.type as type",      //奖品类型
			"apwac.name as item_name", //奖品名称
			"JSON_UNQUOTE(apwh.params_3->'$.point_type') as point_type",            //抽奖资源类型
			"(case when apwh.params_2 > 0 then apwh.params_2 else 0 end) as point", //抽奖消耗
			"apwc.point_type as  old_point_type",                                   //兼容旧数据，抽奖资源类型
			"apwc.point as old_point",                                              //兼容旧数据，抽奖消耗
			"i.inner_price as prize_value",                                         //奖品价值
		).
		Table("user_activity apwh").
		Joins("left join activity AS a ON apwh.activity_id = a.id").
		Joins("left join users AS u ON apwh.user_id = u.id").
		Joins("left join activity_prize_wheel_config AS apwc ON apwh.params = apwc.id").
		Joins("left join activity_prize_wheel_award_config AS apwac ON JSON_UNQUOTE(apwh.params_3->'$.award_id') = apwac.id").
		Joins("left join item AS i ON apwac.value = i.id").
		Where("a.key='PrizeWheel'").
		Where("apwh.created_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Scopes(database.ScopeQuery(queryParams))
}
