package dao

import (
	"context"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type StepByStepDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewStepByStepDao(center *gorm.DB, log *logger.Logger) *StepByStepDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".StepByStepDao")))
	return &StepByStepDao{
		center: center,
		logger: log,
	}
}

func (d *StepByStepDao) LogList(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, paper app.Pager) (data []map[string]any, summary map[string]any, err error) {
	err = d.center.Table("(?) as t", d.logAll(dateTimeRange, queryParams)).
		Select("count(0) as total, sum(point) as point, sum(inner_price) as inner_price").
		Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("LogList summary err: %v", err)
		return
	}

	err = d.logAll(dateTimeRange, queryParams).
		Scopes(database.Paginate(paper.Page, paper.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("LogList list err: %v", err)
		return
	}

	return
}

func (d *StepByStepDao) LogAll(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.logAll(dateTimeRange, queryParams).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("LogAll err: %v", err)
		return
	}

	return
}

func (d *StepByStepDao) logAll(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.all(dateTimeRange, queryParams).
		Select(
			"ua.params as id",
			"ua.created_at",
			"ua.user_id",
			"u.nickname as user_name",
			"JSON_UNQUOTE(ua.params_3->'$.config_id') as period",
			"sc.point_type",
			"sc.point",
			"JSON_UNQUOTE(ua.params_3->'$.step_no') as step_no",
			"JSON_UNQUOTE(ua.params_3->'$.cell_no') as cell_no",
			"ifnull(sum(case when ac.type = 20 then i.inner_price * ac.num else 0 end), 0) as inner_price",
		).
		Group("ua.params, ua.created_at, ua.user_id, u.nickname, JSON_UNQUOTE(ua.params_3->'$.config_id'), sc.point_type, sc.point, JSON_UNQUOTE(ua.params_3->'$.step_no'), JSON_UNQUOTE(ua.params_3->'$.cell_no')").
		Order("ua.created_at desc, ua.user_id")
}

func (d *StepByStepDao) Detail(queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.center.
		Select(
			"ac.type as award_type",
			"ac.value as award_value",
			"ac.name as award_name",
			"ac.params as award_params",
			"ac.num as award_num",
			"(case when ac.type = 20 then i.inner_price else 0 end) as inner_price",
		).
		Table("activity_step_by_step_award_config ac").
		Joins("left join item i on ac.type = 20 and ac.value = i.id").
		Scopes(database.ScopeQuery(queryParams)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Detail err: %v", err)
		return
	}

	return
}

func (d *StepByStepDao) DetailAll(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.all(dateTimeRange, queryParams).
		Select(
			"ua.created_at",
			"ua.user_id",
			"u.nickname as user_name",
			"JSON_UNQUOTE(ua.params_3->'$.config_id') as period",
			"JSON_UNQUOTE(ua.params_3->'$.step_no') as step_no",
			"JSON_UNQUOTE(ua.params_3->'$.cell_no') as cell_no",
			"ac.type as award_type",
			"ac.value as award_value",
			"ac.name as award_name",
			"ac.params as award_params",
			"ac.num as award_num",
			"ifnull(case when ac.type = 20 then i.inner_price else 0 end, 0) as inner_price",
		).
		Order("ua.created_at desc, ua.user_id").
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("DetailAll err: %v", err)
		return
	}

	return
}

func (d *StepByStepDao) all(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Table("activity a, user_activity ua, users u, activity_step_by_step_step_config sc, activity_step_by_step_award_config ac").
		Joins("left join item i on ac.type = 20 and ac.value = i.id").
		Where("a.key = 'StepByStep'").
		Where("a.id = ua.activity_id").
		Where("ua.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ua.user_id = u.id").
		Where("sc.config_id = cast(JSON_UNQUOTE(ua.params_3->'$.config_id') as signed)").
		Where("sc.step_no = cast(JSON_UNQUOTE(ua.params_3->'$.step_no') as signed)").
		Where("ac.cell_config_id = cast(ua.params as signed)").
		Scopes(database.ScopeQuery(queryParams))
}
