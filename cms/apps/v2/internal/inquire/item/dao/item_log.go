package dao

import (
	"fmt"
	"strings"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"

	"gorm.io/gorm"
)

type ItemLog struct {
	ID             string `gorm:"column:id; type:varchar(255)" json:"id"`
	DateTime       string `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID         int64  `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName       string `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	LogType        int    `gorm:"column:log_type; type:int" json:"log_type"` // 101,102,103,200
	LogTypeName    string `gorm:"column:log_type_name; type:longtext;" json:"log_type_name"`
	Period         int    `gorm:"column:period; type:int;" json:"period"`
	BetNums        int    `gorm:"column:bet_nums; type:int" json:"bet_nums"`
	LevelType      int    `gorm:"column:level_type; type:int" json:"level_type"` // 0,1,2,3,4
	UpdateAmount   int64  `gorm:"column:update_amount; type:bigint" json:"update_amount"`
	ShowPrice      int64  `gorm:"column:show_price; type:bigint" json:"show_price"`
	InnerPrice     int64  `gorm:"column:inner_price; type:bigint" json:"inner_price"`
	RecyclingPrice int64  `gorm:"column:recycling_price; type:bigint" json:"recycling_price"`
}

// list
// gacha_user_record union all (market_order join market_offer) union all (market_order join market_offer)
// id(gacha_user_record.id / market_offer.id), dateTime, user_id, user_name, log_type, level_type, update_amount, show_price, inner_price,
// level_type 一番赏 last (最后一个); 洞洞乐 扭蛋机 lucky （随机抽一个）
// TODO 优化 all 和 list 分开
func (d *ItemDao) GetLog(dateTimeRange [2]time.Time, logTypeList []int, paramsGroup AllRequestParamsGroup, pager *app.Pager) (summary map[string]any, data []*ItemLog, err error) {
	var maxNums int
	if pager != nil {
		maxNums = pager.Page * pager.PageSize
	}

	summary = make(map[string]any)

	if !(paramsGroup.BetFlag || paramsGroup.MarketFlag || paramsGroup.OrderFlag || paramsGroup.ActivityFlag.Flag || paramsGroup.AdminFlag) {
		err = fmt.Errorf("invalid listType: %v", logTypeList)
		d.logger.Errorf("GetLog: %v", err)
		return nil, nil, err
	}

	var dbList []any
	if paramsGroup.BetFlag {
		dbList = append(dbList, d.getLogBetDB(dateTimeRange, paramsGroup))
	}
	if paramsGroup.MarketFlag {
		dbList = append(dbList, d.getLogMarketCreatorDB(dateTimeRange, paramsGroup))
		dbList = append(dbList, d.getLogMarketOffererDB(dateTimeRange, paramsGroup))
	}
	if paramsGroup.OrderFlag {
		dbList = append(dbList, d.getLogOrderDB(dateTimeRange, paramsGroup))
	}
	if paramsGroup.ActivityFlag.Flag {
		if paramsGroup.ActivityFlag.CostAward {
			dbList = append(dbList, d.getLogActivityCostAwardDB(dateTimeRange, paramsGroup))
		}
		if paramsGroup.ActivityFlag.CostRank {
			dbList = append(dbList, d.getLogActivityCostRankDB(dateTimeRange, paramsGroup))
		}
		if paramsGroup.ActivityFlag.ItemExchange {
			dbList = append(dbList, d.getLogActivityItemExchangeDB(dateTimeRange, paramsGroup))
		}
		if paramsGroup.ActivityFlag.PrizeWheel {
			dbList = append(dbList, d.getLogActivityPrizeWheelDB(dateTimeRange, paramsGroup))
		}
	}
	if paramsGroup.AdminFlag {
		dbList = append(dbList, d.getLogAdminDB(dateTimeRange, paramsGroup))
	}

	sqlList := make([]string, len(dbList))
	sqlListPager := make([]string, len(dbList))
	var orderStr = "date_time desc, log_type, level_type, user_id"
	for ind := range sqlList {
		sqlList[ind] = "?"
		sqlListPager[ind] = fmt.Sprintf("(? order by %s limit %d)", orderStr, maxNums)
	}

	summaryDB := d.center.Table("("+strings.Join(sqlList, " union all ")+") as t", dbList...)
	err = summaryDB.
		Select(
			"count(0) as total",
			"count(distinct user_id) as user_cnt",
			"sum(bet_nums) as bet_nums",
			"sum(update_amount) as update_amount",
			"sum(show_price) as show_price",
			"sum(inner_price) as inner_price",
			"sum(recycling_price) as recycling_price",
		).
		Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("GetLog Agg: %v", err)
		return nil, nil, err
	}

	// 针对分页优化 // 能有那么一点效果吧... // TODO 可以考虑 在翻页的情况下 增加 当前页最小/大时间作为 条件 // 用先查 id 再查 数据的方式 改起来很麻烦, sql 很丑...
	allDB := summaryDB
	if pager != nil {
		allDB = d.center.Table("("+strings.Join(sqlListPager, " union all ")+") as t", dbList...)
	}

	err = allDB.
		Select("*").
		Order(orderStr).
		Scopes(func(d *gorm.DB) *gorm.DB {
			if pager != nil {
				return database.Paginate(pager.Page, pager.PageSize)(d)
			}
			return d
		}).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetLog Find: %v", err)
		return nil, nil, err
	}

	return
}

// bet
func (d *ItemDao) getLogBetDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.ItemParams...)
	queryParams = append(queryParams, paramsGroup.AmountParams...)
	having, havingParams := queryParams.GetQuerySqlParams()

	return d.center.
		Select(
			"t.id",
			"t.date_time",
			"t.user_id",
			"t.user_name",
			"t.log_type",
			"t.log_type_name",
			"t.period",
			"(case j.level_type when 1 then t.bet_nums else 0 end) as bet_nums",
			"j.level_type",
			"(case j.level_type when 1 then t.update_amount else 0 end) as update_amount",
			"sum(i.show_price * j.nums) as show_price",
			"sum(i.inner_price * j.nums) as inner_price",
			"sum(i.recycling_price * j.nums) as recycling_price",
		).
		Table("(?) as t, ? as j, item i", d.getLogBetGachaDB(dateTimeRange, paramsGroup), d.getLogBetItemJsonDB()).
		Where("j.item_id = i.id").
		Group("t.id,t.date_time,t.user_id,t.user_name,t.log_type,t.log_type_name,t.period,t.bet_nums,j.level_type,(case j.level_type when 1 then t.update_amount else 0 end)").
		Having(having, havingParams...)
}

func (d *ItemDao) getLogBetGachaDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			"gur.id",
			fmt.Sprintf("DATE_FORMAT(gur.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"gur.user_id",
			"u.nickname as user_name",
			"(gur.gacha_type+100) as log_type",
			"gur.gacha_name as log_type_name",
			"gur.period",
			"gur.count as bet_nums",
			"bl.update_amount",
			"gur.items",
		).
		Table("users u, gacha_user_record gur").
		Joins("left join balance_log bl on gur.request_id = bl.request_id and gur.user_id = bl.user_id and bl.source_type = (gur.gacha_type+100)").
		Where("gur.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("gur.user_id = u.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.GachaParams))
}

func (d *ItemDao) getLogBetItemJsonDB() *gorm.DB {
	return d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
		JSON_UNQUOTE(t.items), 
		'$[*]' COLUMNS(
				nums int path '$.Nums',
				item_id bigint path '$.ItemID',
				level_type int path '$.LevelType'
			)
	)
	`, "\t", " "), "\n", " "))
}

// market
func (d *ItemDao) getLogMarketCreatorDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	having, havingParams := paramsGroup.AmountParams.GetQuerySqlParams()

	return d.center.
		Select(
			"muo.id",
			fmt.Sprintf("DATE_FORMAT(muo.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"200 as log_type",
			"'创建者' as log_type_name",
			"0 as period",
			"0 as bet_nums",
			"0 as level_type",
			"(muo.offer_amount-mo.ask_price) as update_amount",
			"0 as show_price",
			"0 as inner_price",
			"0 as recycling_price",
		).
		Table("market_order mo, market_user_offer muo, users u").
		Where("muo.state = 2").
		Where("muo.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("muo.order_id = mo.id").
		Where("mo.user_id = u.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Having(having, havingParams...)
}

func (d *ItemDao) getLogMarketOffererDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	having, havingParams := paramsGroup.AmountParams.GetQuerySqlParams()

	return d.center.
		Select(
			"muo.id",
			fmt.Sprintf("DATE_FORMAT(muo.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"200 as log_type",
			"'交易者' as log_type_name",
			"0 as period",
			"0 as bet_nums",
			"0 as level_type",
			"(mo.ask_price-muo.offer_amount) as update_amount",
			"0 as show_price",
			"0 as inner_price",
			"0 as recycling_price",
		).
		Table("market_order mo, market_user_offer muo, users u").
		Where("muo.state = 2").
		Where("muo.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("muo.order_id = mo.id").
		Where("muo.user_id = u.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Having(having, havingParams...)
}

// order
func (d *ItemDao) getLogOrderDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.ItemParams...)
	queryParams = append(queryParams, paramsGroup.AmountParams...)
	having, havingParams := queryParams.GetQuerySqlParams()
	return d.center.
		Select(
			"t.id",
			"t.date_time",
			"t.user_id",
			"t.user_name",
			"t.log_type",
			"t.log_type_name",
			"t.period",
			"count(j.item_id) as bet_nums",
			"t.level_type",
			"t.update_amount",
			"sum(i.show_price) as show_price",
			"sum(i.inner_price) as inner_price",
			"sum(i.recycling_price) as recycling_price",
		).
		Table("(?) as t, ? as j, item i", d.getLogOrderLogDB(dateTimeRange, paramsGroup), d.getLogOrderItemJsonDB()).
		Where("j.item_id = i.id").
		Group("t.id,t.date_time,t.user_id,t.user_name,t.log_type,t.log_type_name,t.period,t.bet_nums,t.level_type,t.update_amount").
		Having(having, havingParams...)
}

func (d *ItemDao) getLogOrderLogDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			"o.id as id",
			fmt.Sprintf("DATE_FORMAT(FROM_UNIXTIME(CAST((o.delivery_time / 1000) AS SIGNED)), '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"300 as log_type", // 发货
			"'' as log_type_name",
			"0 as period",
			"0 as bet_nums",
			"0 as level_type",
			"o.pay_amount as update_amount",
			"o.order_items as items",
		).
		Table("`order` o, users u").
		Where("o.delivery_time between ? and ?", dateTimeRange[0].UnixMilli(), dateTimeRange[1].Add(time.Second-time.Millisecond).UnixMilli()).
		Where("o.state in (4,5)"). // 已完成
		Where("o.user_id = u.id").
		Group(fmt.Sprintf("DATE_FORMAT(FROM_UNIXTIME(CAST((o.delivery_time / 1000) AS SIGNED)), '%s'), o.delivery_time, u.id, u.nickname, o.id", pkg.SQL_DATE_TIME_FORMAT)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams))
}

func (d *ItemDao) getLogOrderItemJsonDB() *gorm.DB {
	return d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
		JSON_UNQUOTE(t.items), 
		'$[*]' COLUMNS(
				name varchar(255) path '$.name',
				state int path '$.state',
				item_id bigint path '$.item_id',
				stock_id int path '$.stock_id'
			)
	)
	`, "\t", " "), "\n", " "))
}

// activity
func (d *ItemDao) getLogActivityCostAwardDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.ItemParams...)
	queryParams = append(queryParams, paramsGroup.AmountParams...)
	having, havingParams := queryParams.GetQuerySqlParams()

	return d.center.
		Select(
			"ua.params as id",
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100002 as log_type", // 欧气值
			"(case when ua.params_3->'$.type'='0' then '兑换' when ua.params_3->'$.type'='1' then '购买' else '' end) as log_type_name",
			"0 as period",
			"cast(ua.params_2 as signed) as bet_nums",
			"0 as level_type",
			"max((case when ua.params_3->'$.type'='0' then 0 else cast(ua.params_2 as signed) end) *cac.price) as update_amount",
			"sum(i.show_price * cac.award_num * cast(ua.params_2 as signed)) as show_price",
			"sum(i.inner_price * cac.award_num * cast(ua.params_2 as signed)) as inner_price",
			"sum(i.recycling_price * cac.award_num * cast(ua.params_2 as signed)) as recycling_price",
		).
		Table("activity a, user_activity ua, users u, activity_cost_award_config cac, item i").
		Where("a.name = '欧气值'").
		Where("a.id = ua.activity_id").
		Where("ua.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = cac.config_id").
		Where("cac.award_type = 20").
		Where("cac.award_value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Group(fmt.Sprintf("ua.params, DATE_FORMAT(ua.created_at, '%s'), u.id, u.nickname, (case when ua.params_3->'$.type'='0' then '兑换' when ua.params_3->'$.type'='1' then '购买' else '' end), cast(ua.params_2 as signed)", pkg.SQL_DATE_TIME_FORMAT)).
		Having(having, havingParams...)
}

func (d *ItemDao) getLogActivityCostRankDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	//! 窗口函数不能用于 having // You cannot use the alias 'near_no' of an expression containing a window function in this context.
	db := d.center.
		Select(
			"CONCAT_WS('|', ua.params, crac.no_limit) as id", // period | No | no_limit
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"crac.period",
			"sum(i.show_price * crac.num) as show_price",
			"sum(i.inner_price * crac.num) as inner_price",
			"sum(i.recycling_price * crac.num) as recycling_price",
			fmt.Sprintf("ROW_NUMBER() OVER(PARTITION BY ua.params, DATE_FORMAT(ua.created_at, '%s'), u.id, u.nickname, crac.period ORDER BY no_limit) as near_no", pkg.SQL_DATE_TIME_FORMAT),
		).
		Table("activity a, user_activity ua, users u, activity_cost_rank_award_config crac, item i").
		Where("a.name = '消费排行'").
		Where("a.id = ua.activity_id").
		Where("ua.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("ua.user_id = u.id").
		Where("cast(substring_index(ua.params, '|', 1) as SIGNED) = crac.period").
		Where("cast(substring_index(ua.params, '|', -1) as SIGNED) <= crac.no_limit").
		Where("crac.type = 20").
		Where("crac.value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Group(fmt.Sprintf("ua.params, DATE_FORMAT(ua.created_at, '%s'), u.id, u.nickname, crac.period, crac.no_limit", pkg.SQL_DATE_TIME_FORMAT))

	return d.center.
		Select(
			"id",
			"date_time",
			"user_id",
			"user_name",
			"100003 as log_type", // 欧气排名
			"'' as log_type_name",
			"t.period as period",
			"0 as bet_nums",
			"0 as level_type",
			"0 as update_amount",
			"show_price",
			"inner_price",
			"recycling_price",
		).
		Table("(?) t", db).
		Where("near_no = 1").
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) getLogActivityItemExchangeDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			"ua.id as id",
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100004 as log_type", // 物品置换
			"'' as log_type_name",
			"0 as period",
			"ua.params_2 as bet_nums",
			"0 as level_type",
			"0 as update_amount",
			"0 as show_price",
			"0 as inner_price",
			"0 as recycling_price",
		).
		Table("activity as a, user_activity ua, users u").
		Where("a.name = '物品置换'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("ua.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Group(fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s'), ua.id, u.id, u.nickname, ua.params_2", pkg.SQL_DATE_TIME_FORMAT)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams))
}

// admin // TODO 新版本上线后改为 通过 request_id 做id
func (d *ItemDao) getLogAdminDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	having, havingParams := paramsGroup.ItemParams.GetQuerySqlParams()

	return d.center.
		Select(
			"CONCAT_WS('|', cs.create_time, u.id) as id",
			fmt.Sprintf("DATE_FORMAT(cs.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"999999 as log_type", // 管理员发放
			"'' as log_type_name",
			"0 as period",
			"count(0) as bet_nums",
			"0 as level_type",
			"0 as update_amount",
			"sum(i.show_price) as show_price",
			"sum(i.inner_price) as inner_price",
			"sum(i.recycling_price) as recycling_price",
		).
		Table("cabinet_stock cs, users u, item i").
		Where("cs.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Where("cs.theme_type = 999"). // 管理员发放
		Where("cs.user_id = u.id").
		Where("cs.item_id = i.id").
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Group(fmt.Sprintf("DATE_FORMAT(cs.created_at, '%s'), cs.create_time, u.id, u.nickname", pkg.SQL_DATE_TIME_FORMAT)).
		Having(having, havingParams...)
}

// 转盘抽奖
func (d *ItemDao) getLogActivityPrizeWheelDB(dateTimeRange [2]time.Time, paramsGroup AllRequestParamsGroup) *gorm.DB {
	fmt.Println("======================", paramsGroup)
	return d.center.
		Select(
			"apwh.id as id",
			fmt.Sprintf("DATE_FORMAT(apwh.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100005 as log_type", // 转盘抽奖
			"apwc.name as log_type_name",
			"apwc.period as period",
			"1 as bet_nums",
			"0 as level_type",
			"0 as update_amount",
			"i.show_price as show_price",
			"i.inner_price as inner_price",
			"i.recycling_price as recycling_price",
		).
		Table("activity as a, user_activity ua, users u,activity_prize_wheel_history apwh,activity_prize_wheel_award_config apwac,item i,activity_prize_wheel_config apwc").
		Where("a.name = '抽奖转盘'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("apwh.user_id = u.id").
		Where("apwh.award_id = apwac.id").
		Where("apwh.config_id = apwc.id").
		Where("apwac.value = i.id").
		Where("apwh.created_at between ? and ?", dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second-time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)).
		Group(fmt.Sprintf("DATE_FORMAT(apwh.created_at, '%s'), apwh.id, u.id, u.nickname,apwc.name,apwc.period, i.show_price,i.inner_price,i.recycling_price", pkg.SQL_DATE_TIME_FORMAT)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}
