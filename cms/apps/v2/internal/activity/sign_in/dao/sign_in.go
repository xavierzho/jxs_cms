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

type SignInDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewSignInDao(center *gorm.DB, log *logger.Logger) *SignInDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".SignInDao")))
	return &SignInDao{
		center: center,
		logger: log,
	}
}

func (d *SignInDao) List(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, paper app.Pager) (data []map[string]any, total int64, err error) {
	err = d.all(dateTimeRange, queryParams).
		Count(&total).
		Scopes(database.Paginate(paper.Page, paper.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List err: %v", err)
		return
	}
	return
}

func (d *SignInDao) All(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) (data []map[string]any, err error) {
	err = d.all(dateTimeRange, queryParams).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All err: %v", err)
		return
	}

	return
}

func (d *SignInDao) all(dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Select(
			"ua.created_at",
			"ua.user_id",
			"u.nickname as user_name",
			"JSON_UNQUOTE(ua.params_3->'$.type') as type",
			"JSON_UNQUOTE(ua.params_3->'$.value') as value",
			"JSON_UNQUOTE(ua.params_3->'$.day_no') as day_no",
			"JSON_UNQUOTE(ua.params_3->'$.sign_in_type') as sign_in_type",
		).
		Table("activity a, user_activity ua, users u").
		Where("a.key = 'SignIn'").
		Where("a.id = ua.activity_id").
		Where("ua.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ua.user_id = u.id").
		Scopes(database.ScopeQuery(queryParams)).
		Order("ua.created_at desc, ua.user_id")
}
