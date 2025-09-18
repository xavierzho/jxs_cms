package form

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/item/dao"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type DetailRequest struct {
	ID        string `form:"id" binding:"required"`
	LogType   int    `form:"log_type" binding:"required"`
	LevelType int    `form:"level_type"` // 0,1,2,3,4
}

func (q *DetailRequest) Parse() (queryParams, betQueryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if q.LogType == 100002 { // 欧气值
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.config_id = ?",
			Value:  []any{convert.GetInt64(q.ID)},
		})
	} else if q.LogType == 100003 { // 欧气值排名
		idList := strings.Split(q.ID, "|")
		if len(idList) != 3 {
			return nil, nil, fmt.Errorf("id is invalid")
		}

		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.period = ? and t.no_limit = ?",
			Value:  []any{idList[0], idList[2]},
		})
	} else if q.LogType == 999999 { // 管理员添加
		idList := strings.Split(q.ID, "|")
		if len(idList) != 2 {
			return nil, nil, fmt.Errorf("id is invalid")
		}

		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.create_time = ? and t.user_id = ?",
			Value:  []any{idList[0], idList[1]},
		})
	} else if q.LogType == 300 { //发货
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.id = ?",
			Value:  []any{convert.GetInt64(q.ID)},
		})
	} else if q.LogType == 100005 { //转盘抽奖
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.id = ?",
			Value:  []any{convert.GetInt64(q.ID)},
		})

	} else {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "t.id = ?",
			Value:  []any{convert.GetInt64(q.ID)},
		})
	}

	betQueryParams = append(betQueryParams, database.QueryWhere{
		Prefix: "j.level_type = ?",
		Value:  []any{q.LevelType},
	})

	return
}

func (q *DetailRequest) Valid() (err error) {
	switch q.LogType {
	case 101, 102, 103, 104:
	case 200:
	case 300:
	case 100002, 100003, 100004, 100005:
	case 999999:
	default:
		return fmt.Errorf("not expected log type: %d", q.LogType)
	}

	switch q.LevelType {
	case 0, 1, 2, 3, 4:
	default:
		return fmt.Errorf("not expected level type: %d", q.LevelType)
	}

	return nil
}

type DetailAllRequest struct {
	LogAllRequest
}

func (q *DetailAllRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	if dateTimeRange, paramsGroup, err = q.LogAllRequest.Parse(); err != nil {
		return
	}
	paramsGroup.DateTimeParams = append(paramsGroup.DateTimeParams, database.QueryWhere{
		Prefix: "t.created_at between ? and ?",
		Value:  []any{dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second - time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)},
	})

	paramsGroup.ActivityDateTimeParams = append(paramsGroup.ActivityDateTimeParams, database.QueryWhere{
		Prefix: "ua.created_at between ? and ?",
		Value:  []any{dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second - time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)},
	})

	return
}

type BetItemDetail struct {
	// GachaName string `json:"gacha_name"`
	BoxOutNo int64 `json:"box_out_no"`
	// BetNums   int    `json:"bet_nums"`
	cForm.Item
	Nums int `json:"nums"`
	No   int `json:"no"`
}

func FormatBetItemDetail(data []*dao.BetItemDetail) (result []*BetItemDetail) {
	for _, item := range data {
		result = append(result, &BetItemDetail{
			// GachaName: item.GachaName,
			BoxOutNo: item.BoxOutNo,
			// BetNums:   item.BetNums,
			Item: cForm.Item{
				ItemID:         strconv.FormatInt(item.Item.ItemID, 10),
				ItemName:       item.Item.ItemName,
				LevelName:      item.Item.LevelName,
				CoverThumb:     item.Item.CoverThumb,
				ShowPrice:      util.ConvertAmount2Decimal(item.Item.ShowPrice),
				InnerPrice:     util.ConvertAmount2Decimal(item.Item.InnerPrice),
				RecyclingPrice: util.ConvertAmount2Decimal(item.Item.RecyclingPrice),
			},
			Nums: item.Item.Nums,
			No:   item.Item.No,
		})
	}

	return
}

type MarketItemDetail struct {
	UserID   int64           `json:"user_id"`
	UserName string          `json:"user_name"`
	Amount   decimal.Decimal `json:"amount"`
	cForm.Item
	Nums int `json:"nums"`
}

func FormatMarketItemDetail(data []*dao.MarketItemDetail) (result []*MarketItemDetail) {
	for _, item := range data {
		result = append(result, &MarketItemDetail{
			UserID:   item.UserID,
			UserName: item.UserName,
			Amount:   util.ConvertAmount2Decimal(item.Amount),
			Item: cForm.Item{
				ItemID:         strconv.FormatInt(item.Item.ItemID, 10),
				ItemName:       item.Item.ItemName,
				LevelName:      item.Item.LevelName,
				CoverThumb:     item.Item.CoverThumb,
				ShowPrice:      util.ConvertAmount2Decimal(item.Item.ShowPrice),
				InnerPrice:     util.ConvertAmount2Decimal(item.Item.InnerPrice),
				RecyclingPrice: util.ConvertAmount2Decimal(item.Item.RecyclingPrice),
			},
			Nums: item.Item.Nums,
		})
	}

	return
}

type Item struct {
	cForm.Item
	Nums int `json:"nums"`
}

func FormatItem(data []*dao.Item) (result []*Item) {
	for _, item := range data {
		result = append(result, &Item{
			Item: cForm.Item{
				ItemID:         strconv.FormatInt(item.Item.ItemID, 10),
				ItemName:       item.Item.ItemName,
				LevelName:      item.Item.LevelName,
				CoverThumb:     item.Item.CoverThumb,
				ShowPrice:      util.ConvertAmount2Decimal(item.Item.ShowPrice),
				InnerPrice:     util.ConvertAmount2Decimal(item.Item.InnerPrice),
				RecyclingPrice: util.ConvertAmount2Decimal(item.Item.RecyclingPrice),
			},
			Nums: item.Nums,
		})
	}

	return
}

type ItemDetail struct {
	DateTime    string `json:"date_time"`
	UserID      int64  `json:"user_id"`
	UserName    string `json:"user_name"`
	LogType     int    `json:"log_type"`
	LogTypeStr  string `json:"log_type_str"`
	LogTypeName string `json:"log_type_name"`
	BoxOutNo    string `json:"box_out_no"`
	cForm.Item
	Nums int `json:"nums"`
	No   int `json:"no"`
}

func FormatItemDetail(ctx context.Context, data []*dao.ItemDetail) (result []*ItemDetail, err error) {
	for _, item := range data {
		var logTypeName = item.LogTypeName
		if item.Period != 0 {
			if item.LogType == 100003 { // 消费排行
				logTypeNameList := strings.Split(item.LogTypeName, "|")
				if len(logTypeNameList) != 2 {
					return nil, fmt.Errorf("100003's LogTypeName invalid: %s", item.LogTypeName)
				}
				logTypeName = fmt.Sprintf("第%d期 第%s名", item.Period, logTypeNameList[1])
			} else {
				logTypeName = strings.Join([]string{item.LogTypeName, fmt.Sprintf("第%d期", item.Period)}, " ")
			}
		}

		result = append(result, &ItemDetail{
			DateTime:    item.DateTime,
			UserID:      item.UserID,
			UserName:    item.UserName,
			LogType:     item.LogType,
			LogTypeStr:  global.I18n.T(ctx, "source_type", convert.GetString(item.LogType)),
			LogTypeName: logTypeName,
			BoxOutNo:    strconv.FormatInt(item.BoxOutNo, 10),
			Item: cForm.Item{
				ItemID:         strconv.FormatInt(item.Item.ItemID, 10),
				ItemName:       item.Item.ItemName,
				LevelName:      item.Item.LevelName,
				CoverThumb:     item.Item.CoverThumb,
				ShowPrice:      util.ConvertAmount2Decimal(item.Item.ShowPrice),
				InnerPrice:     util.ConvertAmount2Decimal(item.Item.InnerPrice),
				RecyclingPrice: util.ConvertAmount2Decimal(item.Item.RecyclingPrice),
			},
			Nums: item.Nums,
			No:   item.No,
		})
	}

	return
}

func FormatItemDetail2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.ItemDetail) (excelModel *excel.Excel[*ItemDetail], err error) {
	data, err := FormatItemDetail(ctx, _data)
	if err != nil {
		return
	}

	reflectMap := map[string]func(source *ItemDetail) any{
		"时间":    func(source *ItemDetail) any { return source.DateTime },
		"用户ID":  func(source *ItemDetail) any { return source.UserID },
		"用户昵称":  func(source *ItemDetail) any { return source.UserName },
		"项目类型":  func(source *ItemDetail) any { return source.LogTypeStr },
		"项目名称":  func(source *ItemDetail) any { return source.LogTypeName },
		"箱号":    func(source *ItemDetail) any { return source.BoxOutNo },
		"物品ID":  func(source *ItemDetail) any { return source.ItemID },
		"物品名称":  func(source *ItemDetail) any { return source.ItemName },
		"物品等级":  func(source *ItemDetail) any { return source.LevelName },
		"物品封面图": func(source *ItemDetail) any { return source.CoverThumb },
		"物品展示价": func(source *ItemDetail) any { return source.ShowPrice },
		"物品成本价": func(source *ItemDetail) any { return source.InnerPrice },
		"物品回收价": func(source *ItemDetail) any { return source.RecyclingPrice },
		"数量":    func(source *ItemDetail) any { return source.Nums },
		"第几抽":   func(source *ItemDetail) any { return source.No },
	}

	excelModel = &excel.Excel[*ItemDetail]{
		FileName:   fmt.Sprintf("user_item_detail_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"用户物品详情日志"},
		SheetNameWithHead: map[string][]string{
			"用户物品详情日志": {
				"时间",
				"用户ID", "用户昵称",
				"项目类型", "项目名称", "箱号",
				"物品ID", "物品名称", "物品等级", "物品封面图",
				"物品展示价", "物品成本价", "物品回收价",
				"数量", "第几抽",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*ItemDetail]{
			"用户物品详情日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*ItemDetail]{
			"用户物品详情日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
