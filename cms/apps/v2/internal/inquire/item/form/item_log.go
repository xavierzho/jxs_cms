package form

import (
	"context"
	"fmt"
	"strings"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/item/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 范围兼容 只输入一个; 前端设置值为 2**32 占位
type LogRequest struct {
	*app.Pager
	LogAllRequest
}

func (q *LogRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	q.Pager.Parse()
	return q.LogAllRequest.Parse()
}

type LogAllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest       // UsersParams
	LogTypeList           []int `form:"log_type_list[]"` // dbBranch && GachaParams
	betFlag               bool
	marketFlag            bool
	activityFlag          dao.ActivityFlag
	adminFlag             bool
	orderFlag             bool
	taskFlag              bool
	publicizeFlag         bool
	GachaName             string      `form:"gacha_name"`            // 仅 GachaParams
	UpdateAmountRange     *[2]float64 `form:"update_amount_range[]"` // 仅 AmountParams
	ShowPriceRange        *[2]float64 `form:"show_price_range[]"`    // 仅 ItemParams
	InnerPriceRange       *[2]float64 `form:"inner_price_range[]"`   // 仅 ItemParams
}

func (q *LogAllRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	// UsersParams
	if paramsGroup.UsersParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	paramsGroup.BetFlag = q.betFlag
	paramsGroup.MarketFlag = q.marketFlag
	paramsGroup.ActivityFlag = q.activityFlag
	paramsGroup.AdminFlag = q.adminFlag
	paramsGroup.OrderFlag = q.orderFlag
	paramsGroup.TaskFlag = q.taskFlag
	paramsGroup.PublicizeFlag = q.publicizeFlag

	// GachaParams
	{
		if q.GachaName != "" {
			paramsGroup.GachaParams = append(paramsGroup.GachaParams, database.QueryWhere{
				Prefix: "gacha_name = ?",
				Value:  []any{q.GachaName},
			})
		}
		if len(q.LogTypeList) != 0 {
			typeList := []int{}
			for _, i := range q.LogTypeList {
				switch i {
				case 101, 102, 103, 104, 105, 106:
					typeList = append(typeList, i%100)
				}
			}

			if len(typeList) != 0 {
				paramsGroup.GachaParams = append(paramsGroup.GachaParams, database.QueryWhere{
					Prefix: "gacha_type in ?",
					Value:  []any{typeList},
				})
			}
		}
	}

	// AmountParams
	{
		if q.UpdateAmountRange != nil {
			paramsGroup.AmountParams = append(paramsGroup.AmountParams, database.QueryWhere{
				Prefix: "update_amount between ? and ?",
				Value:  []any{util.ReconvertAmount2Decimal(q.UpdateAmountRange[0]).IntPart(), util.ReconvertAmount2Decimal(q.UpdateAmountRange[1]).IntPart()},
			})
		}
	}

	// ItemParams
	{
		if q.ShowPriceRange != nil {
			paramsGroup.ItemParams = append(paramsGroup.ItemParams, database.QueryWhere{
				Prefix: "show_price between ? and ?",
				Value:  []any{util.ReconvertAmount2Decimal(q.ShowPriceRange[0]).IntPart(), util.ReconvertAmount2Decimal(q.ShowPriceRange[1]).IntPart()},
			})
		}

		if q.InnerPriceRange != nil {
			paramsGroup.ItemParams = append(paramsGroup.ItemParams, database.QueryWhere{
				Prefix: "inner_price between ? and ?",
				Value:  []any{util.ReconvertAmount2Decimal(q.InnerPriceRange[0]).IntPart(), util.ReconvertAmount2Decimal(q.InnerPriceRange[1]).IntPart()},
			})
		}
	}

	return
}

func (q *LogAllRequest) Valid() (err error) {
	if len(q.LogTypeList) == 0 {
		q.betFlag, q.marketFlag, q.adminFlag, q.orderFlag, q.publicizeFlag = true, true, true, true, true
		q.activityFlag = dao.ActivityFlag{
			Flag:           true,
			CostAward:      true,
			CostRank:       true,
			ItemExchange:   true,
			PrizeWheel:     true, // 转盘抽奖
			StepByStep:     true,
			SignIn:         true,
			LuckyNum:       true,
			RedemptionCode: true,
		}
		q.taskFlag = true
	} else {
		for _, i := range q.LogTypeList {
			switch i {
			case 101, 102, 103, 104, 105, 106:
				q.betFlag = true
			case 200:
				q.marketFlag = true
			case 300:
				q.orderFlag = true
			case 701:
				q.publicizeFlag = true
			case 100002:
				q.activityFlag.Flag = true
				q.activityFlag.CostAward = true
			case 100003:
				q.activityFlag.Flag = true
				q.activityFlag.CostRank = true
			case 100004:
				q.activityFlag.Flag = true
				q.activityFlag.ItemExchange = true
			case 100005:
				q.activityFlag.Flag = true
				q.activityFlag.PrizeWheel = true // 转盘抽奖
			case 100006:
				q.activityFlag.Flag = true
				q.activityFlag.StepByStep = true
			case 100007:
				q.activityFlag.Flag = true
				q.activityFlag.SignIn = true
			case 100008:
				q.activityFlag.Flag = true
				q.activityFlag.LuckyNum = true
			case 100010:
				q.activityFlag.Flag = true
				q.activityFlag.RedemptionCode = true
			case 100011:
				q.activityFlag.Flag = true
				q.activityFlag.Lottery = true
			case 200000:
				q.taskFlag = true
			case 999999:
				q.adminFlag = true
			default:
				return fmt.Errorf("not expected log type: %d", q.LogTypeList)
			}
		}
	}

	if q.UpdateAmountRange != nil && q.UpdateAmountRange[1] < q.UpdateAmountRange[0] {
		return fmt.Errorf("invalid UpdateAmountRange: %v", q.UpdateAmountRange)
	}

	if q.ShowPriceRange != nil && q.ShowPriceRange[1] < q.ShowPriceRange[0] {
		return fmt.Errorf("invalid ShowPriceRange: %v", q.ShowPriceRange)
	}

	if q.InnerPriceRange != nil && q.InnerPriceRange[1] < q.InnerPriceRange[0] {
		return fmt.Errorf("invalid InnerPriceRange: %v", q.InnerPriceRange)
	}

	return nil
}

type ItemLog struct {
	ID             string          `json:"id"`
	DateTime       string          `json:"date_time"`
	UserID         int64           `json:"user_id"`
	UserName       string          `json:"user_name"`
	LogType        int             `json:"log_type"`
	LogTypeStr     string          `json:"log_type_str"`
	LogTypeName    string          `json:"log_type_name"`
	BetNums        int             `json:"bet_nums"`
	LevelType      int             `json:"level_type"`
	LevelTypeStr   string          `json:"level_type_str"`
	UpdateAmount   decimal.Decimal `json:"update_amount"`
	ShowPrice      decimal.Decimal `json:"show_price"`
	InnerPrice     decimal.Decimal `json:"inner_price"`
	RecyclingPrice decimal.Decimal `json:"recycling_price"`
}

func FormatLog(ctx context.Context, _summary map[string]any, data []*dao.ItemLog) (summary map[string]any, result []*ItemLog, err error) {
	if _summary != nil {
		summary = _summary
		summary["update_amount"] = util.ConvertAmount2Decimal(_summary["update_amount"])
		summary["show_price"] = util.ConvertAmount2Decimal(_summary["show_price"])
		summary["inner_price"] = util.ConvertAmount2Decimal(_summary["inner_price"])
		summary["recycling_price"] = util.ConvertAmount2Decimal(_summary["recycling_price"])
	}

	for _, item := range data {
		var logTypeName = item.LogTypeName
		if item.Period != 0 {
			switch item.LogType {
			case 100003:
				idStrList := strings.Split(item.ID, "|")
				if len(idStrList) != 3 {
					return nil, nil, fmt.Errorf("100003's id invalid: %s", item.ID)
				}
				logTypeName = fmt.Sprintf("第%d期 第%s名", item.Period, idStrList[1])
			case 100006, 100007:
				logTypeName = strings.Join([]string{fmt.Sprintf("第%d期", item.Period), item.LogTypeName}, " ")
			default:
				logTypeName = strings.Join([]string{item.LogTypeName, fmt.Sprintf("第%d期", item.Period)}, " ")
			}
		}

		result = append(result, &ItemLog{
			ID:             item.ID,
			DateTime:       item.DateTime,
			UserID:         item.UserID,
			UserName:       item.UserName,
			LogType:        item.LogType,
			LogTypeStr:     global.I18n.T(ctx, "source_type", convert.GetString(item.LogType)),
			LogTypeName:    logTypeName,
			BetNums:        item.BetNums,
			LevelType:      item.LevelType,
			LevelTypeStr:   global.I18n.T(ctx, "item.levelType", convert.GetString(item.LevelType)),
			UpdateAmount:   util.ConvertAmount2Decimal(item.UpdateAmount),
			ShowPrice:      util.ConvertAmount2Decimal(item.ShowPrice),
			InnerPrice:     util.ConvertAmount2Decimal(item.InnerPrice),
			RecyclingPrice: util.ConvertAmount2Decimal(item.RecyclingPrice),
		})
	}

	return
}

func FormatLog2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.ItemLog) (excelModel *excel.Excel[*ItemLog], err error) {
	_, data, err := FormatLog(ctx, nil, _data)
	if err != nil {
		return nil, err
	}

	reflectMap := map[string]func(source *ItemLog) any{
		"时间":   func(source *ItemLog) any { return source.DateTime },
		"用户ID": func(source *ItemLog) any { return source.UserID },
		"用户昵称": func(source *ItemLog) any { return source.UserName },
		"项目类型": func(source *ItemLog) any { return source.LogTypeStr },
		"项目名称": func(source *ItemLog) any { return source.LogTypeName },
		"数量":   func(source *ItemLog) any { return source.BetNums },
		"奖品类型": func(source *ItemLog) any { return source.LevelTypeStr },
		"余额变动": func(source *ItemLog) any { return source.UpdateAmount },
		"展示价":  func(source *ItemLog) any { return source.ShowPrice },
		"成本价":  func(source *ItemLog) any { return source.InnerPrice },
		"回收价":  func(source *ItemLog) any { return source.RecyclingPrice },
	}
	excelModel = &excel.Excel[*ItemLog]{
		FileName:   fmt.Sprintf("user_item_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"用户物品日志"},
		SheetNameWithHead: map[string][]string{
			"用户物品日志": {
				"时间", "用户ID", "用户昵称", "项目类型", "项目名称", "数量", "奖品类型",
				"余额变动", "展示价", "成本价", "回收价",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*ItemLog]{
			"用户物品日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*ItemLog]{
			"用户物品日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

type GetRevenueRequest struct {
	cForm.UserInfoRequest // UsersParam
}

type ItemRevenue struct {
	UpdateAmount7Day    decimal.Decimal `json:"update_amount_7_day"`
	RecyclingPrice7Day  decimal.Decimal `json:"recycling_price_7_day"`
	Revenue7Day         decimal.Decimal `json:"revenue_7_day"`
	UpdateAmount15Day   decimal.Decimal `json:"update_amount_15_day"`
	RecyclingPrice15Day decimal.Decimal `json:"recycling_price_15_day"`
	Revenue15Day        decimal.Decimal `json:"revenue_15_day"`
	UpdateAmount30Day   decimal.Decimal `json:"update_amount_30_day"`
	RecyclingPrice30Day decimal.Decimal `json:"recycling_price_30_day"`
	Revenue30Day        decimal.Decimal `json:"revenue_30_day"`
}

func (q *GetRevenueRequest) Parse() (paramsGroup dao.AllRequestParamsGroup, err error) {
	// UsersParams
	if paramsGroup.UsersParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}
	return
}
func FormatRevenue(ctx context.Context, data *dao.ItemRevenue) (result []*ItemRevenue, err error) {
	result = append(result, &ItemRevenue{
		Revenue7Day:         data.Revenue7Day,
		UpdateAmount7Day:    data.UpdateAmount7Day,
		RecyclingPrice7Day:  data.RecyclingPrice7Day,
		Revenue15Day:        data.Revenue15Day,
		UpdateAmount15Day:   data.UpdateAmount15Day,
		RecyclingPrice15Day: data.RecyclingPrice15Day,
		Revenue30Day:        data.Revenue30Day,
		UpdateAmount30Day:   data.UpdateAmount30Day,
		RecyclingPrice30Day: data.RecyclingPrice30Day,
	})
	return
}
