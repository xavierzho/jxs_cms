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

const (
	CostAwardLogType_Normal          = 0   // 消费返欧气值
	CostAwardLogType_Invite          = 1   // 邀请用户消费返欧气值
	CostAwardLogType_Accept          = 100 // 欧气值兑换
	CostAwardLogType_Turntable       = 101 // 转盘抽奖
	CostAwardLogType_StepByStep      = 102 // 步步高升
	CostAwardLogType_SignIn          = 103 // 签到
	CostAwardLogType_LuckyNum        = 104 // 幸运数
	CostAwardLogType_CostAwardOffset = 105 // 欧气值抵扣
	CostAwardLogType_Recall          = 106 // 好友召回
	CostAwardLogType_RedemptionCode  = 107 // 兑换码
	CostAwardLogType_Task_CostAmount = 201 // 任务 抽赏有送
	CostAwardLogType_Task_PrizeValue = 202 // 任务 SP悬赏
	CostAwardLogType_Task_Week       = 203 // 任务 周任务（统一使用这个）
	CostAwardLogType_Admin           = 999 // 管理员手动修改
)

type ListLogRequest struct {
	app.Pager
	AllLogRequest
}

func (q *ListLogRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()

	return q.AllLogRequest.Parse()
}

type AllLogRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	LogType []uint32 `form:"log_type[]"`
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
			Prefix: "caul.log_type in ?",
			Value:  []any{q.LogType},
		})
	}

	return
}

func (q *AllLogRequest) Valid() error {
	for _, logType := range q.LogType {
		switch logType {
		case CostAwardLogType_Normal:
		case CostAwardLogType_Invite:
		case CostAwardLogType_Accept:
		case CostAwardLogType_Turntable:
		case CostAwardLogType_StepByStep:
		case CostAwardLogType_SignIn:
		case CostAwardLogType_LuckyNum:
		case CostAwardLogType_CostAwardOffset:
		case CostAwardLogType_Recall:
		case CostAwardLogType_RedemptionCode:
		case CostAwardLogType_Task_CostAmount:
		case CostAwardLogType_Task_PrizeValue:
		case CostAwardLogType_Task_Week:
		case CostAwardLogType_Admin:
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
