package dao

import (
	"context"
	"fmt"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type CostAward struct {
	iDao.DailyModel
	GetUserCnt          uint  `gorm:"column:get_user_cnt; type:int unsigned; default:0; comment:获得用户数" json:"get_user_cnt"`
	GetAmount           uint  `gorm:"column:get_amount; type:bigint; default:0; comment:获得总额" json:"get_amount"`
	AcceptUserCnt       uint  `gorm:"column:accept_user_cnt; type:int unsigned; default:0; comment:领取用户数" json:"accept_user_cnt"`
	AcceptAmount        uint  `gorm:"column:accept_amount; type:bigint; default:0; comment:领取总额" json:"accept_amount"`
	AwardAmount         int64 `gorm:"column:award_amount; type:bigint; default:0; comment:现金奖励总额" json:"award_amount"`
	AwardItemShowPrice  int64 `gorm:"column:award_item_show_price; type:bigint; default:0; comment:物品奖励展示价总额" json:"award_item_show_price"`
	AwardItemInnerPrice int64 `gorm:"column:award_item_inner_price; type:bigint; default:0; comment:物品奖励成本价总额" json:"award_item_inner_price"`
}

func (CostAward) TableName() string {
	return "cost_award" //获取表名
}

type CostAwardDao struct {
	*iDao.DailyModelDao[*CostAward]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewCostAwardDao(engine, center *gorm.DB, log *logger.Logger) *CostAwardDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".CostAwardDao")))
	return &CostAwardDao{
		DailyModelDao: iDao.NewDailyModelDao[*CostAward](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *CostAwardDao) Generate(cDate time.Time) (dataLog, dataAward *CostAward, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataLog, err = d.generateLog(cDate)
		return err
	})
	eg.Go(func() (err error) {
		dataAward, err = d.generateAward(cDate)
		return err
	})

	err = eg.Wait()

	return
}

func (d *CostAwardDao) generateLog(cDate time.Time) (data *CostAward, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(caul.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"count(distinct case when update_point > 0 then user_id else null end) as get_user_cnt",
			"sum(case when update_point > 0 then update_point else 0 end) as get_amount",
			"count(distinct case when log_type = 100 then user_id else null end) as accept_user_cnt",
			"sum(case when log_type = 100 then -update_point else 0 end) as accept_amount",
		).
		Table("activity_cost_award_user_log caul, users u").
		Where("caul.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("caul.user_id = u.id and u.is_admin=0").
		Group(fmt.Sprintf("date_format(caul.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateLog err: %v", err)
		return
	}

	return
}

func (d *CostAwardDao) generateAward(cDate time.Time) (data *CostAward, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(ua.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(case award_type when 0 then award_value else 0 end) as award_amount",
			"sum(case award_type when 20 then cac.award_num * ifnull(i.show_price, 0) else 0 end) as award_item_show_price",
			"sum(case award_type when 20 then cac.award_num * ifnull(i.inner_price, 0) else 0 end) as award_item_inner_price",
		).
		Table("users u, activity a, user_activity ua, activity_cost_award_config cac").
		Joins("left join item i on cac.award_type = 20 and cac.award_value = i.id").
		Where("a.name = '欧气值'").
		Where("a.id = ua.activity_id").
		Where("ua.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ua.params_3->'$.type' = '0'").
		Where("cast(ua.params as SIGNED) = cac.config_id").
		Where("ua.user_id = u.id and u.is_admin=0").
		Group(fmt.Sprintf("date_format(ua.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateAward err: %v", err)
		return
	}

	return
}
