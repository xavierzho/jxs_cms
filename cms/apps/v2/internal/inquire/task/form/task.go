package form

import (
	"context"
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/task/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

// 范围兼容 只输入一个; 前端设置值为 2**32 占位
type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	q.Pager.Parse()
	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	TaskTypeList []int    `form:"task_type_list[]"` // dbBranch && GachaParams
	TaskKeyList  []string `form:"task_key_list[]"`
}

func (q *AllRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
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

	{
		if len(q.TaskTypeList) != 0 {
			paramsGroup.TaskTypeParams = append(paramsGroup.TaskTypeParams, database.QueryWhere{
				Prefix: "t.`type` in ?",
				Value:  []any{q.TaskTypeList},
			})
		}
	}
	{
		if len(q.TaskKeyList) != 0 {
			paramsGroup.TaskKeyParams = append(paramsGroup.TaskKeyParams, database.QueryWhere{
				Prefix: "t.`key` in ?",
				Value:  []any{q.TaskKeyList},
			})
		}
	}

	return
}

func (q *AllRequest) Valid() (err error) {
	if len(q.TaskTypeList) != 0 {
		for _, i := range q.TaskTypeList {
			switch i {
			case 1, 2, 3, 4, 5:
			default:
				return fmt.Errorf("not expected task type: %+v", q.TaskTypeList)
			}
		}
	}
	if len(q.TaskKeyList) != 0 {
		for _, i := range q.TaskKeyList {
			switch i {
			case "CostAmount", "PrizeValue", "Week1", "Week2", "Week3", "Week4", "Week5", "Week6", "Week7", "Weekend":
			default:
				return fmt.Errorf("not expected task key: %+v", q.TaskKeyList)
			}
		}
	}
	return nil
}

type TaskSvc struct {
	TaskID                    string          `json:"task_id"`
	DateTime                  string          `json:"date_time"`
	UserID                    int64           `json:"user_id"`
	UserName                  string          `json:"user_name"`
	TaskTypeStr               string          `json:"task_type_str"`
	TaskKeyStr                string          `json:"task_key_str"`
	TaskTypeName              string          `json:"task_name"`
	RewardValueItem           decimal.Decimal `json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `json:"reward_value_cost_award_point"`
}

type AwardDetailSvc struct {
	AwardType                 string          `json:"award_type"`
	AwardName                 string          `json:"award_name"`
	RewardValueItem           decimal.Decimal `json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `json:"reward_value_cost_award_point"`
	AwardNum                  int64           `json:"award_num"`
}

type TaskAwardDetailSvc struct {
	DateTime                  string          `json:"date_time"`
	UserID                    int64           `json:"user_id"`
	UserName                  string          `json:"user_name"`
	TaskTypeStr               string          `json:"task_type_str"`
	TaskKeyStr                string          `json:"task_key_str"`
	TaskTypeName              string          `json:"task_name"`
	AwardType                 string          `json:"award_type"`
	AwardName                 string          `json:"award_name"`
	RewardValueItem           decimal.Decimal `json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `json:"reward_value_cost_award_point"`
	AwardNum                  int64           `json:"award_num"`
}

func Format(ctx context.Context, _summary map[string]any, data []*dao.TaskList) (summary map[string]any, result []*TaskSvc, err error) {
	summary = _summary
	if summary != nil {
		summary["reward_value_item"] = util.ConvertAmount2Decimal(summary["reward_value_item"])
		summary["reward_value_cost_award_point"] = util.ConvertAmount2Decimal(summary["reward_value_cost_award_point"]).Mul(cForm.COST_AWARD_POINT_STEP)
	}

	for _, item := range data {
		result = append(result, &TaskSvc{
			TaskID:                    item.TaskID,
			DateTime:                  item.DateTime,
			UserID:                    item.UserID,
			UserName:                  item.UserName,
			TaskTypeStr:               global.I18n.T(ctx, "task.type", fmt.Sprintf("%d", item.TaskType)),
			TaskKeyStr:                global.I18n.T(ctx, "task.key", fmt.Sprintf("%s", item.TaskKey)),
			TaskTypeName:              item.TaskName,
			RewardValueItem:           util.ConvertAmount2Decimal(item.RewardValueItem),
			RewardValueCostAwardPoint: util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP),
		})
	}

	return
}

type DetailRequest struct {
	TaskID string `form:"task_id" binding:"required"`
}

func FormatAwardDetail(ctx context.Context, data []*dao.AwardDetail) (result []*AwardDetailSvc, err error) {

	for _, item := range data {
		if item.AwardNum == 0 {
			item.AwardNum = 1
		}
		awardDetail := AwardDetailSvc{
			AwardType: global.I18n.T(ctx, "common.award_type", fmt.Sprintf("%d", item.AwardType)),
			AwardName: item.AwardName,
			AwardNum:  item.AwardNum,
		}
		if item.AwardType == cForm.AwardType_CostAwardPoint {
			awardDetail.RewardValueCostAwardPoint = util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP)
		} else if item.AwardType == cForm.AwardType_Item {
			awardDetail.RewardValueItem = util.ConvertAmount2Decimal(item.RewardValueItem)
		}
		result = append(result, &awardDetail)
	}

	return
}

func FormatList(ctx context.Context, _summary map[string]any, data []*dao.TaskList) (summary map[string]any, result []*TaskSvc, err error) {
	if _summary != nil {
		summary = _summary
		summary["reward_value_item"] = util.ConvertAmount2Decimal(_summary["reward_value_item"])
	}

	for _, item := range data {
		result = append(result, &TaskSvc{
			DateTime:                  item.DateTime,
			UserID:                    item.UserID,
			UserName:                  item.UserName,
			TaskTypeStr:               global.I18n.T(ctx, "task.type", fmt.Sprintf("%d", item.TaskType)),
			TaskKeyStr:                global.I18n.T(ctx, "task.key", fmt.Sprintf("%s", item.TaskKey)),
			TaskTypeName:              item.TaskName,
			RewardValueItem:           util.ConvertAmount2Decimal(item.RewardValueItem),
			RewardValueCostAwardPoint: util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP),
		})
	}

	return
}

func FormatList2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.TaskList) (excelModel *excel.Excel[*TaskSvc], err error) {
	_, data, err := FormatList(ctx, nil, _data)
	if err != nil {
		return nil, err
	}

	reflectMap := map[string]func(source *TaskSvc) any{
		"时间":        func(source *TaskSvc) any { return source.DateTime },
		"用户ID":      func(source *TaskSvc) any { return source.UserID },
		"用户昵称":      func(source *TaskSvc) any { return source.UserName },
		"类型":        func(source *TaskSvc) any { return source.TaskTypeStr },
		"KEY":       func(source *TaskSvc) any { return source.TaskKeyStr },
		"名称":        func(source *TaskSvc) any { return source.TaskTypeName },
		"奖励价值(物品)":  func(source *TaskSvc) any { return source.RewardValueItem },
		"奖励价值(欧气值)": func(source *TaskSvc) any { return source.RewardValueCostAwardPoint },
	}
	excelModel = &excel.Excel[*TaskSvc]{
		FileName:   fmt.Sprintf("inquire_task_list_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"任务记录"},
		SheetNameWithHead: map[string][]string{
			"任务记录": {
				"时间", "用户ID", "用户昵称", "类型", "KEY", "名称", "奖励价值(物品)",
				"奖励价值(欧气值)",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*TaskSvc]{
			"任务记录": data,
		},
		ReflectMap: map[string]excel.RowReflect[*TaskSvc]{
			"任务记录": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

func FormatTaskAwardDetail2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.TaskListAwardDetail) (excelModel *excel.Excel[*TaskAwardDetailSvc], err error) {
	data, err := FormatTaskAwardDetail(ctx, _data)
	if err != nil {
		return
	}

	reflectMap := map[string]func(source *TaskAwardDetailSvc) any{
		"时间":        func(source *TaskAwardDetailSvc) any { return source.DateTime },
		"用户ID":      func(source *TaskAwardDetailSvc) any { return source.UserID },
		"用户昵称":      func(source *TaskAwardDetailSvc) any { return source.UserName },
		"任务类型":      func(source *TaskAwardDetailSvc) any { return source.TaskKeyStr },
		"任务KEY":     func(source *TaskAwardDetailSvc) any { return source.TaskKeyStr },
		"任务名称":      func(source *TaskAwardDetailSvc) any { return source.TaskTypeName },
		"奖励类型":      func(source *TaskAwardDetailSvc) any { return source.AwardType },
		"奖励名称":      func(source *TaskAwardDetailSvc) any { return source.AwardName },
		"奖励价值(物品)":  func(source *TaskAwardDetailSvc) any { return source.RewardValueItem },
		"奖励价值(欧气值)": func(source *TaskAwardDetailSvc) any { return source.RewardValueCostAwardPoint },
		"奖励数量":      func(source *TaskAwardDetailSvc) any { return source.AwardNum },
	}

	excelModel = &excel.Excel[*TaskAwardDetailSvc]{
		FileName:   fmt.Sprintf("inquire_task_award_detail_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"任务记录详情日志"},
		SheetNameWithHead: map[string][]string{
			"任务记录详情日志": {
				"时间",
				"用户ID", "用户昵称",
				"任务类型", "任务KEY", "任务名称",
				"奖励类型", "奖励名称", "奖励价值(物品)", "奖励价值(欧气值)", "奖励数量",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*TaskAwardDetailSvc]{
			"任务记录详情日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*TaskAwardDetailSvc]{
			"任务记录详情日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

func FormatTaskAwardDetail(ctx context.Context, data []*dao.TaskListAwardDetail) (result []*TaskAwardDetailSvc, err error) {
	for _, item := range data {
		if item.AwardNum == 0 {
			item.AwardNum = 1
		}
		result = append(result, &TaskAwardDetailSvc{
			DateTime:                  item.DateTime,
			UserID:                    item.UserID,
			UserName:                  item.UserName,
			TaskTypeStr:               global.I18n.T(ctx, "task.type", fmt.Sprintf("%d", item.TaskType)),
			TaskKeyStr:                global.I18n.T(ctx, "task.key", item.TaskKey),
			TaskTypeName:              item.TaskName,
			AwardType:                 global.I18n.T(ctx, "common.award_type", fmt.Sprintf("%d", item.AwardType)),
			AwardName:                 item.AwardName,
			AwardNum:                  item.AwardNum,
			RewardValueItem:           util.ConvertAmount2Decimal(item.RewardValueItem),
			RewardValueCostAwardPoint: util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP),
		})
	}

	return
}
