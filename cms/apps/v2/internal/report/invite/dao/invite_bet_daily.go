package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type InviteBetDaily struct {
	iDao.DailyModel
	TotalAmount int64 `gorm:"column:total_amount; type:bigint" json:"total_amount"`
	Amount      int64 `gorm:"column:amount; type:bigint" json:"amount"`
	Difference  int64 `gorm:"column:difference; type:bigint" json:"difference"`
}

func (InviteBetDaily) TableName() string {
	return "invite_bet_daily"
}

type InviteBetDailyDao struct {
	*iDao.DailyModelDao[*InviteBetDaily]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewInviteBetDailyDao(engine, center *gorm.DB, log *logger.Logger) *InviteBetDailyDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".NewInviteBetDailyDao")))
	return &InviteBetDailyDao{
		DailyModelDao: iDao.NewDailyModelDao[*InviteBetDaily](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *InviteBetDailyDao) Generate(cDate time.Time) (data []*InviteBetDaily, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(-update_amount) as total_amount",
			"sum(case when ui.parent_user_id <> 0 then -update_amount else 0 end) as amount",
			"sum(case when ui.parent_user_id <> 0 then 0 else -update_amount end) as difference",
		).
		Table("users u, balance_log bl").
		Joins("left join user_invite ui on bl.user_id = ui.user_id").
		Where("bl.user_id = u.id and u.is_admin = 0").
		Where("bl.source_type BETWEEN 100 AND 199").
		Where("bl.created_at BETWEEN ? AND ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Group(fmt.Sprintf("date_format(bl.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate err: %v", err)
	}

	return
}
