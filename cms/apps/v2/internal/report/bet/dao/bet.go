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

type Bet struct {
	iDao.DailyTypeModel
	UserCnt         uint `gorm:"column:user_cnt; default:0;" json:"user_cnt"`
	BetNums         uint `gorm:"column:bet_nums; default:0;" json:"bet_nums"`
	BoxCntRemaining uint `gorm:"column:box_cnt_remaining; default:0;" json:"box_cnt_remaining"`
	BoxCntNew       uint `gorm:"column:box_cnt_new; default:0;" json:"box_cnt_new"`
	BoxCntClose     uint `gorm:"column:box_cnt_close; default:0;" json:"box_cnt_close"`
	Amount          uint `gorm:"column:amount; default:0;" json:"amount"`
	AmountWeChat    uint `gorm:"column:amount_wechat; default:0;" json:"amount_wechat"`
	AmountAli       uint `gorm:"column:amount_ali; default:0;" json:"amount_ali"`
}

func (Bet) TableName() string {
	return "bet"
}

type BetDao struct {
	*iDao.DailyTypeModelDao[*Bet]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewBetDao(engine, center *gorm.DB, log *logger.Logger) *BetDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".BetDao")))
	return &BetDao{
		DailyTypeModelDao: iDao.NewDailyTypeModelDao[*Bet](engine, log),
		engine:            engine,
		center:            center,
		logger:            log,
	}
}

func (d *BetDao) Generate(cDate time.Time) (dataBet, dataAmount, dataBox, dataPay []*Bet, err error) {
	eg := errgroup.Group{}

	eg.Go(func() (err error) {
		dataBet, err = d.generateBet(cDate)
		return err
	})
	eg.Go(func() (err error) {
		dataAmount, err = d.generateAmount(cDate)
		return err
	})
	eg.Go(func() (err error) {
		dataBox, err = d.generateBox(cDate)
		return err
	})
	eg.Go(func() (err error) {
		dataPay, err = d.generatePay(cDate)
		return err
	})

	err = eg.Wait()

	return
}

func (d *BetDao) generateBet(cDate time.Time) (data []*Bet, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(gur.created_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"gur.gacha_type as data_type",
			"count(distinct gur.user_id) as user_cnt",
			"sum(gur.count) as bet_nums",
		).
		Table("gacha_user_record gur, users u").
		Where("gur.created_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("gur.user_id = u.id").
		Where("u.role = 0").
		Group(fmt.Sprintf("date_format(gur.created_at, '%s'), gur.gacha_type", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateBet: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *BetDao) generateAmount(cDate time.Time) (data []*Bet, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(bl.finish_at, '%s') as date", pkg.SQL_DATE_FORMAT),
			"bl.source_type as data_type",
			"sum(-bl.update_amount) as amount",
		).
		Table("balance_log bl, users u").
		Where("bl.finish_at between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("bl.source_type between 100 and 199 and bl.update_amount <= 0").
		Where("bl.user_id = u.id").
		Where("u.role = 0").
		Group(fmt.Sprintf("date_format(bl.finish_at, '%s'), bl.source_type", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateAmount: %v", err)
		return nil, err
	}

	return data, nil
}

// 剩余箱子 不算 下架 部分（服务端 箱子 下架部分代码需要修改）
func (d *BetDao) generateBox(cDate time.Time) (data []*Bet, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("'%s' as date", cDate.Format(pkg.DATE_FORMAT)),
			"gm.type as data_type",
			"count( case when gb.state = 1 and gm.flag & 1 = 1 then 1 else null end) as box_cnt_remaining",
			fmt.Sprintf("count( case when datediff(gb.created_at, '%s') = 0 then 1 else null end) as box_cnt_new", cDate.Format(pkg.DATE_FORMAT)),
			fmt.Sprintf("count( case when datediff(gb.updated_at, '%s') = 0 and gb.state = 2 then 1 else null end) as box_cnt_close", cDate.Format(pkg.DATE_FORMAT)),
		).
		Table("gacha_box gb, gacha_machine gm").
		Where("gb.deleted_at is null").
		Where("gb.gacha_id = gm.id").
		Where("gm.deleted_at is null").
		Group("gm.type").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generateBox: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *BetDao) generatePay(cDate time.Time) (data []*Bet, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(ppo.finish_time, '%s') as date", pkg.SQL_DATE_FORMAT),
			"gm.type as data_type",
			"sum(case ppo.platform_id when 'wechatapp' then ppo.amount-ppo.refund_amount when 'wechatjs' then ppo.amount-ppo.refund_amount else 0 end) as  amount_wechat",
			"sum(case ppo.platform_id when 'alipay' then ppo.amount-ppo.refund_amount else 0 end) as  amount_ali",
		).
		Table("pay_payment_order ppo, users u, gacha_machine gm").
		Where("ppo.finish_time between ? and ?", cDate.Format(pkg.DATE_TIME_MIL_FORMAT), cDate.Add(24*time.Hour-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ppo.pay_source_type = 100").
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
		Where("ppo.user_id = u.id").
		Where("u.role = 0").
		Where("ppo.pay_source_id = gm.id").
		Group(fmt.Sprintf("date_format(ppo.finish_time, '%s'), gm.type", pkg.SQL_DATE_FORMAT)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("generatePay: %v", err)
		return nil, err
	}

	return data, nil
}
