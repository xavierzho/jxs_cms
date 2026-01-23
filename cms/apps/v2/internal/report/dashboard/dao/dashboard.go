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
	if err = d.engine.Clauses(clause.OnConflict{UpdateAll: true}).Omit("created_at").Create(data).Error; err != nil {
		d.logger.Errorf("Save: %v", err)
		return err
	}
	return nil
}

func (d *DashboardDao) Generate(startTime, endTime time.Time) (dataGroup DashboardGroup, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataGroup.NewUser, err = d.generateNewUserCnt(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.ActiveUser, err = d.generateActiveUserCnt(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.Pating, err = d.generatePating(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.Pay, err = d.generatePay(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.Recharge, err = d.generateRecharge(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.Draw, err = d.generateDraw(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.RechargeRefund, err = d.generateRechargeRefund(startTime, endTime)
		return err
	})
	eg.Go(func() (err error) {
		dataGroup.SavingRefund, err = d.generateSavingRefund(startTime, endTime)
		return err
	})

	err = eg.Wait()

	return
}

// 注册用户数
func (d *DashboardDao) generateNewUserCnt(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select("count(distinct u.id) as new_user_cnt").
		Table("users u").
		Where("u.created_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateNewUserCnt: %v", err)
		return nil, err
	}

	return
}

// 活跃用户数
func (d *DashboardDao) generateActiveUserCnt(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select("count(distinct l.user_id) as active_user_cnt").
		Table("logon_logs l, users u").
		Where("l.created_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("l.user_id = u.id").
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateActiveUserCnt: %v", err)
		return nil, err
	}

	return
}

// 参与用户 分新旧 // ! 当时间范围 大于一天时 pating_user_cnt pating_user_cnt_new 会重复计算一个用户 // 该函数结果不用于summary部分
func (d *DashboardDao) generatePating(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.Raw(fmt.Sprintf(`
select
	count(distinct t.user_id) as pating_user_cnt,
	count(distinct case t.is_new when 1 then t.user_id else null end) as pating_user_cnt_new
from
	(
	select distinct
		(case when datediff(tv.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		tv.user_id
	from market_order tv, users u -- 集市 创建者
	where tv.created_at between '%[1]s' and '%s' and tv.user_id = u.id and u.role = 0
	union
	select distinct
		(case when datediff(muo.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		muo.user_id
	from market_user_offer muo, users u -- 集市 交易者
	where muo.created_at between '%[1]s' and '%s' and muo.user_id = u.id and u.role = 0
	union
	SELECT distinct
		(case when datediff(FROM_UNIXTIME(left(o.pay_time, 10)), u.created_at) = 0 then 1 else 0 end) is_new,
		o.user_id
	FROM %s o, users u -- 发货用户
	Where o.pay_time between %[4]d and %d and o.user_id = u.id and u.role = 0
	union
	select distinct
		(case when datediff(bl.created_at, u.created_at) = 0 then 1 else 0 end) is_new,
		bl.user_id
	from balance_log bl, users u
	where
		bl.created_at between '%[1]s' and '%s'
		and (bl.source_type between 100 and 199 or bl.source_type in (601)) -- 抽赏 + 商城
		and bl.update_amount <= 0
		and bl.user_id = u.id and u.role = 0
	) t
	`,
		startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT),
		"`order`",
		startTime.UnixMilli(), endTime.UnixMilli(),
	)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generatePating: %v", err)
		return nil, err
	}

	return data, nil
}

// 付费用户 分新旧
func (d *DashboardDao) generatePay(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select(
			"count(distinct bl.user_id) as pay_user_cnt",
			"count(distinct case when datediff(bl.finish_at, u.created_at) = 0 then bl.user_id else null end) pay_user_cnt_new",
		).
		Table("balance_log bl, users u").
		Where("bl.finish_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("(bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0").
		Where("bl.user_id = u.id").
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generatePay: %v", err)
		return nil, err
	}

	return data, nil
}

// 充值 分渠道 // 包括储值
func (d *DashboardDao) generateRecharge(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select(
			"count(distinct ppo.user_id) as recharge_user_cnt",
			"count(distinct case when datediff(ppo.finish_time, u.created_at) = 0 then ppo.user_id else null end) recharge_user_cnt_new",
			"sum(ppo.amount) as recharge_amount",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.amount when 'wechatjs' then ppo.amount else 0 end) as  recharge_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.amount else 0 end) as  recharge_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then ppo.amount else 0 end) as  recharge_amount_huifu",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.role = 0").
		Where("ppo.finish_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		// Where("ppo.pay_source_type IN (100,201,202)"). -- 所有充值都计算
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRecharge: %v", err)
		return nil, err
	}

	return
}

// 退款(￥)
func (d *DashboardDao) generateDraw(startTime, endTime time.Time) (data *Dashboard, err error) {
	err = d.center.
		Select(
			"cast(sum(pdo.amount) as UNSIGNED) as draw_amount",
		).
		Table("pay_payout_order pdo, users u").
		Where("pdo.finish_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("pdo.state in (6, 12)").
		Where("pdo.user_id = u.id").
		Where("u.role = 0").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateDraw: %v", err)
		return nil, err
	}

	return data, nil
}

// 原路退款 - 分渠道
func (d *DashboardDao) generateRechargeRefund(startTime, endTime time.Time) (data *Dashboard, err error) {
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
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRechargeRefund: %v", err)
		return nil, err
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
