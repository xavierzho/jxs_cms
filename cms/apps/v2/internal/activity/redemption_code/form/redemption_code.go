package form

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/activity/redemption_code/dao"
	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type LogRequest struct {
	app.Pager
	AllRequest
}

func (q *LogRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	q.Pager.Parse()
	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	Name string `form:"name"`
	Code string `form:"code"`
}

func (q *AllRequest) Parse() (dateTimeRange [2]time.Time, paramsGroup dao.AllRequestParamsGroup, err error) {
	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	// UsersParams
	if paramsGroup.UsersParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	{
		if q.Name != "" {
			paramsGroup.NameParams = append(paramsGroup.NameParams, database.QueryWhere{
				Prefix: "c.name like ?",
				Value:  []any{"%" + q.Name + "%"},
			})
		}
	}
	{
		if q.Code != "" {
			paramsGroup.CodeParams = append(paramsGroup.CodeParams, database.QueryWhere{
				Prefix: "c.code = ?",
				Value:  []any{q.Code},
			})
		}
	}

	return
}

type RedemptionCodeSvc struct {
	LogID                     string          `json:"log_id"`
	RedemptionCodeID          string          `json:"redemption_code_id"`
	DateTime                  string          `json:"date_time"`
	UserID                    int64           `json:"user_id"`
	UserName                  string          `json:"user_name"`
	Code                      string          `json:"code"`
	Name                      string          `json:"name"`
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

type RedemptionCodeAwardDetailSvc struct {
	DateTime                  string          `json:"date_time"`
	UserID                    int64           `json:"user_id"`
	UserName                  string          `json:"user_name"`
	Code                      string          `json:"code"`
	Name                      string          `json:"name"`
	AwardType                 string          `json:"award_type"`
	AwardName                 string          `json:"award_name"`
	RewardValueItem           decimal.Decimal `json:"reward_value_item"`
	RewardValueCostAwardPoint decimal.Decimal `json:"reward_value_cost_award_point"`
	AwardNum                  int64           `json:"award_num"`
}

func Format(ctx context.Context, _summary map[string]any, data []*dao.RedemptionCodeLog) (summary map[string]any, result []*RedemptionCodeSvc, err error) {
	summary = _summary
	if summary != nil {
		summary["reward_value_item"] = util.ConvertAmount2Decimal(summary["reward_value_item"])
		summary["reward_value_cost_award_point"] = util.ConvertAmount2Decimal(summary["reward_value_cost_award_point"]).Mul(cForm.COST_AWARD_POINT_STEP)
	}

	for _, item := range data {
		result = append(result, &RedemptionCodeSvc{
			LogID:                     item.LogID,
			RedemptionCodeID:          item.RedemptionCodeID,
			DateTime:                  item.DateTime,
			UserID:                    item.UserID,
			UserName:                  item.UserName,
			Name:                      item.Name,
			Code:                      item.Code,
			RewardValueItem:           util.ConvertAmount2Decimal(item.RewardValueItem),
			RewardValueCostAwardPoint: util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP),
		})
	}

	return
}

type DetailRequest struct {
	LogID string `form:"log_id" binding:"required"`
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

func FormatLog2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.RedemptionCodeLog) (excelModel *excel.Excel[*RedemptionCodeSvc], err error) {
	_, data, err := Format(ctx, nil, _data)
	if err != nil {
		return nil, err
	}

	reflectMap := map[string]func(source *RedemptionCodeSvc) any{
		"时间":        func(source *RedemptionCodeSvc) any { return source.DateTime },
		"用户ID":      func(source *RedemptionCodeSvc) any { return source.UserID },
		"用户昵称":      func(source *RedemptionCodeSvc) any { return source.UserName },
		"名称":        func(source *RedemptionCodeSvc) any { return source.Name },
		"兑换码":       func(source *RedemptionCodeSvc) any { return source.Code },
		"奖励价值(物品)":  func(source *RedemptionCodeSvc) any { return source.RewardValueItem },
		"奖励价值(欧气值)": func(source *RedemptionCodeSvc) any { return source.RewardValueCostAwardPoint },
	}
	excelModel = &excel.Excel[*RedemptionCodeSvc]{
		FileName:   fmt.Sprintf("activity_redemption_code_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"兑换记录"},
		SheetNameWithHead: map[string][]string{
			"兑换记录": {
				"时间", "用户ID", "用户昵称", "名称", "兑换码", "奖励价值(物品)", "奖励价值(欧气值)",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*RedemptionCodeSvc]{
			"兑换记录": data,
		},
		ReflectMap: map[string]excel.RowReflect[*RedemptionCodeSvc]{
			"兑换记录": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

func FormatAwardDetail2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.RedemptionCodeLogAwardDetail) (excelModel *excel.Excel[*RedemptionCodeAwardDetailSvc], err error) {
	data, err := FormatLogAwardDetail(ctx, _data)
	if err != nil {
		return
	}

	reflectMap := map[string]func(source *RedemptionCodeAwardDetailSvc) any{
		"时间":        func(source *RedemptionCodeAwardDetailSvc) any { return source.DateTime },
		"用户ID":      func(source *RedemptionCodeAwardDetailSvc) any { return source.UserID },
		"用户昵称":      func(source *RedemptionCodeAwardDetailSvc) any { return source.UserName },
		"名称":        func(source *RedemptionCodeAwardDetailSvc) any { return source.Name },
		"兑换码":       func(source *RedemptionCodeAwardDetailSvc) any { return source.Code },
		"奖励类型":      func(source *RedemptionCodeAwardDetailSvc) any { return source.AwardType },
		"奖励名称":      func(source *RedemptionCodeAwardDetailSvc) any { return source.AwardName },
		"奖励价值(物品)":  func(source *RedemptionCodeAwardDetailSvc) any { return source.RewardValueItem },
		"奖励价值(欧气值)": func(source *RedemptionCodeAwardDetailSvc) any { return source.RewardValueCostAwardPoint },
		"奖励数量":      func(source *RedemptionCodeAwardDetailSvc) any { return source.AwardNum },
	}

	excelModel = &excel.Excel[*RedemptionCodeAwardDetailSvc]{
		FileName:   fmt.Sprintf("inquire_task_award_detail_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"兑换码记录详情日志"},
		SheetNameWithHead: map[string][]string{
			"兑换码记录详情日志": {
				"时间",
				"用户ID", "用户昵称",
				"名称", "兑换码",
				"奖励类型", "奖励名称", "奖励价值(物品)", "奖励价值(欧气值)", "奖励数量",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*RedemptionCodeAwardDetailSvc]{
			"兑换码记录详情日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*RedemptionCodeAwardDetailSvc]{
			"兑换码记录详情日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

func FormatLogAwardDetail(ctx context.Context, data []*dao.RedemptionCodeLogAwardDetail) (result []*RedemptionCodeAwardDetailSvc, err error) {
	for _, item := range data {
		if item.AwardNum == 0 {
			item.AwardNum = 1
		}
		result = append(result, &RedemptionCodeAwardDetailSvc{
			DateTime:                  item.DateTime,
			UserID:                    item.UserID,
			UserName:                  item.UserName,
			Code:                      item.Code,
			Name:                      item.Name,
			AwardType:                 global.I18n.T(ctx, "common.award_type", fmt.Sprintf("%d", item.AwardType)),
			AwardName:                 item.AwardName,
			AwardNum:                  item.AwardNum,
			RewardValueItem:           util.ConvertAmount2Decimal(item.RewardValueItem),
			RewardValueCostAwardPoint: util.ConvertAmount2Decimal(item.RewardValueCostAwardPoint).Mul(cForm.COST_AWARD_POINT_STEP),
		})
	}

	return
}
