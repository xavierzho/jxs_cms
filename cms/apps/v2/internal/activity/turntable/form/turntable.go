package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/activity/turntable/dao"
	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 奖品类型
type Type int64

const (
	PrizeType_None   Type = 0
	PrizeType_Coin   Type = 1
	PrizeType_Coupon Type = 10
	PrizeType_Goods  Type = 20
)

type GenerateRequest struct {
	iForm.DateRangeRequest
}

type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateRange [2]time.Time, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateRangeRequest
	cForm.UserInfoRequest
	Type Type   `form:"type"`
	Name string `form:"name"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateRange, err = q.DateRangeRequest.Parse(); err != nil {
		return
	}

	return dateRange, nil
}

func (q *AllRequest) Valid() error {
	switch q.Type {
	case PrizeType_Coin, PrizeType_None:
	case PrizeType_Coupon, PrizeType_Goods:
	default:
		return fmt.Errorf("not expected Type: %v", q.Type)
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
	PointType     int64           `json:"point_type"`      //抽奖支付类型
	PointTypeName string          `json:"point_type_name"` //抽奖支付类型名称
	Point         int64           `json:"point"`           //抽奖消耗
	PrizeValue    decimal.Decimal `json:"prize_value"`     // 奖品价值
}

func Format(dateRange [2]time.Time, _summary map[string]any, data []*dao.Turntable) (summary map[string]any, result []*Turntable) {
	summary = _summary
	if summary != nil {
		summary["amount"] = util.ConvertAmount2Decimal(summary["amount"])
	}
	PointTypeName := map[int64]string{
		10: "现金点",
		11: "欧气值",
	}

	TypeName := map[int64]string{
		0:  "潮币",
		10: "优惠券",
		20: "物品",
	}
	var ConvertPoint int64
	var totalAmount decimal.Decimal
	for _, item := range data {
		if item.PointType == 10 {
			ConvertPoint = item.Point / 100 // 现金点1元=100点
		} else {
			ConvertPoint = item.Point / 10 // 欧气值1元=10点
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
			TypeName:      TypeName[item.Type],
			ItemId:        item.ItemId,
			ItemName:      item.ItemName,
			PointType:     item.PointType,
			PointTypeName: PointTypeName[item.PointType],
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

func Format2Excel(dateRange [2]time.Time, _data []*dao.Turntable) (excelModel *excel.Excel[*Turntable], err error) {
	_, data := Format(dateRange, nil, _data)

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
