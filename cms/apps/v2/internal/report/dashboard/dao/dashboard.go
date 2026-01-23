package dao

import (
	"context"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Dashboard struct {
	iDao.DailyModel
	NewUserCnt           int   `gorm:"column:new_user_cnt; type:int;" json:"new_user_cnt"`
	ActiveUserCnt        int   `gorm:"column:active_user_cnt; type:int;" json:"active_user_cnt"`
	PatingUserCnt        int   `gorm:"column:pating_user_cnt; type:int;" json:"pating_user_cnt"`
	PatingUserCntNew     int   `gorm:"column:pating_user_cnt_new; type:int;" json:"pating_user_cnt_new"`
	PayUserCnt           int   `gorm:"column:pay_user_cnt; type:int;" json:"pay_user_cnt"`
	PayUserCntNew        int   `gorm:"column:pay_user_cnt_new; type:int;" json:"pay_user_cnt_new"`
	RechargeUserCnt      int   `gorm:"column:recharge_user_cnt; type:int;" json:"recharge_user_cnt"`
	RechargeUserCntNew   int   `gorm:"column:recharge_user_cnt_new; type:int;" json:"recharge_user_cnt_new"`
	RechargeAmount       int64 `gorm:"column:recharge_amount; type:bigint;" json:"recharge_amount"`
	RechargeAmountWeChat int64 `gorm:"column:recharge_amount_wechat; type:bigint;" json:"recharge_amount_wechat"`
	RechargeAmountAli    int64 `gorm:"column:recharge_amount_ali; type:bigint;" json:"recharge_amount_ali"`
	RechargeAmountHuiFu  int64 `gorm:"column:recharge_amount_huifu; type:bigint;" json:"recharge_amount_huifu"`
	DrawAmount           int64 `gorm:"column:draw_amount; type:bigint;" json:"draw_amount"`
	RefundAmountWeChat   int64 `gorm:"-" json:"refund_amount_wechat"`
	RefundAmountAli      int64 `gorm:"-" json:"refund_amount_ali"`
	RefundAmountHuiFu    int64 `gorm:"-" json:"refund_amount_huifu"`
	MarketOrderCnt       int   `gorm:"-" json:"market_order_cnt"`
	MarketAmount0        uint  `gorm:"-" json:"market_amount_0"`
	MarketAmount1        uint  `gorm:"-" json:"market_amount_1"`
	MarketAmount2        uint  `gorm:"-" json:"market_amount_2"`
}

func (Dashboard) TableName() string {
	return "dashboard"
}

type DashboardGroup struct {
	NewUser        *Dashboard
	ActiveUser     *Dashboard
	Pating         *Dashboard
	Pay            *Dashboard
	Recharge       *Dashboard
	Draw           *Dashboard
	RechargeRefund *Dashboard
	SavingRefund   *Dashboard
}

type DashboardDao struct {
	*iDao.DailyModelDao[*Dashboard]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewDashboardDao(engine, center *gorm.DB, log *logger.Logger) *DashboardDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".DashboardDao")))
	return &DashboardDao{
		DailyModelDao: iDao.NewDailyModelDao[*Dashboard](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *DashboardDao) Save(data *Dashboard) (err error) {
	if err = d.engine.Clauses(clause.OnConflict{UpdateAll: true}).
		Omit("created_at", "refund_amount_wechat", "refund_amount_ali", "refund_amount_huifu", "market_order_cnt", "market_amount_0", "market_amount_1", "market_amount_2").
		Create(data).Error; err != nil {
		d.logger.Errorf("Save: %v", err)
		return err
	}
	return nil
}

// ... (Generate function omitted for brevity as it is unchanged here, but verify context)

// ...

// 原路退款 - 分渠道
func (d *DashboardDao) generateRechargeRefund(startTime, endTime time.Time) (data *Dashboard, err error) {
	var result struct {
		DrawAmount         int64
		RefundAmountWeChat int64
		RefundAmountAli    int64
		RefundAmountHuiFu  int64
	}

	err = d.center.
		Select(
			"sum(refund_amount) as draw_amount",
			"sum(case platform_id when 'wechatapp' then refund_amount when 'wechatjs' then refund_amount else 0 end) as refund_amount_wechat",
			"sum(case platform_id when 'alipay' then refund_amount else 0 end) as refund_amount_ali",
			"sum(case platform_id when 'huifu' then refund_amount else 0 end) as refund_amount_huifu",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.role = 0").
		Where("ppo.refund_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status = 9").
		Where("ppo.pay_source_type <> 12"). // 金币储值 不算在 Refund 中
		Find(&result).Error
	if err != nil {
		d.logger.Errorf("generateRechargeRefund: %v", err)
		return nil, err
	}

	data = &Dashboard{
		DrawAmount:         result.DrawAmount,
		RefundAmountWeChat: result.RefundAmountWeChat,
		RefundAmountAli:    result.RefundAmountAli,
		RefundAmountHuiFu:  result.RefundAmountHuiFu,
	}

	return data, nil
}

// 储值退款
func (d *DashboardDao) generateSavingRefund(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select("sum(amount) as draw_amount").
		Table("refund_order_detail rod, users u").
		Where("rod.refund_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("rod.status = 3").
		Where("rod.user_id = u.id").
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateSavingRefund: %v", err)
		return nil, err
	}

	return
}
