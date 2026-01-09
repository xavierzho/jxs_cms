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

// 付费数据 单位: 分
// comment: 每次运行 传入 当前日期
type Pay struct {
	iDao.DailyModel
	Amount                     uint `gorm:"column:amount; default:0" json:"amount"` // 总消费金额: 潮玩+集市+运费+商城
	AmountBet                  uint `gorm:"column:amount_bet; default:0" json:"amount_bet"`
	AmountNew                  uint `gorm:"column:amount_new; default:0" json:"amount_new"`
	AmountOld                  uint `gorm:"column:amount_old; default:0" json:"amount_old"`
	UserCnt                    uint `gorm:"column:user_cnt; default:0" json:"user_cnt"` // 总消费人数
	UserCntNew                 uint `gorm:"column:user_cnt_new; default:0" json:"user_cnt_new"`
	UserCntOld                 uint `gorm:"column:user_cnt_old; default:0" json:"user_cnt_old"`
	UserCntFirst               uint `gorm:"column:user_cnt_first; default:0" json:"user_cnt_first"`                               // 首次消费人数
	RefundAmount               uint `gorm:"column:refund_amount; default:0" json:"refund_amount"`                                 // 退款(潮币)金额
	RefundUserCnt              uint `gorm:"column:refund_user_cnt; default:0" json:"refund_user_cnt"`                             // 退款(潮币)人数
	RechargeAmount             uint `gorm:"column:recharge_amount; default:0" json:"recharge_amount"`                             // 充值金额
	RechargeAmountWeChat       uint `gorm:"column:recharge_amount_wechat; default:0" json:"recharge_amount_wechat"`               // 充值金额(微信)
	RechargeAmountAli          uint `gorm:"column:recharge_amount_ali; default:0" json:"recharge_amount_ali"`                     // 充值金额(支付宝)
	RechargeAmountHuiFu        uint `gorm:"column:recharge_amount_huifu; default:0" json:"recharge_amount_huifu"`                 // 充值金额(汇付)
	RechargeRefundAmount       uint `gorm:"column:recharge_refund_amount; default:0" json:"recharge_refund_amount"`               // 充值退款金额
	RechargeRefundAmountWeChat uint `gorm:"column:recharge_refund_amount_wechat; default:0" json:"recharge_refund_amount_wechat"` // 充值退款金额(微信)
	RechargeRefundAmountAli    uint `gorm:"column:recharge_refund_amount_ali; default:0" json:"recharge_refund_amount_ali"`       // 充值退款金额(支付宝)
	RechargeRefundAmountHuiFu  uint `gorm:"column:recharge_refund_amount_huifu; default:0" json:"recharge_refund_amount_huifu"`   // 充值退款金额(汇付)
	DiscountAmount             uint `gorm:"column:discount_amount; default:0" json:"discount_amount"`                             // 总折扣金额
	DiscountAmountWeChat       uint `gorm:"column:discount_amount_wechat; default:0" json:"discount_amount_wechat"`               // 总折扣金额(微信)
	DiscountAmountAli          uint `gorm:"column:discount_amount_ali; default:0" json:"discount_amount_ali"`                     // 总折扣金额(支付宝)
	DiscountAmountHuiFu        uint `gorm:"column:discount_amount_huifu; default:0" json:"discount_amount_huifu"`                 // 总折扣金额(汇付)
	SavingAmount               uint `gorm:"column:saving_amount; default:0" json:"saving_amount"`                                 // 储值金额
	SavingAmountWeChat         uint `gorm:"column:saving_amount_wechat; default:0" json:"saving_amount_wechat"`                   // 储值金额(微信)
	SavingAmountAli            uint `gorm:"column:saving_amount_ali; default:0" json:"saving_amount_ali"`                         // 储值金额(支付宝)
	SavingAmountHuiFu          uint `gorm:"column:saving_amount_huifu; default:0" json:"saving_amount_huifu"`                     // 储值金额(汇付)
	SavingRefundAmount         uint `gorm:"column:saving_refund_amount; default:0" json:"saving_refund_amount"`                   // 储值退款金额
	SavingRefundAmountWeChat   uint `gorm:"column:saving_refund_amount_wechat; default:0" json:"saving_refund_amount_wechat"`     // 储值退款金额(微信)
	SavingRefundAmountAli      uint `gorm:"column:saving_refund_amount_ali; default:0" json:"SavingRefundAmountAli"`              // 储值退款金额(支付宝)
	SavingRefundAmountHuiFu    uint `gorm:"column:saving_refund_amount_huifu; default:0" json:"saving_refund_amount_huifu"`       // 储值退款金额(汇付)
}

func (Pay) TableName() string {
	return "revenue_pay"
}

type PayGroup struct {
	Pay            *Pay
	Refund         *Pay
	Recharge       *Pay
	RechargeRefund *Pay
	Saving         *Pay
	SavingRefund   *Pay
}

type PayDao struct {
	*iDao.DailyModelDao[*Pay]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewPayDao(engine, center *gorm.DB, log *logger.Logger) *PayDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".PayDao")))
	return &PayDao{
		DailyModelDao: iDao.NewDailyModelDao[*Pay](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *PayDao) Generate(cDate time.Time) (payGroup PayGroup, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		payGroup.Pay, err = d.generatePay(cDate)
		return err
	})

	eg.Go(func() (err error) {
		payGroup.Refund, err = d.generateRefund(cDate)
		return err
	})

	eg.Go(func() (err error) {
		payGroup.Recharge, err = d.generateRecharge(cDate)
		return err
	})

	eg.Go(func() (err error) {
		payGroup.RechargeRefund, err = d.generateRechargeRefund(cDate)
		return err
	})

	eg.Go(func() (err error) {
		payGroup.Saving, err = d.generateSaving(cDate)
		return err
	})

	eg.Go(func() (err error) {
		payGroup.SavingRefund, err = d.generateSavingRefund(cDate)
		return err
	})

	err = eg.Wait()

	return
}

// 需要限制 bl.update_amount <= 0; 退款的订单直接完成了，而消费的订单只有在成功时才会有完成时间
func (d *PayDao) generatePay(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.finish_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"cast(sum(-bl.update_amount) as UNSIGNED) as amount",
			"cast(sum(case when bl.source_type between 100 and 199 then -bl.update_amount else 0 end) as UNSIGNED) as amount_bet",
			"cast(sum(case when datediff(bl.finish_at, u.created_at) = 0 then -bl.update_amount else 0 end) as UNSIGNED) as amount_new",
			"cast(sum(case when datediff(bl.finish_at, u.created_at) <> 0 then -bl.update_amount else 0 end) as UNSIGNED) as amount_old",
			"count(distinct u.id) as user_cnt",
			"count(distinct (case when datediff(bl.finish_at, u.created_at) = 0 then u.id else null end)) as user_cnt_new",
			"count(distinct (case when datediff(bl.finish_at, u.created_at) <> 0 then u.id else null end)) as user_cnt_old",
		).
		Table("balance_log bl").
		Joins("join users u on bl.user_id = u.id and u.role = 0").
		Where("bl.finish_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("(bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0").
		Group(fmt.Sprintf("date_format(bl.finish_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generatePay: %v", err)
		return nil, err
	}

	return data, nil
}

// 这个是潮币的退款 // 不应包括 3
func (d *PayDao) generateRefund(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.finish_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"cast(sum(bl.update_amount) as UNSIGNED) as refund_amount",
			"count(distinct u.id) as refund_user_cnt",
		).
		Table("balance_log bl").
		Joins("join users u on bl.user_id = u.id and u.role = 0").
		Where("bl.finish_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("bl.source_type in (201, 202, 301) and bl.update_amount > 0").
		Group(fmt.Sprintf("date_format(bl.finish_at, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRefund: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *PayDao) generateRecharge(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(ppo.finish_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(ppo.amount) as recharge_amount",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.amount when 'wechatjs' then ppo.amount else 0 end) as recharge_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.amount else 0 end) as recharge_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then ppo.amount else 0 end) as recharge_amount_huifu",
			"sum(ppo.discount_really) as discount_amount",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.discount_really when 'wechatjs' then ppo.discount_really else 0 end) as discount_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.discount_really else 0 end) as discount_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then ppo.discount_really else 0 end) as discount_amount_huifu",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.role = 0").
		Where("ppo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
		Where("ppo.pay_source_type <> 12"). // 金币储值单独统计
		Group(fmt.Sprintf("date_format(ppo.finish_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRecharge: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *PayDao) generateRechargeRefund(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(ppo.refund_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(refund_amount) as recharge_refund_amount",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.refund_amount when 'wechatjs' then ppo.refund_amount else 0 end) as recharge_refund_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.refund_amount else 0 end) as recharge_refund_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then ppo.refund_amount else 0 end) as recharge_refund_amount_huifu",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.role = 0").
		Where("ppo.refund_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status = 9").
		Where("ppo.pay_source_type <> 12"). // 金币储值 不算在 Refund 中
		Group(fmt.Sprintf("date_format(ppo.refund_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRechargeRefund: %v", err)
		return nil, err
	}

	return data, nil
}

// 储值
func (d *PayDao) generateSaving(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(ppo.finish_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(ppo.amount) as saving_amount",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.amount when 'wechatjs' then ppo.amount else 0 end) as saving_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.amount else 0 end) as saving_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then ppo.amount else 0 end) as saving_amount_huifu",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.role = 0").
		Where("ppo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
		Where("ppo.pay_source_type = 12").
		Group(fmt.Sprintf("date_format(ppo.finish_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateSaving: %v", err)
		return nil, err
	}

	return data, nil
}

// 储值退款
func (d *PayDao) generateSavingRefund(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(rod.refund_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"sum(rod.amount) as saving_refund_amount",
			"sum(case ppo.platform_id when 'wechatapp' then rod.amount when 'wechatjs' then rod.amount else 0 end) as saving_refund_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then rod.amount else 0 end) as saving_refund_amount_ali",
			"sum(case ppo.platform_id when 'huifu' then rod.amount else 0 end) as saving_refund_amount_huifu",
		).
		Table("refund_order_detail rod, users u, pay_payment_order ppo").
		Where("rod.refund_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("rod.status = 3").
		Where("rod.user_id = u.id").
		Where("u.role = 0").
		Where("rod.pay_order_id = ppo.id").
		Group(fmt.Sprintf("date_format(rod.refund_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateSavingRefund: %v", err)
		return nil, err
	}

	return data, nil
}
