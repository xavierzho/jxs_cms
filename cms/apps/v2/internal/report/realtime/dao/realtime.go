package dao

import (
	"context"
	"fmt"
	"time"

	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type RealtimeDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewRealtimeDao(center *gorm.DB, log *logger.Logger) *RealtimeDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".RealtimeDao")))
	return &RealtimeDao{
		center: center,
		logger: log,
	}
}

// 活跃用户数
func (d *RealtimeDao) GetActiveUserCnt(startTime, endTime time.Time) (count int64, err error) {
	err = d.center.
		Select("count(distinct l.user_id) as user_cnt").
		Table("logon_logs l, users u").
		Where("l.user_id = u.id").
		Where("l.created_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Scan(&count).Error
	if err != nil {
		d.logger.Errorf("GetActiveUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}

// 参与用户数 // TODO 考虑 将 结果 改为 list 兼容 后续新增 类型，仅在 service 中处理 类型
func (d *RealtimeDao) GetPatingUserCnt(startTime, endTime time.Time) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	err = d.center.Raw(fmt.Sprintf(`
select
	count(distinct t.user_id) as user_cnt,
	count(distinct case t.biz_type when 101 then t.user_id else null end) as user_cnt_101,
	count(distinct case t.biz_type when 102 then t.user_id else null end) as user_cnt_102,
	count(distinct case t.biz_type when 103 then t.user_id else null end) as user_cnt_103,
	count(distinct case t.biz_type when 104 then t.user_id else null end) as user_cnt_104,
	count(distinct case t.biz_type when 200 then t.user_id else null end) as user_cnt_200,
	count(distinct case t.biz_type when 300 then t.user_id else null end) as user_cnt_300,
	count(distinct case t.biz_type when 601 then t.user_id else null end) as user_cnt_600
from
	(
	select distinct 200 as biz_type, tv.user_id -- 创建集市订单的用户
	from market_order tv, users u
	where tv.created_at between '%[1]s' and '%s' and tv.user_id = u.id and u.is_admin = 0
	union
	select distinct 200 as biz_type, muo.user_id
	from market_user_offer muo, users u
	where muo.created_at between '%[1]s' and '%s' and muo.user_id = u.id and u.is_admin = 0
	union
	SELECT distinct 300 as biz_type, o.user_id
	FROM %s o, users u -- 发货用户
	Where o.pay_time between %[4]d and %d and o.user_id = u.id and u.is_admin = 0
	union
	select distinct bl.source_type as biz_type, bl.user_id
	from balance_log bl, users u
	where bl.created_at between '%[1]s' and '%s'
		and (bl.source_type between 100 and 199 or bl.source_type in (601)) -- 抽赏 + 商城
		and bl.update_amount <= 0
		and bl.user_id = u.id and u.is_admin = 0
	) t	
	`,
		startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT),
		"`order`",
		startTime.UnixMilli(), endTime.UnixMilli(),
	)).
		Scan(&data).Error
	if err != nil {
		d.logger.Errorf("GetPatingUserCnt: %v", err)
		return nil, err
	}

	return data, nil
}

// 付费用户数
func (d *RealtimeDao) GetPayData(startTime, endTime time.Time) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	err = d.center.
		Select(
			"cast(sum(-bl.update_amount) as UNSIGNED) as amount",
			"count(distinct bl.user_id) as user_cnt",
		).
		Table("balance_log bl, users u").
		Where("bl.user_id = u.id").
		Where("u.is_admin = 0").
		Where("bl.finish_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("(bl.source_type between 100 and 199 or bl.source_type in (201,202,300,301,302,303,304,601)) and bl.update_amount <= 0").
		Scan(&data).Error
	if err != nil {
		d.logger.Errorf("GetPayData: %v", err)
		return nil, err
	}

	return data, nil
}

// 充值
func (d *RealtimeDao) GetRechargeData(startTime, endTime time.Time) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	err = d.center.
		Select(
			"sum(ppo.amount) as amount",
			"count(distinct ppo.user_id) as user_cnt",
		).
		Table("pay_payment_order ppo").
		Joins("join users u on ppo.user_id = u.id and u.is_admin = 0").
		Where("ppo.finish_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		// Where("ppo.pay_source_type IN (100,201,202)"). -- 所有充值都计算
		Where("ppo.status in (4,7,8,9,10,11,12,13,14)").
		Scan(&data).Error
	if err != nil {
		d.logger.Errorf("GetRechargeData: %v", err)
		return nil, err
	}

	return
}

// 退款(￥)用户数
func (d *RealtimeDao) GetDrawData(startTime, endTime time.Time) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	err = d.center.
		Select(
			"cast(sum(pdo.amount) as UNSIGNED) as amount",
			"count(distinct pdo.user_id) as user_cnt",
		).
		Table("pay_payout_order pdo, users u").
		Where("pdo.finish_time between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("pdo.state in (6, 12)").
		Where("pdo.user_id = u.id").
		Where("u.is_admin = 0").
		Scan(&data).Error
	if err != nil {
		d.logger.Errorf("GetDrawData: %v", err)
		return nil, err
	}

	return data, nil
}

// 注册用户数
func (d *RealtimeDao) GetNewUserCnt(startTime, endTime time.Time) (count int64, err error) {
	err = d.center.
		Select("count(distinct u.id) as user_cnt").
		Table("users u").
		Where("u.created_at between ? and ?", startTime.Format(pkg.DATE_TIME_MIL_FORMAT), endTime.Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("u.is_admin = 0").
		Scan(&count).Error
	if err != nil {
		d.logger.Errorf("GetNewUserCnt: %v", err)
		return 0, err
	}

	return count, nil
}
