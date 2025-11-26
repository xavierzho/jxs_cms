package dao

import (
	"fmt"
	"strings"

	cDao "data_backend/apps/v2/internal/common/dao"
	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"

	"gorm.io/gorm"
)

// base; activity
type Item struct {
	cDao.Item
	Nums int `gorm:"column:nums; type:int" json:"nums"`
	No   int `gorm:"column:no; type:int" json:"no"`
}

type BetItemDetail struct {
	// GachaName string `gorm:"column:gacha_name; type:longtext" json:"gacha_name"`
	BoxOutNo int64 `gorm:"column:box_out_no; type:bigint" json:"box_out_no"`
	// BetNums   int    `gorm:"column:bet_nums; type:int" json:"bet_nums"`
	Item
}

type MarketItemDetail struct {
	UserID   int64  `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName string `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	Amount   int64  `gorm:"column:amount; type:bigint" json:"amount"`
	Item
}

// export
type ItemDetail struct {
	DateTime    string `gorm:"column:date_time; type:varchar(19)" json:"date_time"`
	UserID      int64  `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName    string `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	LogType     int    `gorm:"column:log_type; type:int" json:"log_type"` // 101,102,103,200
	LogTypeName string `gorm:"column:log_type_name; type:longtext;" json:"log_type_name"`
	Period      int    `gorm:"column:period; type:int;" json:"period"`
	BoxOutNo    int64  `gorm:"column:box_out_no; type:bigint" json:"box_out_no"`
	Item
}

// ! 仅1w条 // 数量过多会占用 大量内存
func (d *ItemDao) GetDetail(logTypeList []int, paramsGroup AllRequestParamsGroup, pager *app.Pager) (data []*ItemDetail, total int64, err error) {
	if !(paramsGroup.BetFlag || paramsGroup.MarketFlag || paramsGroup.ActivityFlag.Flag || paramsGroup.AdminFlag || paramsGroup.OrderFlag || paramsGroup.TaskFlag || paramsGroup.PublicizeFlag) {
		err = fmt.Errorf("invalid listType: %v", logTypeList)
		d.logger.Errorf("GetDetail: %v", err)
		return nil, 0, err
	}

	var dbList []any
	if paramsGroup.BetFlag {
		dbList = append(dbList, d.getDetailBetDB(paramsGroup))
	}
	if paramsGroup.MarketFlag {
		dbList = append(dbList, d.getDetailMarketCreatorDB(paramsGroup))
		dbList = append(dbList, d.getDetailMarketOffererDB(paramsGroup))
	}
	if paramsGroup.OrderFlag {
		dbList = append(dbList, d.getDetailOrderDB(paramsGroup))
	}
	if paramsGroup.ActivityFlag.Flag {
		if paramsGroup.ActivityFlag.CostAward {
			dbList = append(dbList, d.getDetailActivityCostAwardDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.CostRank {
			dbList = append(dbList, d.getDetailActivityCostRankDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.ItemExchange {
			dbList = append(dbList, d.getDetailActivityItemExchangeOutDB(paramsGroup))
			dbList = append(dbList, d.getDetailActivityItemExchangeIntoDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.PrizeWheel {
			dbList = append(dbList, d.getDetailActivityPrizeWheelDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.StepByStep {
			dbList = append(dbList, d.getDetailActivityStepByStepDB(paramsGroup))
			dbList = append(dbList, d.getDetailActivityStepByStepRankDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.SignIn {
			dbList = append(dbList, d.getDetailActivitySignInDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.LuckyNum {
			dbList = append(dbList, d.getDetailActivityLuckyNumDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.RedemptionCode {
			dbList = append(dbList, d.getDetailActivityRedemptionCodeDB(paramsGroup))
		}
		if paramsGroup.ActivityFlag.Lottery {
			dbList = append(dbList, d.getDetailActivityLotteryDB(paramsGroup))
		}
	}
	if paramsGroup.TaskFlag {
		dbList = append(dbList, d.getDetailTaskDB(paramsGroup))
	}
	if paramsGroup.AdminFlag {
		dbList = append(dbList, d.getDetailAdminDB(paramsGroup))
	}
	if paramsGroup.PublicizeFlag {
		dbList = append(dbList, d.getDetailPublicizeDB(paramsGroup))
	}

	sqlList := make([]string, len(dbList))
	for ind := range sqlList {
		sqlList[ind] = "?"
	}

	// ! 正式服：十分奇怪 count(user_name) 的情况下 sql 执行方式与其他时候不同，会使用 item 索引，其他情况没有使用这个索引 导致 很慢（通过 EXPLAIN ANALYZE 分析） // 测试服都会走索引
	err = d.center.Table("("+strings.Join(sqlList, " union all ")+") as t", dbList...).Select("count(user_name)").Scan(&total).Error
	if err != nil {
		d.logger.Errorf("GetDetail Count: %v", err)
		return nil, 0, err
	}

	allDB := d.center.Table("("+strings.Join(sqlList, " union all ")+") as t", dbList...)

	// ! 一个奇怪的 bug 使用 select * 部分情况（sqlList仅有一个元素时）会导致 `Unknown column 't.date_time' in 'order clause'` 但 sql 直接执行是正确的
	err = allDB.
		Select(
			"date_time",
			"user_id",
			"user_name",
			"log_type",
			"log_type_name",
			"period",
			"box_out_no",
			"amount",
			"item_id",
			"item_name",
			"level_name",
			"cover_thumb",
			"show_price",
			"inner_price",
			"recycling_price",
			"nums",
			"no",
		).
		// Count(&total). // ? 这个地方使用 count 结果 count(*) 语句不走item索引
		Order("t.date_time desc, t.log_type, t.user_id, t.level_name, t.inner_price desc, t.nums desc").
		Scopes(func(d *gorm.DB) *gorm.DB {
			if pager != nil {
				return database.Paginate(pager.Page, pager.PageSize)(d)
			}
			return d
		}).
		Limit(10000).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetail Find: %v", err)
		return nil, 0, err
	}

	return
}

// 抽赏
func (d *ItemDao) GetDetailBet(queryParams database.QueryWhereGroup) (data []*BetItemDetail, err error) {
	err = d.getDetailBetDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("j.level_name, i.inner_price desc, j.nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailBet: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailBetDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	jsonDB := d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
		JSON_UNQUOTE(t.items), 
		'$[*]' COLUMNS(
				no int path '$.No',
				nums int path '$.Nums',
				item_id bigint path '$.ItemID',
				level_type int path '$.LevelType',
				level_name longtext path '$.LevelName'
			)
	)
	`, "\t", " "), "\n", " "))

	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"(t.gacha_type+100) as log_type",
			"t.gacha_name as log_type_name",
			"t.period",
			"t.box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"j.level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"j.nums",
			"j.no",
		).
		Table("users u, gacha_user_record t, ? as j, item i", jsonDB).
		Where("u.id = t.user_id").
		Where("j.item_id = i.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.GachaParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams))
}

// 集市
func (d *ItemDao) GetDetailMarket(queryParams database.QueryWhereGroup) (dataCreator, dataOfferer []*MarketItemDetail, err error) {
	dataCreator, err = d.getDetailMarketCreator(queryParams)
	if err != nil {
		return
	}

	dataOfferer, err = d.getDetailMarketOfferer(queryParams)
	if err != nil {
		return
	}

	return
}

func (d *ItemDao) getDetailMarketCreator(queryParams database.QueryWhereGroup) (data []*MarketItemDetail, err error) {
	err = d.getDetailMarketCreatorDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("j.inner_price desc, j.nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("getDetailMarketCreator: %v", err)
		return
	}

	return
}

// detail 中的金额 为加价金额 而不是 成交后的变动金额
func (d *ItemDao) getDetailMarketCreatorDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.DateTimeParams...)
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.ItemParams...)
	queryParams = append(queryParams, paramsGroup.OtherParams...)
	whereParams, sqlParam := queryParams.GetQuerySqlParams()
	if whereParams != "" {
		whereParams = " and " + whereParams
	}

	jsonDB := d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf(`
	select
		t.id, i.id as item_id, i.name as item_name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price, count(i.id) as nums
	from 
		users u,
		market_order mo,
		market_user_offer t,
		JSON_TABLE(
			JSON_UNQUOTE(mo.order_items), 
			'$[*]' COLUMNS(
					item_id bigint path '$.ItemID'
				)
		) as j,
		item i
	where
		t.state = 2
		and t.order_id = mo.id
		and mo.user_id = u.id
		and j.item_id = i.id %s
	group by t.id, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price
	`, whereParams,
	), "\t", " "), "\n", " "), sqlParam...)

	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"200 as log_type",
			"'创建者' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"mo.ask_price as amount",
			"j.item_id",
			"j.item_name",
			"'' as level_name",
			"j.cover_thumb",
			"j.show_price",
			"j.inner_price",
			"j.recycling_price",
			"j.nums",
			"0 as no",
		).
		Table("users u, market_order mo, market_user_offer t").
		Joins("left join (?) as j on t.id = j.id", jsonDB).
		Where("t.state = 2").
		Where("t.order_id = mo.id").
		Where("mo.user_id = u.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams))
}

func (d *ItemDao) getDetailMarketOfferer(queryParams database.QueryWhereGroup) (data []*MarketItemDetail, err error) {
	err = d.getDetailMarketOffererDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("j.inner_price desc, j.nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("getDetailMarketOfferer: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailMarketOffererDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	var queryParams database.QueryWhereGroup
	queryParams = append(queryParams, paramsGroup.DateTimeParams...)
	queryParams = append(queryParams, paramsGroup.UsersParams...)
	queryParams = append(queryParams, paramsGroup.ItemParams...)
	queryParams = append(queryParams, paramsGroup.OtherParams...)
	whereParams, sqlParam := queryParams.GetQuerySqlParams()
	if whereParams != "" {
		whereParams = " and " + whereParams
	}

	jsonDB := d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(fmt.Sprintf(`
	select
		t.id, i.id as item_id, i.name as item_name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price, count(i.id) as nums
	from 
		users u, 
		market_user_offer t,
		JSON_TABLE(
			JSON_UNQUOTE(t.offer_items), 
			'$[*]' COLUMNS(
					item_id bigint path '$.ItemID'
				)
		) as j,
		item i
	where
		t.state = 2
		and t.user_id = u.id
		and j.item_id = i.id %s
	group by t.id, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price
	`, whereParams,
	), "\t", " "), "\n", " "), sqlParam...)

	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"200 as log_type",
			"'交易者' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"t.offer_amount as amount",
			"j.item_id",
			"j.item_name",
			"'' as level_name",
			"j.cover_thumb",
			"j.show_price",
			"j.inner_price",
			"j.recycling_price",
			"j.nums",
			"0 as no",
		).
		Table("users u, market_user_offer t").
		Joins("left join (?) as j on t.id = j.id", jsonDB).
		Where("t.state = 2").
		Where("t.user_id = u.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams))
}

// 发货
func (d *ItemDao) GetDetailOrder(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.getDetailOrderDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("i.inner_price desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailOrder: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailOrderDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	jsonDB := d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
		JSON_UNQUOTE(t.order_items), 
		'$[*]' COLUMNS(
				name varchar(255) path '$.name',
				state int path '$.state',
				item_id bigint path '$.item_id',
				stock_id int path '$.stock_id'
			)
	)
	`, "\t", " "), "\n", " "))

	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(FROM_UNIXTIME(CAST((t.delivery_time / 1000) AS SIGNED)), '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"300 as log_type",
			"'' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"count(j.item_id) as nums",
			"0 as no",
		).
		Table("users u, `order` t, ? as j, item i", jsonDB).
		Where("u.id = t.user_id").
		Where("t.state in (4,5)").
		Where("j.item_id = i.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams)).
		Group(fmt.Sprintf(
			"DATE_FORMAT(FROM_UNIXTIME(CAST((t.delivery_time / 1000) AS SIGNED)), '%s'), u.id, u.nickname, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price",
			pkg.SQL_DATE_TIME_FORMAT,
		))
}

// 邮件
func (d *ItemDao) GetDetailPublicize(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"a.num as nums",
		).
		Table("publicize_mail_log l,publicize_mail_attachment a, item i").
		Where("l.mail_id = a.config_id").
		Where("a.type = 20").
		Where("a.value = i.id").
		Scopes(database.ScopeQuery(queryParams)).
		Order("i.inner_price desc, a.num desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("getDetailAdmin: %v", err)
		return nil, err
	}
	return
}

func (d *ItemDao) getDetailPublicizeDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"701 as log_type",
			"m.title as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"a.num as nums",
			"0 as no",
		).
		Table("publicize_mail_log l").
		Joins("LEFT JOIN publicize_mail_attachment a ON a.config_id = l.mail_id").
		Joins("LEFT JOIN publicize_mail m ON m.id = l.mail_id").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN item i ON i.id = a.value").
		Where("a.type = 20").
		Scopes(database.ScopeQuery(paramsGroup.PublicizeDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

// 管理员 添加
func (d *ItemDao) GetDetailAdmin(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"gl.level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"count(t.id) as nums",
		).
		Table("cabinet_stock t, item i, gacha_level gl").
		Where("t.theme_type = 999").
		Where("t.item_id = i.id").
		Where("t.level_index=gl.level_index").
		Scopes(database.ScopeQuery(queryParams)).
		Group("i.id, i.name, gl.level_name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price").
		Order("gl.level_name, i.inner_price desc, count(t.id) desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("getDetailAdmin: %v", err)
		return nil, err
	}

	return
}

func (d *ItemDao) getDetailAdminDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"999999 as log_type", // 欧气值
			"'' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"cast(gl.level_name as char) as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"count(t.id) as nums",
			"0 as no",
		).
		Table("cabinet_stock t, users u, item i, gacha_level gl").
		Where("t.theme_type = 999").
		Where("t.user_id = u.id").
		Where("t.item_id = i.id").
		Where("t.level_index=gl.level_index").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Group(fmt.Sprintf("DATE_FORMAT(t.created_at, '%s'), u.id, u.nickname, i.id, i.name, cast(gl.level_name as char), i.cover_thumb, i.show_price, i.inner_price, i.recycling_price", pkg.SQL_DATE_TIME_FORMAT))
}

// 活动
func (d *ItemDao) GetDetailActivityCostAwardConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.award_num as nums",
		).
		Table("activity_cost_award_config t, item i").
		Where("t.award_type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("t.award_value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityCostAward err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityCostAwardDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100002 as log_type", // 欧气值
			"(case when ua.params_3->'$.type'='0' then '兑换' when ua.params_3->'$.type'='1' then '购买' else '' end) as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.award_num*cast(ua.params_2 as signed) as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_cost_award_config t, item i").
		Where("a.key = 'CostAward'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = t.config_id").
		Where("t.award_type = 20").
		Where("t.award_value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivityCostRankConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.num as nums",
		).
		Table("activity_cost_rank_award_config t, item i").
		Where("t.type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("t.value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityCostRank err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityCostRankDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	db := d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"ua.params",
			"t.period",
			"t.no_limit",
			fmt.Sprintf("ROW_NUMBER() OVER(PARTITION BY DATE_FORMAT(ua.created_at, '%s'), u.id, u.nickname, ua.params, t.period ORDER BY no_limit) as near_no", pkg.SQL_DATE_TIME_FORMAT),
		).
		Table("activity a, user_activity ua, users u, activity_cost_rank_award_config t").
		Where("a.key = 'CostRank'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(substring_index(ua.params, '|', 1) as SIGNED) = t.period").
		Where("cast(substring_index(ua.params, '|', -1) as SIGNED) <= t.no_limit").
		Where("t.type = 20").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Group(fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s'), u.id, u.nickname, ua.params, t.period, t.no_limit", pkg.SQL_DATE_TIME_FORMAT))

	return d.center.
		Select(
			"date_time",
			"user_id",
			"user_name",
			"100003 as log_type", // 欧气排名
			"t.params as log_type_name",
			"t.period as period",
			"0 as box_out_no",
			"0 as amount",
			"crac.value as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"crac.num as nums",
			"0 as no",
		).
		Table("(?) t, activity_cost_rank_award_config crac, item i", db).
		Where("t.near_no = 1").
		Where("t.period = crac.period").
		Where("t.no_limit = crac.no_limit").
		Where("crac.type = 20").
		Where("crac.value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivityItemExchange(queryParams database.QueryWhereGroup) (dataCreator, dataOfferer []*MarketItemDetail, err error) {
	dataCreator, err = d.GetDetailActivityItemExchangeOut(queryParams)
	if err != nil {
		return
	}

	dataOfferer, err = d.GetDetailActivityItemExchangeInto(queryParams)
	if err != nil {
		return
	}

	return
}

func (d *ItemDao) GetDetailActivityItemExchangeOut(queryParams database.QueryWhereGroup) (data []*MarketItemDetail, err error) {
	err = d.getDetailActivityItemExchangeOutDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("inner_price desc, nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityItemExchangeOut: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityItemExchangeOutDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	jsonDB := d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
        JSON_UNQUOTE(JSON_EXTRACT(t.params_3, '$.stock_id')),
        '$[*]' COLUMNS(
            stock_id BIGINT PATH '$'
        )
    )
`, "\t", " "), "\n", " "))
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100004 as log_type",
			"'消耗' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"count(j.stock_id) as nums",
			"0 as no",
		).
		Table("activity as a, `user_activity` t, ? as j, users u, cabinet_stock as cs, item i", jsonDB).
		Where("a.key = 'ItemExchange'").
		Where("a.id = t.activity_id").
		Where("t.user_id = u.id").
		Where("j.stock_id = cs.id").
		Where("cs.item_id = i.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams)).
		Group(fmt.Sprintf("DATE_FORMAT(t.created_at, '%s'), u.id, u.nickname, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price", pkg.SQL_DATE_TIME_FORMAT))
}

func (d *ItemDao) GetDetailActivityItemExchangeInto(queryParams database.QueryWhereGroup) (data []*MarketItemDetail, err error) {
	err = d.getDetailActivityItemExchangeIntoDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("inner_price desc, nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityItemExchangeInto: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityItemExchangeIntoDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100004 as log_type",
			"'获得' as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.params_2 as nums",
			"0 as no",
		).
		Table("activity as a, `user_activity` t, users u, activity_item_exchange_config as iec, item i").
		Where("a.key = 'ItemExchange'").
		Where("a.id = t.activity_id").
		Where("t.user_id = u.id").
		Where("t.params = iec.id").
		Where("iec.item_id = i.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams)).
		Group(fmt.Sprintf("DATE_FORMAT(t.created_at, '%s'), u.id, u.nickname, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price, t.params_2", pkg.SQL_DATE_TIME_FORMAT))
}

func (d *ItemDao) GetDetailActivityPrizeWheelConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.getDetailActivityPrizeWheelDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("inner_price desc, nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityPrizeWheelConfig: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityPrizeWheelDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100005 as log_type",
			"apwc.name as log_type_name",
			"apwc.period as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"count(i.id) as nums",
			"0 as no",
		).
		Table("users u,activity_prize_wheel_history t, activity_prize_wheel_award_config apwac, item i,activity_prize_wheel_config apwc").
		Where("t.award_id = apwac.id").
		Where("apwac.type = 20").
		Where("apwac.value = i.id").
		Where("t.user_id = u.id").
		Where("t.config_id = apwc.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams)).
		Group(fmt.Sprintf("DATE_FORMAT(t.created_at, '%s'), u.id, u.nickname, apwc.name, apwc.period, i.id, i.name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price", pkg.SQL_DATE_TIME_FORMAT))
}

func (d *ItemDao) GetDetailActivityStepByStepConfig(levelType int, queryParams database.QueryWhereGroup) (data []*Item, err error) {
	if levelType == 0 {
		err = d.center.
			Select(
				"i.id as item_id",
				"i.name as item_name",
				"i.cover_thumb",
				"i.show_price",
				"i.inner_price",
				"i.recycling_price",
				"t.num as nums",
			).
			Table("activity_step_by_step_award_config t, item i").
			Where("t.type = 20").
			Scopes(database.ScopeQuery(queryParams)).
			Where("t.value = i.id").
			Order("inner_price desc, nums desc").
			Find(&data).Error
	} else if levelType == 1 {
		err = d.center.
			Select(
				"i.id as item_id",
				"i.name as item_name",
				"i.cover_thumb",
				"i.show_price",
				"i.inner_price",
				"i.recycling_price",
				"t.num as nums",
			).
			Table("activity_step_by_step_rank_award_log t, item i").
			Where("t.type = 20").
			Scopes(database.ScopeQuery(queryParams)).
			Where("t.value = i.id").
			Order("inner_price desc, nums desc").
			Find(&data).Error
	}

	if err != nil {
		d.logger.Errorf("GetDetailActivityStepByStepConfig err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityStepByStepDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100006 as log_type", // 步步高升
			"CONCAT(JSON_UNQUOTE(ua.params_3->'$.step_no'), '层 ', JSON_UNQUOTE(ua.params_3->'$.cell_no'), '格') as log_type_name",
			"JSON_UNQUOTE(ua.params_3->'$.config_id') as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"c.num as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_step_by_step_award_config c, item i").
		Where("a.key = 'StepByStep'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = c.cell_config_id").
		Where("c.type = 20").
		Where("c.value = i.id").
		Where("JSON_EXTRACT(ua.params_3, '$.type') = ?", "1"). //1步步高升开奖
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) getDetailActivityStepByStepRankDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100006 as log_type", // 步步高升
			"CONCAT('第', ua.params_2, '名奖励') as log_type_name",
			"JSON_UNQUOTE(ua.params_3->'$.config_id') as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"c.num as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_step_by_step_rank_award_log c, item i").
		Where("a.key = 'StepByStep'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = c.config_id").
		Where("c.type = 20").
		Where("c.value = i.id").
		Where("JSON_EXTRACT(ua.params_3, '$.type') = ?", "2").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivitySignInConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.num as nums",
		).
		Table("activity_sign_in_day_config t, item i").
		Where("t.type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("t.value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityStepByStepConfig err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivitySignInDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100007 as log_type", // 签到
			"CONCAT(JSON_UNQUOTE(ua.params_3->'$.day_no'), '天') as log_type_name",
			"ua.params as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"c.num as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_sign_in_day_config c, item i").
		Where("a.key = 'SignIn'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("JSON_UNQUOTE(ua.params_3->'$.value') = c.value").
		Where("cast(ua.params as SIGNED) = c.config_id").
		Where("c.type = 20").
		Where("c.value = i.id").
		Where("c.deleted_at is NULL").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivityLuckyNumConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.num as nums",
		).
		Table("activity_lucky_num_award_config t, item i").
		Where("t.type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("t.value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityLuckyNumConfig err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityLuckyNumDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100008 as log_type", // 幸运数
			"tc.name as log_type_name",
			"0 as period", // id 作期数
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"c.num as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_lucky_num_target_config tc, activity_lucky_num_award_config c, item i").
		Where("a.key = 'LuckyNum'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = tc.id").
		Where("cast(ua.params as SIGNED) = c.target_id").
		Where("c.type = 20").
		Where("c.value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivityRedemptionCodeConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"t.num as nums",
		).
		Table("activity_redemption_code_award_log t, item i").
		Where("t.type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("t.value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityRedemptionCodeConfig err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityRedemptionCodeDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(ua.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"100010 as log_type",
			"ua.params_2 as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"c.num as nums",
			"0 as no",
		).
		Table("activity a, user_activity ua, users u, activity_redemption_code_log tc, activity_redemption_code_award_log c, item i").
		Where("a.key = 'RedemptionCode'").
		Where("a.id = ua.activity_id").
		Where("ua.user_id = u.id").
		Where("cast(ua.params as SIGNED) = tc.id").
		Where("cast(ua.params as SIGNED) = c.config_id").
		Where("c.type = 20").
		Where("c.value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.ActivityDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}

func (d *ItemDao) GetDetailActivityLotteryConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.getDetailActivityLotteryDB(AllRequestParamsGroup{OtherParams: queryParams}).Order("inner_price desc, nums desc").Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailActivityLotteryConfig: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailActivityLotteryDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(t.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT), //格式化时间戳
			"u.id as user_id",
			"u.nickname as user_name",
			"100011 as log_type",
			"concat_ws(' ', config.name, t.period) as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"award.num as nums",
			"0 as no",
		).
		Table("users u, activity_lottery_config config, activity_lottery_history t, activity_lottery_award award, item i").
		Where("u.id = t.user_id").
		Where("t.config_id = config.id").
		Where("t.config_id = award.config_id").
		Where("award.deleted_at is null").
		Where("award.type = 20").
		Where("award.value = i.id").
		Scopes(database.ScopeQuery(paramsGroup.DateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams)).
		Scopes(database.ScopeQuery(paramsGroup.OtherParams))
}

// 任务
func (d *ItemDao) GetDetailTaskConfig(queryParams database.QueryWhereGroup) (data []*Item, err error) {
	err = d.center.
		Select(
			"i.id as item_id",
			"i.name as item_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"a.num as nums",
		).
		Table("task_user_log l").
		Joins(`
			LEFT JOIN JSON_TABLE(
				JSON_EXTRACT(l.params_2, '$.award_id'),
				'$[*]' COLUMNS (
					award_id BIGINT PATH '$'
				)
			) AS jt ON TRUE
		`).
		Joins("LEFT JOIN task_award a ON a.id = jt.award_id").
		Joins("LEFT JOIN item i ON i.id = a.value").
		Where("a.type = 20").
		Scopes(database.ScopeQuery(queryParams)).
		Where("a.value = i.id").
		Order("inner_price desc, nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetailTaskConfig err: %v", err)
		return
	}

	return
}

func (d *ItemDao) getDetailTaskDB(paramsGroup AllRequestParamsGroup) *gorm.DB {
	return d.center.
		Select(
			fmt.Sprintf("DATE_FORMAT(l.created_at, '%s') as date_time", pkg.SQL_DATE_TIME_FORMAT),
			"u.id as user_id",
			"u.nickname as user_name",
			"200000 as log_type",
			"t.name as log_type_name",
			"0 as period",
			"0 as box_out_no",
			"0 as amount",
			"i.id as item_id",
			"i.name as item_name",
			"'' as level_name",
			"i.cover_thumb",
			"i.show_price",
			"i.inner_price",
			"i.recycling_price",
			"a.num as nums",
			"0 as no",
		).
		Table("task_user_log l").
		Joins("LEFT JOIN users u ON u.id = l.user_id").
		Joins("LEFT JOIN task t ON t.id = l.task_id").
		Joins(`
			LEFT JOIN JSON_TABLE(
				JSON_EXTRACT(l.params_2, '$.award_id'),
				'$[*]' COLUMNS (
					award_id BIGINT PATH '$'
				)
			) AS jt ON TRUE
		`).
		Joins("LEFT JOIN task_award a ON a.id = jt.award_id").
		Joins("LEFT JOIN item i ON i.id = a.value").
		Where("a.type = 20").
		Scopes(database.ScopeQuery(paramsGroup.TaskDateTimeParams)).
		Scopes(database.ScopeQuery(paramsGroup.UsersParams)).
		Scopes(database.ScopeQuery(paramsGroup.ItemParams))
}
