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
type InviteBet struct {
	Date           string    `gorm:"column:date; type:varchar(10); primary_key;" json:"date" form:"date"`
	CreatedAt      time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	UserID         int64     `gorm:"column:user_id; type:bigint; primary_key;" json:"user_id"`
	UserName       string    `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	ParentUserID   int64     `gorm:"column:parent_user_id; type:bigint" json:"parent_user_id"`
	ParentUserName string    `gorm:"column:parent_user_name; type:varchar(64)" json:"parent_user_name"`
	Amount         int64     `gorm:"column:amount; type:bigint" json:"amount"`
}

func (InviteBet) TableName() string {
	return "invite_bet"
}

type InviteBetDao struct {
	*iDao.DailyModelDao[*InviteBet]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewInviteBetDao(engine, center *gorm.DB, log *logger.Logger) *InviteBetDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".InviteBetDao")))
	return &InviteBetDao{
		DailyModelDao: iDao.NewDailyModelDao[*InviteBet](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *InviteBetDao) Generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*InviteBet, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"ui.user_id", "u.nickname as user_name",
			"ui.parent_user_id", "up.nickname as parent_user_name",
			"-sum(bl.update_amount) as amount",
		).
		Table("user_invite ui, balance_log bl, users u, users up").
		Where("ui.parent_user_id <> 0").
		Where("ui.user_id = bl.user_id").
		Where("bl.source_type between 100 and 199").
		Where("ui.user_id = u.id").
		Where("u.is_admin = 0").
		Where("ui.parent_user_id = up.id").
		Where(fmt.Sprintf("bl.created_at between '%s' and '%s'", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT))).
		Scopes(database.ScopeQuery(queryParams)).
		Group(fmt.Sprintf("date_format(bl.created_at, '%s'), ui.user_id, u.nickname, ui.parent_user_id, up.nickname", pkg.SQL_DATE_FORMAT)).
		Order("`date`, amount DESC").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
	}

	return
}
