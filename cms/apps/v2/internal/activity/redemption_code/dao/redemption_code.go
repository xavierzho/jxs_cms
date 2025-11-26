package dao

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AllRequestParamsGroup struct {
	DateTimeParams database.QueryWhereGroup
	UsersParams    database.QueryWhereGroup
	NameParams     database.QueryWhereGroup
	CodeParams     database.QueryWhereGroup
}

type RedemptionCodeLog struct {
	LogID                     string          `gorm:"column:log_id; type:varchar(255)" json:"log_id"`
	RedemptionCodeID          string          `gorm:"column:redemption_code_id; type:varchar(255)" json:"redemption_code_id"`
	Code                      string          `gorm:"column:code; type:varchar(32)" json:"code"`
	Name                      string          `gorm:"column:name; type:varchar(32)" json:"name"`
	DateTime                  string          `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID                    int64           `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName                  string          `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	RewardValueItem           decimal.Decimal `gorm:"column:reward_value_item; type:int" json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `gorm:"column:reward_value_cost_award_point; type:int" json:"reward_value_cost_award_point"`
}

type AwardDetail struct {
	AwardType                 int64           `gorm:"column:award_type; type:bigint" json:"award_type"`
	AwardName                 string          `gorm:"column:award_name; type:varchar(64)" json:"award_name"`
	RewardValueItem           decimal.Decimal `gorm:"column:reward_value_item; " json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `gorm:"column:reward_value_cost_award_point;" json:"reward_value_cost_award_point"`
	AwardNum                  int64           `gorm:"column:award_num; type:bigint" json:"award_num"`
}

type RedemptionCodeLogAwardDetail struct {
	LogID                     string          `gorm:"column:log_id; type:varchar(255)" json:"log_id"`
	RedemptionCodeID          string          `gorm:"column:redemption_code_id; type:varchar(255)" json:"redemption_code_id"`
	Code                      string          `gorm:"column:code; type:varchar(32)" json:"code"`
	Name                      string          `gorm:"column:name; type:varchar(32)" json:"name"`
	DateTime                  string          `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID                    int64           `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName                  string          `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	AwardType                 int64           `gorm:"column:award_type; type:bigint" json:"award_type"`
	AwardName                 string          `gorm:"column:award_name; type:varchar(64)" json:"award_name"`
	RewardValueItem           decimal.Decimal `gorm:"column:reward_value_item; " json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `gorm:"column:reward_value_cost_award_point;" json:"reward_value_cost_award_point"`
	AwardNum                  int64           `gorm:"column:award_num; type:bigint" json:"award_num"`
}

type RedemptionCodeDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewRedemptionCodeDao(center *gorm.DB, log *logger.Logger) *RedemptionCodeDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".RedemptionCodeDao")))
	return &RedemptionCodeDao{
		center: center,
		logger: log,
	}
}

func (d *RedemptionCodeDao) GetLog(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup, pager *app.Pager) (summary map[string]any, data []*RedemptionCodeLog, err error) {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.NameParams...)
	queryParams = append(queryParams, paramsGroup.CodeParams...)

	subQuery := d.center.
		Select(
			"l.id as log_id",
			"l.config_id as redemption_code_id",
			"c.code",
			"c.name",
			fmt.Sprintf("date_format(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"l.user_id",
			"u.nickname as user_name",
		).
		Table("activity_redemption_code_log l").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN activity_redemption_code c ON c.id = l.config_id").
		Where("l.created_at BETWEEN ? AND ?",
			dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT),
			dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		).
		Scopes(database.ScopeQuery(queryParams))

	allDB := d.center.
		Table("(?) AS redemption_code", subQuery).
		Joins("LEFT JOIN activity_redemption_code_award_log a ON a.config_id = redemption_code.log_id AND a.deleted_at IS NULL").
		Joins("LEFT JOIN item i ON i.id = a.value AND a.type = ?", form.AwardType_Item).
		Select("redemption_code.*,SUM(IF(a.type = ?, i.inner_price*a.num, 0)) AS reward_value_item,SUM(IF(a.type = ?, a.value, 0)) AS reward_value_cost_award_point", form.AwardType_Item, form.AwardType_CostAwardPoint).
		Group("redemption_code.log_id")
	err = allDB.
		Order("date_time DESC").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if pager != nil {
				return database.Paginate(pager.Page, pager.PageSize)(db)
			}
			return db
		}).
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("GetLog Find: %v", err)
		return nil, nil, err
	}

	summaryDB := d.center.Table("(?) AS redemption_code", allDB).Select("COUNT(1) as total,SUM(reward_value_item) AS reward_value_item,SUM(reward_value_cost_award_point) AS reward_value_cost_award_point")
	err = summaryDB.Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("GetList Agg: %v", err)
		return nil, nil, err
	}
	return
}

func (d *RedemptionCodeDao) GetAwardDetail(logID string) (data []*AwardDetail, err error) {
	err = d.center.Table("activity_redemption_code_award_log as a").
		Joins("LEFT JOIN item i ON i.id = a.value AND a.type = ?", form.AwardType_Item).
		Where("a.deleted_at IS NULL").Where("a.config_id = ?", logID).Order("a.id asc").Select("`type` as award_type,a.`name` as award_name,IF(a.type = ?, i.inner_price, 0) AS reward_value_item,IF(a.type = ?, a.value, 0) AS reward_value_cost_award_point,`num` as award_num", form.AwardType_Item, form.AwardType_CostAwardPoint).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetAwardDetail: %v", err)
		return
	}
	return
}

func (d *RedemptionCodeDao) GetLogAwardDetail(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup, pager *app.Pager) (data []*RedemptionCodeLogAwardDetail, err error) {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.NameParams...)
	queryParams = append(queryParams, paramsGroup.CodeParams...)

	subQuery := d.center.
		Select(
			"l.config_id as redemption_code_id",
			"c.code",
			"c.name",
			fmt.Sprintf("date_format(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"l.user_id",
			"u.nickname as user_name",
			"a.type AS award_type",
			"a.name AS award_name",
			fmt.Sprintf("IF(a.type = %d, i.inner_price, 0) AS reward_value_item", form.AwardType_Item),
			fmt.Sprintf("IF(a.type = %d, a.value, 0) AS reward_value_cost_award_point", form.AwardType_CostAwardPoint),
			"a.num AS award_num",
		).
		Table("activity_redemption_code_log l").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN activity_redemption_code c ON c.id = l.config_id").
		Joins("LEFT JOIN activity_redemption_code_award_log a ON a.config_id = l.id AND a.deleted_at IS NULL").
		Joins("LEFT JOIN item i ON i.id = a.value AND a.type = ?", form.AwardType_Item).
		Where("l.created_at BETWEEN ? AND ?",
			dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT),
			dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		).
		Scopes(database.ScopeQuery(queryParams))

	err = subQuery.
		Order("date_time DESC").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if pager != nil {
				return database.Paginate(pager.Page, pager.PageSize)(db)
			}
			return db
		}).
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("GetLogAwardDetail Find: %v", err)
		return nil, err
	}

	return
}
