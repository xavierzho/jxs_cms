package form

import (
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type ListLogRequest struct {
	app.Pager
	AllLogRequest
	BalanceType uint `form:"balance_type" json:"balance_type"`
}

func (q *ListLogRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.AllLogRequest.Parse()
}

type AllLogRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	LogType     []int32 `form:"log_type[]"`
	BalanceType uint    `form:"balance_type" json:"balance_type"`
}

func (q *AllLogRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	if len(q.LogType) != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "bl.source_type in ?",
			Value:  []any{q.LogType},
		})
	}

	return
}

func (q *AllLogRequest) Valid() error {
	for _, logType := range q.LogType {
		switch logType {
		case 100002, 100009, 100008, 100010, 200000:
		default:
			return fmt.Errorf("not expected LogType: %d", q.LogType)
		}
	}

	return nil
}

type CostAwardLog struct {
	CreatedAt   string          `json:"created_at"`
	UserID      int64           `json:"user_id"`
	UserName    string          `json:"user_name"`
	LogTypeStr  string          `json:"log_type_str"`
	UpdatePoint decimal.Decimal `json:"update_point"`
	BeforePoint decimal.Decimal `json:"before_point"`
	AfterPoint  decimal.Decimal `json:"after_point"`
}

func FormatLog(ctx *gin.Context, _summary map[string]any, data []map[string]any) (summary map[string]any, result []*CostAwardLog) {
	summary = _summary
	if summary != nil {
		summary["update_amount"] = util.ConvertAmount2Decimal(summary["update_amount"])
	}

	for _, item := range data {
		result = append(result, &CostAwardLog{
			CreatedAt:   convert.GetString(item["created_at"]),
			UserID:      convert.GetInt64(item["user_id"]),
			UserName:    convert.GetString(item["user_name"]),
			LogTypeStr:  global.I18n.T(ctx, "source_type", convert.GetString(item["source_type"])),
			UpdatePoint: util.ConvertAmount2Decimal(item["update_amount"]),
			BeforePoint: util.ConvertAmount2Decimal(item["before_balance"]),
			AfterPoint:  util.ConvertAmount2Decimal(item["after_balance"]),
		})
	}

	return
}

func FormatLog2Excel(ctx *gin.Context, dateTimeRange [2]time.Time, _data []map[string]any) (excelModel *excel.Excel[*CostAwardLog], err error) {
	_, data := FormatLog(ctx, nil, _data)

	reflectMap := map[string]func(source *CostAwardLog) any{
		"时间":     func(source *CostAwardLog) any { return source.CreatedAt },
		"用户ID":   func(source *CostAwardLog) any { return source.UserID },
		"用户昵称":   func(source *CostAwardLog) any { return source.UserName },
		"变动类型":   func(source *CostAwardLog) any { return source.LogTypeStr },
		"欧气值变动前": func(source *CostAwardLog) any { return source.BeforePoint },
		"欧气值变动后": func(source *CostAwardLog) any { return source.AfterPoint },
		"欧气值变动":  func(source *CostAwardLog) any { return source.UpdatePoint },
	}

	excelModel = &excel.Excel[*CostAwardLog]{
		FileName:   fmt.Sprintf("cost_award_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"欧气值流水日志"},
		SheetNameWithHead: map[string][]string{
			"欧气值流水日志": {
				"时间", "用户ID", "用户昵称",
				"变动类型", "欧气值变动前", "欧气值变动后", "欧气值变动",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*CostAwardLog]{
			"欧气值流水日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*CostAwardLog]{
			"欧气值流水日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
