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

func (d *CostAwardLogDao) List(dateTimeRange [2]time.Time, balanceType uint, queryParams database.QueryWhereGroup, paper app.Pager) (summary map[string]any, data []map[string]any, err error) {
	err = d.all(dateTimeRange, balanceType, queryParams).
		Select(
			"count(0) as total",
			"count(distinct bl.user_id) as user_cnt",
			"sum(bl.update_amount) as update_point",
		).
		Find(&summary).Error
	if err != nil {
		d.logger.Errorf("List summary err: %v", err)
		return
	}

	err = d.all(dateTimeRange, balanceType, queryParams).
		Select("bl.*, u.nickname as user_name").
		Order("bl.created_at desc, bl.user_id").
		Scopes(database.Paginate(paper.Page, paper.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List data err: %v", err)
		return
	}

	return
}

func (d *CostAwardLogDao) All(dateTimeRange [2]time.Time, balanceType uint, queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.all(dateTimeRange, balanceType, queryParams).
		Select("bl.*, u.nickname as user_name").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All err: %v", err)
		return
	}

	return
}

func (d *CostAwardLogDao) all(dateTimeRange [2]time.Time, balanceType uint, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Table("balance_log bl, users u").
		Where("type = 3").
		Where("bl.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("bl.user_id = u.id").
		Where("bl.type = ?", balanceType).
		Scopes(database.ScopeQuery(queryParams))
}
