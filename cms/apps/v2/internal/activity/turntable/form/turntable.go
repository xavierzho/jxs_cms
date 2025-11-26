package form

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/activity/turntable/dao"
	cForm "data_backend/apps/v2/internal/common/form"
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

type GenerateRequest struct {
	iForm.DateRangeRequest
}

type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateRangeRequest
	cForm.UserInfoRequest
	Name      string           `form:"name"`
	PointType *cForm.PointType `form:"point_type"`
	Type      cForm.AwardType  `form:"type"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateRange, err = q.DateRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	if q.Name != "" {
		queryParams = append(queryParams, database.QueryWhere{Prefix: "apwc.name = ?", Value: []any{q.Name}})
	}
	//AND ((JSON_UNQUOTE(apwh.params_3->'$.point_type') IS NULL AND apwc.point_type=11) OR JSON_UNQUOTE(apwh.params_3->'$.point_type') = 'PointType_CostAwardPoint')
	if q.PointType != nil {
		switch *q.PointType {
		case cForm.PointType_Coin:
			queryParams = append(queryParams, database.QueryWhere{Prefix: "(JSON_UNQUOTE(apwh.params_3->'$.point_type') = ?)", Value: []any{"PointType_Coin"}})
		case cForm.PointType_AmountPoint:
			queryParams = append(queryParams, database.QueryWhere{Prefix: "((JSON_UNQUOTE(apwh.params_3->'$.point_type') IS NULL AND apwc.point_type= ?) OR JSON_UNQUOTE(apwh.params_3->'$.point_type') = 'PointType_AmountPoint')", Value: []any{int32(*q.PointType)}})
		case cForm.PointType_Free:
			queryParams = append(queryParams, database.QueryWhere{Prefix: "(JSON_UNQUOTE(apwh.params_3->'$.point_type') = ?)", Value: []any{"PointType_Free"}})
		case cForm.PointType_CostAwardPoint:
			queryParams = append(queryParams, database.QueryWhere{Prefix: "((JSON_UNQUOTE(apwh.params_3->'$.point_type') IS NULL AND apwc.point_type= ?) OR JSON_UNQUOTE(apwh.params_3->'$.point_type') = 'PointType_CostAwardPoint')", Value: []any{int32(*q.PointType)}})
		case cForm.PointType_Gold:
			queryParams = append(queryParams, database.QueryWhere{Prefix: "(JSON_UNQUOTE(apwh.params_3->'$.point_type') = ?)", Value: []any{"PointType_Gold"}})
		}

	}
	if q.Type != 0 {
		queryParams = append(queryParams, database.QueryWhere{Prefix: "apwac.type = ?", Value: []any{int32(q.Type)}})
	}

	return dateRange, queryParams, nil
}

func (q *AllRequest) Valid() error {
	if q.PointType != nil {
		if err := q.PointType.Valid(); err != nil {
			return err
		}
	}

	if err := q.Type.Valid(); err != nil {
		return err
	}

	return nil
}

type Turntable struct {
	Date          string          `json:"date"`
	UserID        int64           `json:"user_id"`
	UserName      string          `json:"user_name"`
	Period        int64           `json:"period"`
	Name          string          `json:"name"`
	ItemId        int64           `json:"item_id"`
	Type          int64           `json:"type"`
	TypeName      string          `json:"type_name"`
	ItemName      string          `json:"item_name"`
	PointType     string          `json:"point_type"`      //抽奖支付类型
	PointTypeName string          `json:"point_type_name"` //抽奖支付类型名称
	Point         string          `json:"point"`           //抽奖消耗
	PrizeValue    decimal.Decimal `json:"prize_value"`     // 奖品价值
}

func Format(ctx context.Context, dateRange [2]time.Time, _summary map[string]any, data []*dao.Turntable) (summary map[string]any, result []*Turntable) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	}

	var ConvertPoint string
	var totalAmount decimal.Decimal
	for _, item := range data {
		switch item.PointType {
		case "PointType_AmountPoint":
			item.PointType = "10"
			ConvertPoint = util.ConvertAmount2Decimal(item.Point).String()
		case "PointType_CostAwardPoint":
			item.PointType = "11"
			ConvertPoint = util.ConvertAmount2Decimal(item.Point).Mul(cForm.COST_AWARD_POINT_STEP).String()
		case "PointType_Coin":
			item.PointType = "0"
			ConvertPoint = util.ConvertAmount2Decimal(item.Point).String()
		case "PointType_Gold":
			item.PointType = "2"
			ConvertPoint = util.ConvertAmount2Decimal(item.Point).String()
		case "PointType_Free":
			item.PointType = "100" //免费
			ConvertPoint = "0"
		default:
			if item.OldPointType == 10 {
				item.PointType = "10"
				ConvertPoint = util.ConvertAmount2Decimal(item.Point).String()
			} else if item.OldPointType == 11 {
				item.PointType = "11"
				ConvertPoint = util.ConvertAmount2Decimal(item.Point).Mul(cForm.COST_AWARD_POINT_STEP).String()
			}
		}

		PrizeValue := util.ConvertAmount2Decimal(item.PrizeValue)
		if item.Type == 0 {
			PrizeValue = util.ConvertAmount2Decimal(item.ItemId)
			totalAmount = totalAmount.Add(PrizeValue)
		}
		result = append(result, &Turntable{
			Date:          item.Date,
			UserID:        item.UserID,
			UserName:      item.UserName,
			Period:        item.Period,
			Name:          item.Name,
			Type:          item.Type,
			TypeName:      global.I18n.T(ctx, "activity.award_type", convert.GetString(item.Type)),
			ItemId:        item.ItemId,
			ItemName:      item.ItemName,
			PointType:     item.PointType,
			PointTypeName: global.I18n.T(ctx, "activity.point_type", item.PointType),
			Point:         ConvertPoint,
			PrizeValue:    PrizeValue,
		})
	}
	if summary != nil {
		totalAmount = totalAmount.Add(util.ConvertAmount2Decimal(summary["total_amount"]))
		summary["total_amount"] = totalAmount
	}
	return
}

func Format2Excel(ctx context.Context, dateRange [2]time.Time, _data []*dao.Turntable) (excelModel *excel.Excel[*Turntable], err error) {
	_, data := Format(ctx, dateRange, nil, _data)

	reflectMap := map[string]func(source *Turntable) any{
		"日期":   func(source *Turntable) any { return source.Date },
		"用户ID": func(source *Turntable) any { return source.UserID },
		"用户昵称": func(source *Turntable) any { return source.UserName },
		"项目名称": func(source *Turntable) any { return source.Name },
		"期数":   func(source *Turntable) any { return source.Period },
		"消耗类型": func(source *Turntable) any { return source.PointTypeName },
		"消耗数量": func(source *Turntable) any { return source.Point },
		"奖品类型": func(source *Turntable) any { return source.TypeName },
		"奖品ID": func(source *Turntable) any { return source.ItemId },
		"奖品名称": func(source *Turntable) any { return source.ItemName },
		"奖品价值": func(source *Turntable) any { return source.PrizeValue },
	}

	excelModel = &excel.Excel[*Turntable]{
		FileName:   fmt.Sprintf("Turntable_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"转盘抽奖记录"},
		SheetNameWithHead: map[string][]string{
			"转盘抽奖记录": {
				"日期",
				"用户ID", "用户昵称", "项目名称", "期数", "消耗类型", "消耗数量", "奖品类型", "奖品ID", "奖品名称", "奖品价值",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*Turntable]{
			"转盘抽奖记录": data,
		},
		ReflectMap: map[string]excel.RowReflect[*Turntable]{
			"转盘抽奖记录": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
