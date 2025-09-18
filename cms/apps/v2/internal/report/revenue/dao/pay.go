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
	RechargeRefundAmount       uint `gorm:"column:recharge_refund_amount; default:0" json:"recharge_refund_amount"`               // 充值退款金额
	RechargeRefundAmountWeChat uint `gorm:"column:recharge_refund_amount_wechat; default:0" json:"recharge_refund_amount_wechat"` // 充值退款金额(微信)
	RechargeRefundAmountAli    uint `gorm:"column:recharge_refund_amount_ali; default:0" json:"recharge_refund_amount_ali"`       // 充值退款金额(支付宝)
}

func (Pay) TableName() string {
	return "revenue_pay"
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

func (d *PayDao) Generate(cDate time.Time) (dataPay, dataRefund, dataRecharge, dataRechargeRefund *Pay, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataPay, err = d.generatePay(cDate)
		return err
	})

	eg.Go(func() (err error) {
		dataRefund, err = d.generateRefund(cDate)
		return err
	})

	eg.Go(func() (err error) {
		dataRecharge, err = d.generateRecharge(cDate)
		return err
	})

	eg.Go(func() (err error) {
		dataRechargeRefund, err = d.generateRechargeRefund(cDate)
		return err
	})

	err = eg.Wait()

	return
}

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
		Joins("join users u on bl.user_id = u.id and u.is_admin = 0").
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

func (d *PayDao) generateRefund(cDate time.Time) (data *Pay, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"cast(sum(bl.update_amount) as UNSIGNED) as refund_amount",
			"count(distinct u.id) as refund_user_cnt",
		).
		Table("balance_log bl").
		Joins("join users u on bl.user_id = u.id and u.is_admin = 0").
		Where("bl.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("bl.source_type in (3, 201, 202, 301) and bl.update_amount > 0").
		Group(fmt.Sprintf("date_format(bl.created_at, '%s')", pkg.SQL_DATE_FORMAT)).
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
			"sum(case ppo.platform_id when 'wechatapp' then ppo.amount when 'wechatjs' then ppo.amount else 0 end) as  recharge_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.amount else 0 end) as  recharge_amount_ali",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.is_admin = 0").
		Where("ppo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
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
			"sum(case ppo.platform_id when 'wechatapp' then ppo.refund_amount when 'wechatjs' then ppo.refund_amount else 0 end) as  recharge_refund_amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.refund_amount else 0 end) as  recharge_refund_amount_ali",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.is_admin = 0").
		Where("ppo.refund_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.status = 9").
		Group(fmt.Sprintf("date_format(ppo.refund_time, '%s')", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateRechargeRefund: %v", err)
		return nil, err
	}

	return data, nil
}
