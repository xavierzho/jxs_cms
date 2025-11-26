package dao

import (
	"context"
	"fmt"
	"time"

	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type RevenueDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewRevenueDao(engine *gorm.DB, log *logger.Logger) *RevenueDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".RevenueDao")))
	return &RevenueDao{
		engine: engine,
		logger: log,
	}
}

func (d *RevenueDao) All(dateRange [2]time.Time) (data []map[string]interface{}, err error) {
	err = d.engine.Raw(fmt.Sprintf(`
	select
		bal.date,
		bal.wallet as wallet_balance,
		bal.merchant as merchant_balance,
		bal.gold as gold_balance,
		pay.amount as pay_amount,
		pay.amount_bet as pay_amount_bet,
		pay.recharge_amount as recharge_amount,
		pay.recharge_amount_wechat as recharge_amount_wechat,
		pay.recharge_amount_ali as recharge_amount_ali,
		pay.recharge_refund_amount as recharge_refund_amount,
		pay.recharge_refund_amount_wechat as recharge_refund_amount_wechat,
		pay.recharge_refund_amount_ali as recharge_refund_amount_ali,
		pay.saving_amount as saving_amount,
		pay.saving_amount_wechat as saving_amount_wechat,
		pay.saving_amount_ali as saving_amount_ali,
		pay.saving_refund_amount as saving_refund_amount,
		pay.saving_refund_amount_wechat as saving_refund_amount_wechat,
		pay.saving_refund_amount_ali as saving_refund_amount_ali,
		pay.discount_amount as discount_amount,
		pay.discount_amount_wechat as discount_amount_wechat,
		pay.discount_amount_ali as discount_amount_ali,
		draw.amount as draw_amount,
		draw.tax as tax_amount,
		pay.refund_amount as refund_amount,
		a.active_cnt as active_cnt
	from
		revenue_balance bal
		left join revenue_pay pay on bal.date = pay.date
		left join revenue_draw draw on bal.date = draw.date
		left join revenue_active a on bal.date = a.date
	where bal.date between '%s' and '%s'
	order by bal.date desc
		`,
		dateRange[0].Format(pkg.DATE_FORMAT), dateRange[1].Format(pkg.DATE_FORMAT),
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}
