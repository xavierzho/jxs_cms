package dao

import (
	"strings"

	"data_backend/internal/app"
	"data_backend/pkg/database"

	"gorm.io/gorm"
)

type GachaRevenue struct {
	GachaID              int64  `gorm:"column:gacha_id; type:bigint;" json:"gacha_id"`
	GachaType            int    `gorm:"column:gacha_type; type:int;" json:"gacha_type"`
	GachaName            string `gorm:"column:gacha_name; type:longtext;" json:"gacha_name"`
	Period               int    `gorm:"column:period; type:int;" json:"period"`
	BoxOutNo             int64  `gorm:"column:box_out_no; type:bigint;" json:"box_out_no"`
	BetNums              int    `gorm:"column:bet_nums; type:int;" json:"bet_nums"`
	TotalNums            int    `gorm:"column:total_nums; type:int;" json:"total_nums"`
	Price                int64  `gorm:"column:price; type:bigint;" json:"price"`
	DiscountPrice        int64  `gorm:"column:discount_price; type:bigint;" json:"discount_price"`
	Amount               int64  `gorm:"column:amount; type:bigint;" json:"amount"`
	InnerPriceBetNormal  int64  `gorm:"column:inner_price_bet_normal; type:bigint;" json:"inner_price_bet_normal"`
	InnerPriceLeftNormal int64  `gorm:"column:inner_price_left_normal; type:bigint;" json:"inner_price_left_normal"`
	InnerPriceBetExtra   int64  `gorm:"column:inner_price_bet_extra; type:bigint;" json:"inner_price_bet_extra"`
	InnerPriceLeftExtra  int64  `gorm:"column:inner_price_left_extra; type:bigint;" json:"inner_price_left_extra"`
}

type RevenueRequestParamsGroup struct {
	IsBoxDim    bool
	GMParams    database.QueryWhereGroup
	AwardParams database.QueryWhereGroup
	OutParams   database.QueryWhereGroup
}

// TODO OutParams 字段 错误
func (d *GachaDao) GetRevenue(paramsGroup RevenueRequestParamsGroup, pager *app.Pager) (summary map[string]any, data []*GachaRevenue, err error) {
	summary = make(map[string]any)

	var betDB, amountDB *gorm.DB
	betDB = d.getRevenueBetDB(paramsGroup)
	if paramsGroup.IsBoxDim {
		amountDB = d.getRevenueAmountBoxDB(paramsGroup)
	} else {
		amountDB = d.getRevenueAmountMachineDB(paramsGroup)
	}

	allJoinStr := " on b.gacha_id = a.gacha_id and b.period = a.period"
	if paramsGroup.IsBoxDim {
		allJoinStr += " and b.box_out_no = a.box_out_no"
	}

	err = d.center.
		Table("(?) as b", betDB).
		Joins("left join (?) as a"+allJoinStr, amountDB).
		Scopes(database.ScopeQuery(paramsGroup.OutParams)).
		Select(
			"count(0) as total",
			"sum(b.bet_nums) as bet_nums",
			"sum(b.total_nums) as total_nums",
			"sum(ifnull(a.amount, 0)) as amount",
			"sum((case a.discount_price when 0 then a.price else a.discount_price end)*(b.total_nums-b.bet_nums)) as amount_left",
			"sum(b.inner_price_bet) as inner_price_bet",
			"sum(b.inner_price_left) as inner_price_left",
			"sum(b.inner_price_bet_extra) as inner_price_bet_extra",
			"sum(b.inner_price_left_extra) as inner_price_left_extra",
		).
		Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("GetRevenue Agg: %v", err)
		return nil, nil, err
	}

	allOrderList := []string{"b.created_at desc"}
	if paramsGroup.IsBoxDim {
		allOrderList = append(allOrderList, "b.box_out_no")
	}
	err = d.center.
		Table("(?) as b", betDB).
		Joins("left join (?) as a"+allJoinStr, amountDB).
		Scopes(database.ScopeQuery(paramsGroup.OutParams)).
		Select("b.*, a.price, a.discount_price, a.amount").
		Order(strings.Join(allOrderList, ",")).
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetRevenue Find: %v", err)
		return nil, nil, err
	}

	return
}

// bet
func (d *GachaDao) getRevenueBetDB(paramsGroup RevenueRequestParamsGroup) *gorm.DB {
	selectList := []string{
		"t.created_at, t.gacha_id", "t.gacha_type", "t.gacha_name", "t.period",
		"sum(t.bet_nums) as bet_nums", "sum(t.total_nums) as total_nums",
		"sum(t.inner_price_bet_normal) as inner_price_bet_normal",
		"sum(t.inner_price_left_normal) as inner_price_left_normal",
		"sum(t.inner_price_bet_extra) as inner_price_bet_extra",
		"sum(t.inner_price_left_extra) as inner_price_left_extra",
		"sum(t.inner_price_bet_normal+t.inner_price_bet_extra) as inner_price_bet",
		"sum(t.inner_price_left_normal+t.inner_price_left_extra) as inner_price_left",
	}
	groupList := []string{"t.created_at, t.gacha_id", "t.gacha_type", "t.gacha_name", "t.period"}
	if paramsGroup.IsBoxDim {
		selectList = append(selectList, "t.box_out_no")
		groupList = append(groupList, "t.box_out_no")
	}

	awardHaving, awardHavingParams := paramsGroup.AwardParams.GetQuerySqlParams()

	return d.center.
		Select(selectList).
		Table("(? union all ?) t", d.getRevenueBetNormalDB(paramsGroup), d.getRevenueBetExtraDB(paramsGroup)).
		Group(strings.Join(groupList, ",")).
		Having(awardHaving, awardHavingParams...)
}

func (d *GachaDao) getRevenueBetNormalDB(paramsGroup RevenueRequestParamsGroup) *gorm.DB {
	betSelectList := []string{
		"gm.created_at, gb.gacha_id", "(gm.`type` + 100) as gacha_type", "gm.name as gacha_name", "gm.period",
		"sum(gba.total_nums-gba.left_nums) as bet_nums", "sum(gba.total_nums) as total_nums",
		"sum((gba.total_nums-gba.left_nums)*i.inner_price) as inner_price_bet_normal",
		"sum(gba.left_nums*i.inner_price) as inner_price_left_normal",
		"0 as inner_price_bet_extra",
		"0 as inner_price_left_extra",
	}
	betGroupList := []string{"gm.created_at, gb.gacha_id", "(gm.`type` + 100)", "gm.name", "gm.period"}
	if paramsGroup.IsBoxDim {
		betSelectList = append(betSelectList, "gb.box_out_no")
		betGroupList = append(betGroupList, "gb.box_out_no")
	}

	return d.center.
		Select(betSelectList).
		Table("gacha_box gb, gacha_machine gm, gacha_box_award gba, item i").
		Scopes(database.ScopeQuery(paramsGroup.GMParams)).
		Where("gb.gacha_id = gm.id").
		Where("gb.gacha_id = gba.gacha_id").
		Where("gb.box_index = gba.box_index").
		Where("gba.item_id = i.id").
		Where("gm.deleted_at is null").
		Where("gb.deleted_at is null").
		Where("gba.deleted_at is null").
		Group(strings.Join(betGroupList, ","))
}

func (d *GachaDao) getRevenueBetExtraDB(paramsGroup RevenueRequestParamsGroup) *gorm.DB {
	extraSelectList := []string{
		"gm.created_at, gb.gacha_id", "(gm.`type` + 100) as gacha_type", "gm.name as gacha_name", "gm.period",
		"0 as bet_nums", "0 as total_nums",
		"0 as inner_price_bet_normal",
		"0 as inner_price_left_normal",
		"sum(case gb.state when 2 then ga.total_nums*i.inner_price else 0 end) as inner_price_bet_extra",
		"sum(case gb.state when 2 then 0 else ga.total_nums*i.inner_price end) as inner_price_left_extra",
	}
	extraGroupList := []string{"gm.created_at, gb.gacha_id", "(gm.`type` + 100)", "gm.name", "gm.period"}
	if paramsGroup.IsBoxDim {
		extraSelectList = append(extraSelectList, "gb.box_out_no")
		extraGroupList = append(extraGroupList, "gb.box_out_no")
	}

	return d.center.
		Select(extraSelectList).
		Table("gacha_box gb, gacha_machine gm, gacha_award ga, item i ").
		Scopes(database.ScopeQuery(paramsGroup.GMParams)).
		Where("gb.gacha_id = gm.id").
		Where("gb.gacha_id = ga.gacha_id").
		Where("ga.level_type <> 1").
		Where("ga.item_id = i.id").
		Where("gm.deleted_at is null").
		Where("gb.deleted_at is null").
		Where("ga.deleted_at is null").
		Group(strings.Join(extraGroupList, ","))
}

// amount
func (d *GachaDao) getRevenueAmountMachineDB(paramsGroup RevenueRequestParamsGroup) *gorm.DB {
	return d.center.
		Select("gm.id as gacha_id, gm.period, gm.price, gm.discount_price, -sum(ifnull(bl.update_amount, 0)) as amount").
		Table("gacha_machine gm").
		Joins("left join balance_log bl on cast(gm.id as char CHARACTER SET utf8mb4) COLLATE utf8mb4_unicode_ci = bl.source_id").
		Scopes(database.ScopeQuery(paramsGroup.GMParams)).
		Group("gm.id, gm.period, gm.price, gm.discount_price")
}

func (d *GachaDao) getRevenueAmountBoxDB(paramsGroup RevenueRequestParamsGroup) *gorm.DB {
	return d.center.
		Select("gm.id as gacha_id, gm.period, gb.box_out_no, gm.price, gm.discount_price, -sum(ifnull(bl.update_amount, 0)) as amount").
		Table("gacha_machine gm, gacha_box gb").
		Joins("left join gacha_user_record gur on gb.gacha_id = gur.gacha_id and gb.box_out_no = gur.box_out_no and gur.count <> 0").
		Joins("left join balance_log bl on gur.request_id = bl.request_id and gur.user_id = bl.user_id and bl.source_type = (gur.gacha_type+100)").
		Where("gb.gacha_id = gm.id ").
		Scopes(database.ScopeQuery(paramsGroup.GMParams)).
		Group("gm.id, gm.period, gb.box_out_no, gm.price, gm.discount_price")
}
