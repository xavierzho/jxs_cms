package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type LogListRequest struct {
	app.Pager
	LogAllRequest
}

func (q *LogListRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.LogAllRequest.Parse()
}

type LogAllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	PointType cForm.PointType `form:"point_type"`
}

func (q *LogAllRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if err = q.PointType.Valid(); err != nil {
		return
	}

	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	if q.PointType != 0 {
		queryParams = append(queryParams, database.QueryWhere{Prefix: "sc.point_type = ?", Value: []any{int32(q.PointType)}})
	}

	return
}

type StepByStepLog struct {
	ID         int64           `json:"id"`
	CreatedAt  string          `json:"created_at"`
	UserID     int64           `json:"user_id"`
	UserName   string          `json:"user_name"`
	Period     int64           `json:"period"`
	PointType  cForm.PointType `json:"point_type"`
	Point      decimal.Decimal `json:"point"`
	StepNO     uint32          `json:"step_no"`
	CellNO     uint32          `json:"cell_no"`
	InnerPrice decimal.Decimal `json:"inner_price"`
}

func LogFormat(pointType cForm.PointType, data []map[string]any, _summary map[string]any) (result []*StepByStepLog, summary map[string]any) {
	summary = _summary
	if summary != nil {
		switch pointType {
		case cForm.PointType_AmountPoint:
			summary["point"] = util.ConvertAmount2Decimal(summary["point"])
		case cForm.PointType_CostAwardPoint:
			summary["point"] = util.ConvertAmount2Decimal(summary["point"]).Mul(cForm.COST_AWARD_POINT_STEP)
		default:
			summary["point"] = decimal.Zero
		}
		summary["inner_price"] = util.ConvertAmount2Decimal(summary["inner_price"])
	}

	for _, item := range data {
		pointType := cForm.PointType(convert.GetInt32(item["point_type"]))
		var point decimal.Decimal

		switch pointType {
		case cForm.PointType_AmountPoint:
			point = util.ConvertAmount2Decimal(item["point"])
		case cForm.PointType_CostAwardPoint:
			point = util.ConvertAmount2Decimal(item["point"]).Mul(cForm.COST_AWARD_POINT_STEP)
		}

		result = append(result, &StepByStepLog{
			ID:         convert.GetInt64(item["id"]),
			CreatedAt:  convert.GetString(item["created_at"]),
			UserID:     convert.GetInt64(item["user_id"]),
			UserName:   convert.GetString(item["user_name"]),
			Period:     convert.GetInt64(item["period"]),
			PointType:  pointType,
			Point:      point,
			StepNO:     convert.GetUint32(item["step_no"]),
			CellNO:     convert.GetUint32(item["cell_no"]),
			InnerPrice: util.ConvertAmount2Decimal(item["inner_price"]),
		})
	}

	return
}

func LogFormat2Excel(dateTimeRange [2]time.Time, _data []map[string]any) (excelModel *excel.Excel[*StepByStepLog], err error) {
	data, _ := LogFormat(cForm.PointType_AmountPoint, _data, nil)

	reflectMap := map[string]func(source *StepByStepLog) any{
		"时间":   func(source *StepByStepLog) any { return source.CreatedAt },
		"用户ID": func(source *StepByStepLog) any { return source.UserID },
		"用户昵称": func(source *StepByStepLog) any { return source.UserName },
		"期数":   func(source *StepByStepLog) any { return source.Period },
		"层号":   func(source *StepByStepLog) any { return source.StepNO },
		"格号":   func(source *StepByStepLog) any { return source.CellNO },
	}

	excelModel = &excel.Excel[*StepByStepLog]{
		FileName:   fmt.Sprintf("step_by_step_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"步步高升日志"},
		SheetNameWithHead: map[string][]string{
			"步步高升日志": {
				"时间", "用户ID", "用户昵称",
				"期数", "层号", "格号",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*StepByStepLog]{
			"步步高升日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*StepByStepLog]{
			"步步高升日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

type DetailRequest struct {
	CellConfigId int64 `form:"cell_config_id"`
}

func (q *DetailRequest) Parse() (queryParams database.QueryWhereGroup, err error) {
	queryParams = append(queryParams, database.QueryWhere{
		Prefix: "ac.cell_config_id = ?",
		Value:  []any{q.CellConfigId},
	})

	return
}

type DetailAllRequest struct {
	LogAllRequest
}

type StepByStepAward struct {
	CreatedAt   string          `json:"created_at"`
	UserID      int64           `json:"user_id"`
	UserName    string          `json:"user_name"`
	Period      int64           `json:"period"`
	StepNO      uint32          `json:"step_no"`
	CellNO      uint32          `json:"cell_no"`
	AwardType   cForm.AwardType `json:"award_type"`
	AwardValue  string          `json:"award_value"`
	AwardName   string          `json:"award_name"`
	AwardParams string          `json:"award_params"`
	AwardNum    int             `json:"award_num"`
	InnerPrice  decimal.Decimal `json:"inner_price"`
}

func DetailFormat(data []map[string]any) (result []*StepByStepAward, err error) {
	for _, item := range data {
		awardType := cForm.AwardType(convert.GetInt32(item["award_type"]))
		awardValue := convert.GetString(item["award_value"])
		switch awardType {
		case cForm.AwardType_CostAwardPoint:
			awardValue = util.ConvertAmount2Decimal(item["award_value"]).Mul(cForm.COST_AWARD_POINT_STEP).String()
		default:
		}

		result = append(result, &StepByStepAward{
			CreatedAt:   convert.GetString(item["created_at"]),
			UserID:      convert.GetInt64(item["user_id"]),
			UserName:    convert.GetString(item["user_name"]),
			Period:      convert.GetInt64(item["period"]),
			StepNO:      convert.GetUint32(item["step_no"]),
			CellNO:      convert.GetUint32(item["cell_no"]),
			AwardType:   awardType,
			AwardValue:  awardValue,
			AwardName:   convert.GetString(item["award_name"]),
			AwardParams: convert.GetString(item["award_params"]),
			AwardNum:    convert.GetInt(item["award_num"]),
			InnerPrice:  util.ConvertAmount2Decimal(item["inner_price"]),
		})
	}

	return
}

func DetailFormat2Excel(dateTimeRange [2]time.Time, _data []map[string]any) (excelModel *excel.Excel[*StepByStepAward], err error) {
	data, err := DetailFormat(_data)
	if err != nil {
		return nil, err
	}

	reflectMap := map[string]func(source *StepByStepAward) any{
		"时间":     func(source *StepByStepAward) any { return source.CreatedAt },
		"用户ID":   func(source *StepByStepAward) any { return source.UserID },
		"用户昵称":   func(source *StepByStepAward) any { return source.UserName },
		"期数":     func(source *StepByStepAward) any { return source.Period },
		"层号":     func(source *StepByStepAward) any { return source.StepNO },
		"格号":     func(source *StepByStepAward) any { return source.CellNO },
		"奖励类型":   func(source *StepByStepAward) any { return source.AwardType },
		"奖励值":    func(source *StepByStepAward) any { return source.AwardValue },
		"奖励名称":   func(source *StepByStepAward) any { return source.AwardName },
		"奖励补充参数": func(source *StepByStepAward) any { return source.AwardParams },
		"奖励数量":   func(source *StepByStepAward) any { return source.AwardNum },
		"成本价":    func(source *StepByStepAward) any { return source.InnerPrice },
	}

	excelModel = &excel.Excel[*StepByStepAward]{
		FileName:   fmt.Sprintf("step_by_step_detail_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"步步高升奖励日志"},
		SheetNameWithHead: map[string][]string{
			"步步高升奖励日志": {
				"时间", "用户ID", "用户昵称",
				"期数", "层号", "格号",
				"奖励类型", "奖励值", "奖励名称", "奖励补充参数", "奖励数量", "成本价",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*StepByStepAward]{
			"步步高升奖励日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*StepByStepAward]{
			"步步高升奖励日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
