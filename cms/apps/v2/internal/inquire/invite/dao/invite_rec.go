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

type InviteRec struct {
	Date           string `gorm:"column:date; type:datetime " json:"date"`
	UserID         int64  `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName       string `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	ParentUserID   int64  `gorm:"column:parent_user_id; type:bigint" json:"parent_user_id"`
	ParentUserName string `gorm:"column:parent_user_name; type:varchar(64)" json:"parent_user_name"`
}

type InviteRecDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewInviteRecDao(center *gorm.DB, log *logger.Logger) *InviteRecDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".InviteRecDao")))
	return &InviteRecDao{
		center: center,
		logger: log,
	}
}

func (d *InviteRecDao) ListAndSummary(summaryField []string, dateRange [2]time.Time, queryParams database.QueryWhereGroup, pager app.Pager) (summary map[string]any, data []*InviteRec, err error) {
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

func (d *InviteRecDao) All(dateRange [2]time.Time, queryParams database.QueryWhereGroup) (data []*InviteRec, err error) {
	err = d.all(dateRange, queryParams).Order("`date` desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return
	}

	return
}

func (d *InviteRecDao) all(dateRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("date_format(ui.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"ui.user_id AS user_id",
			"u.nickname AS user_name", //被邀请人名称
			"ui.parent_user_id AS parent_user_id",
			"up.nickname AS parent_user_name", //邀请人名称
		).
		Table("user_invite ui").
		Joins("left join users AS u ON ui.user_id = u.id").
		Joins("left join users AS up ON ui.parent_user_id = up.id").
		Where("ui.created_at between ? and ?", dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ui.parent_user_id <> 0").
		Scopes(database.ScopeQuery(queryParams))
}
