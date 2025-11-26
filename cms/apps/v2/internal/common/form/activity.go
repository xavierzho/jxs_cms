package form

import (
	"fmt"

	"github.com/shopspring/decimal"
)

var COST_AWARD_POINT_STEP = decimal.NewFromInt(10)

type PointType int32

const (
	PointType_Coin           = 0   // 潮币钱包
	PointType_Invite         = 1   // 邀请返佣钱包
	PointType_Gold           = 2   // 金币钱包
	PointType_Point          = 3   // 积分
	PointType_AmountPoint    = 10  // 现金点（活动期间抽赏金额）
	PointType_CostAwardPoint = 11  // 欧气值
	PointType_Free           = 100 // 免费
)

func (t PointType) Valid() (err error) {
	switch t {
	case PointType_Coin, PointType_Invite, PointType_Gold, PointType_Point:
	case PointType_AmountPoint, PointType_CostAwardPoint:
	case PointType_Free:
	default:
		return fmt.Errorf("not expected PointType: %v", t)
	}

	return nil
}

type AwardType int32

const (
	AwardType_Coin           = 0  // 潮币钱包
	AwardType_Invite         = 1  // 邀请返佣钱包
	AwardType_Gold           = 2  // 金币钱包
	AwardType_Point          = 3  // 积分
	AwardType_Coupon         = 10 // 优惠券
	AwardType_Item           = 20 // 物品
	AwardType_CostAwardPoint = 30 // 欧气值
)

func (t AwardType) Valid() (err error) {
	switch t {
	case AwardType_Coin, AwardType_Invite, AwardType_Gold, AwardType_Point:
	case AwardType_Coupon, AwardType_Item, AwardType_CostAwardPoint:
	default:
		return fmt.Errorf("not expected AwardType: %v", t)
	}

	return nil
}
