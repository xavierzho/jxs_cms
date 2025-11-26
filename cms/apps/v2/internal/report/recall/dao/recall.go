package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

// TODO 用户名改为 从数据库中查询
type Recall struct {
	Date           string    `gorm:"column:date; type:varchar(10); primary_key;" json:"date" form:"date"`
	CreatedAt      time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	UserID         int64     `gorm:"column:user_id; type:bigint; primary_key;" json:"user_id"`
	UserName       string    `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	ParentUserID   int64     `gorm:"column:parent_user_id; type:bigint" json:"parent_user_id"`
	ParentUserName string    `gorm:"column:parent_user_name; type:varchar(64)" json:"parent_user_name"`
	Amount         int64     `gorm:"column:amount; type:bigint" json:"amount"`
	Point          int64     `gorm:"column:point; type:bigint;" json:"point"`
}

func (Recall) TableName() string {
	return "recall"
}

type RecallDao struct {
	*iDao.DailyModelDao[*Recall]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewRecallDao(engine, center *gorm.DB, log *logger.Logger) *RecallDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".RecallDao")))
	return &RecallDao{
		DailyModelDao: iDao.NewDailyModelDao[*Recall](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

// func (d *RecallDao) Generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*Recall, err error) {
// 	err = d.center.
// 		Select(
// 			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
// 			"ui.user_id", "u.nickname as user_name",
// 			"ui.parent_user_id", "up.nickname as parent_user_name",
// 			"-sum(bl.update_amount) as amount",
// 		).
// 		Table("user_recall ui").
// 		Joins("JOIN balance_log bl ON ui.user_id = bl.user_id").
// 		Joins("JOIN users u ON ui.user_id = u.id").
// 		Joins("JOIN users up ON ui.parent_user_id = up.id").
// 		Where("ui.parent_user_id <> 0").
// 		Where("bl.source_type between 100 and 199").
// 		Where("u.role = 0").
// 		Where(fmt.Sprintf("bl.created_at between '%s' and '%s'", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))).
// 		Scopes(database.ScopeQuery(queryParams)).
// 		Group(fmt.Sprintf("date_format(bl.created_at, '%s'), ui.user_id, u.nickname, ui.parent_user_id, up.nickname", pkg.SQL_DATE_FORMAT)).
// 		Order("`date`, amount DESC").
// 		Find(&data).Error
// 	if err != nil {
// 		d.logger.Errorf("Generate: %v", err)
// 	}

// 	return
// }

func (d *RecallDao) Generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*Recall, err error) {
	// 计算金额汇总
	subQuery := d.center.
		Select(
			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"ui.user_id", "u.nickname as user_name",
			"ui.parent_user_id", "up.nickname as parent_user_name",
			"-sum(bl.update_amount) as amount",
		).
		Table("user_recall ui").
		Joins("JOIN balance_log bl ON ui.user_id = bl.user_id").
		Joins("JOIN users u ON ui.user_id = u.id").
		Joins("JOIN users up ON ui.parent_user_id = up.id").
		Where("ui.parent_user_id <> 0").
		Where("bl.source_type between 100 and 199").
		Where("u.role = 0").
		Where(fmt.Sprintf("bl.created_at between '%s' and '%s'", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))).
		Scopes(database.ScopeQuery(queryParams)).
		Group(fmt.Sprintf("date_format(bl.created_at, '%s'), ui.user_id, u.nickname, ui.parent_user_id, up.nickname", pkg.SQL_DATE_FORMAT)).
		Order("`date`, amount DESC")

	// 计算积分汇总
	pointQuery := d.center.
		Select(
			fmt.Sprintf("date_format(ac.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"ac.user_id", "sum(ac.update_point) as point",
		).
		Table("activity_cost_award_user_log ac").
		Where("ac.log_type = 106").
		Where(fmt.Sprintf("ac.created_at between '%s' and '%s'", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))).
		Group(fmt.Sprintf("date_format(ac.created_at, '%s'), ac.user_id", pkg.SQL_DATE_FORMAT))

	// 合并
	err = d.center.
		Select(
			"sub.date",
			"sub.user_id",
			"sub.user_name",
			"sub.parent_user_id",
			"sub.parent_user_name",
			"sub.amount",
			"COALESCE(ac.point, 0) as point",
		).
		Table("(?) as sub", subQuery).
		Joins("LEFT JOIN (?) ac ON sub.parent_user_id = ac.user_id AND sub.date = ac.date", pointQuery).
		Find(&data).Error

	if err != nil {
		d.logger.Errorf("Generate: %v", err)
	}

	return
}
