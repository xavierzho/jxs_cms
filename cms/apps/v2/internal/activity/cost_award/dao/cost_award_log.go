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

type CostAwardLogDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewCostAwardLogDao(center *gorm.DB, log *logger.Logger) *CostAwardLogDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".CostAwardLogDao")))
	return &CostAwardLogDao{
		center: center,
		logger: log,
	}
}

func (d *CostAwardLogDao) List(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, paper app.Pager) (summary map[string]any, data []map[string]any, err error) {
	err = d.all(dateTimeRange, queryParams).
		Select(
			"count(0) as total",
			"count(distinct caul.user_id) as user_cnt",
			"sum(caul.update_point) as update_point",
		).
		Find(&summary).Error
	if err != nil {
		d.logger.Errorf("List summary err: %v", err)
		return
	}

	err = d.all(dateTimeRange, queryParams).
		Select("caul.*, u.nickname as user_name").
		Order("caul.created_at desc, caul.user_id").
		Scopes(database.Paginate(paper.Page, paper.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List data err: %v", err)
		return
	}

	return
}

func (d *CostAwardLogDao) All(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.all(dateTimeRange, queryParams).
		Select("caul.*, u.nickname as user_name").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All err: %v", err)
		return
	}

	return
}

func (d *CostAwardLogDao) all(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Table("activity_cost_award_user_log as caul, users u").
		Where("caul.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("caul.user_id = u.id").
		Scopes(database.ScopeQuery(queryParams))
}
