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
	TaskTypeParams database.QueryWhereGroup
	TaskKeyParams  database.QueryWhereGroup
}

type TaskList struct {
	ID                        string          `gorm:"column:id; type:varchar(255)" json:"id"`
	TaskID                    string          `gorm:"column:task_id; type:varchar(255)" json:"task_id"`
	DateTime                  string          `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID                    int64           `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName                  string          `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	TaskKey                   string          `gorm:"column:task_key; type:varchar(64)" json:"task_key"`
	TaskType                  int             `gorm:"column:task_type; type:int" json:"task_type"`
	TaskName                  string          `gorm:"column:task_name; type:varchar(64);" json:"task_name"`
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

type TaskListAwardDetail struct {
	ID                        string          `gorm:"column:id; type:varchar(255)" json:"id"`
	TaskID                    string          `gorm:"column:task_id; type:varchar(255)" json:"task_id"`
	DateTime                  string          `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID                    int64           `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName                  string          `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	TaskKey                   string          `gorm:"column:task_key; type:varchar(64)" json:"task_key"`
	TaskType                  int             `gorm:"column:task_type; type:int" json:"task_type"`
	TaskName                  string          `gorm:"column:task_name; type:varchar(64);" json:"task_name"`
	AwardType                 int64           `gorm:"column:award_type; type:bigint" json:"award_type"`
	AwardName                 string          `gorm:"column:award_name; type:varchar(64)" json:"award_name"`
	RewardValueItem           decimal.Decimal `gorm:"column:reward_value_item; " json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `gorm:"column:reward_value_cost_award_point;" json:"reward_value_cost_award_point"`
	AwardNum                  int64           `gorm:"column:award_num; type:bigint" json:"award_num"`
}

type TaskDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewTaskDao(center *gorm.DB, log *logger.Logger) *TaskDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".TaskDao")))
	return &TaskDao{
		center: center,
		logger: log,
	}
}

func (d *TaskDao) GetList(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup, pager *app.Pager) (summary map[string]any, data []*TaskList, err error) {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.TaskKeyParams...)
	queryParams = append(queryParams, paramsGroup.TaskTypeParams...)

	subQuery := d.center.
		Select(
			"l.id",
			"l.task_id",
			fmt.Sprintf("date_format(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"l.user_id",
			"u.nickname as user_name",
			"t.key AS task_key",
			"t.type AS task_type",
			"t.name AS task_name",
			"jt.award_id",
		).
		Table("task_user_log l").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN task t ON t.id = l.task_id").
		Joins(`
		LEFT JOIN JSON_TABLE(
			JSON_EXTRACT(l.params_2, '$.award_id'),
			'$[*]' COLUMNS (
				award_id BIGINT PATH '$'
			)
		) AS jt ON TRUE
	`).
		Where("l.created_at BETWEEN ? AND ?",
			dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT),
			dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT),
		).
		Scopes(database.ScopeQuery(queryParams))

	allDB := d.center.
		Table("(?) AS task", subQuery).
		Joins("LEFT JOIN task_award a ON a.id = task.award_id").
		Joins("LEFT JOIN item i ON i.id = a.value AND a.type = ?", form.AwardType_Item).
		Select("task.id,task.task_id,task.date_time,task.user_id,task.user_name,task.task_key,task.task_type,task.task_name,SUM(IF(a.type = ?, i.inner_price*a.num, 0)) AS reward_value_item,SUM(IF(a.type = ?, a.value, 0)) AS reward_value_cost_award_point", form.AwardType_Item, form.AwardType_CostAwardPoint).
		Group("task.id")
	err = allDB.
		Order("task.date_time DESC").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if pager != nil {
				return database.Paginate(pager.Page, pager.PageSize)(db)
			}
			return db
		}).
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("GetList Find: %v", err)
		return nil, nil, err
	}

	summaryDB := d.center.Table("(?) AS task", allDB).Select("COUNT(1) as total,SUM(reward_value_item) AS reward_value_item,SUM(reward_value_cost_award_point) AS reward_value_cost_award_point")
	err = summaryDB.Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("GetList Agg: %v", err)
		return nil, nil, err
	}
	return
}

func (d *TaskDao) GetAwardDetail(taskID string) (data []*AwardDetail, err error) {
	subQuery := d.center.
		Table("task_user_log AS l").
		Select("DISTINCT jt.award_id").
		Joins(`
			LEFT JOIN JSON_TABLE(
				JSON_EXTRACT(l.params_2, '$.award_id'),
				'$[*]' COLUMNS (
					award_id BIGINT PATH '$'
				)
			) AS jt ON TRUE
		`).
		Where("l.task_id = ?", taskID)
	err = d.center.
		Table("task_award AS a").
		Joins("LEFT JOIN item i ON i.id = a.value AND a.type = ?", form.AwardType_Item).
		Where("a.id IN (?)", subQuery).
		Order("a.id ASC").
		Select(`
			a.type AS award_type,
			a.name AS award_name,
			IF(a.type = ?, i.inner_price, 0) AS reward_value_item,
			IF(a.type = ?, a.value, 0) AS reward_value_cost_award_point,
			a.num AS award_num
		`, form.AwardType_Item, form.AwardType_CostAwardPoint).
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("GetAwardDetail: %v", err)
		return
	}

	return
}

func (d *TaskDao) GetListAwardDetail(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup, pager *app.Pager) (data []*TaskListAwardDetail, err error) {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.TaskKeyParams...)
	queryParams = append(queryParams, paramsGroup.TaskTypeParams...)

	subQuery := d.center.
		Select(
			"l.id",
			"l.task_id",
			fmt.Sprintf("date_format(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"l.user_id",
			"u.nickname as user_name",
			"t.key AS task_key",
			"t.type AS task_type",
			"t.name AS task_name",
			"a.type AS award_type",
			"a.name AS award_name",
			fmt.Sprintf("IF(a.type = %d, i.inner_price, 0) AS reward_value_item", form.AwardType_Item),
			fmt.Sprintf("IF(a.type = %d, a.value, 0) AS reward_value_cost_award_point", form.AwardType_CostAwardPoint),
			"a.num AS award_num",
		).
		Table("task_user_log l").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN task t ON t.id = l.task_id").
		Joins(`
			LEFT JOIN JSON_TABLE(
				JSON_EXTRACT(l.params_2, '$.award_id'),
				'$[*]' COLUMNS (
					award_id BIGINT PATH '$'
				)
			) AS jt ON TRUE
		`).
		Joins("LEFT JOIN task_award a ON a.id = jt.award_id").
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
		d.logger.Errorf("GetListAwardDetail Find: %v", err)
		return nil, err
	}

	return
}
