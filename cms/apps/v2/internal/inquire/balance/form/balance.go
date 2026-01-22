package form

import (
	"context"
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/balance/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type ListRequest struct {
	*app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()
	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	DateTimeType      dao.DateTimeType `form:"date_time_type"`
	SourceType        []int            `form:"source_type[]"`
	ItemName          string           `form:"item_name"`
	ChannelType       []int            `form:"channel_type[]"` // 1 支付宝 2 微信
	UpdateAmountRange *[2]int64        `form:"update_amount_range[]"`
	PaySourceType     []int            `form:"pay_source_type[]"` // 充值目标
	BalanceType       []int            `form:"balance_type[]"`
}

func (q *AllRequest) Parse() (dateTimeRange [2]time.Time, queryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return dateTimeRange, nil, err
	} else {
		switch q.DateTimeType {
		case dao.DateTimeType_Created:
			queryParams = append(queryParams, database.QueryWhere{
				Prefix: "bl.created_at between ? and ?",
				Value:  []any{dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second - time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)},
			})
		case dao.DateTimeType_Finish:
			queryParams = append(queryParams, database.QueryWhere{
				Prefix: "bl.finish_at between ? and ?",
				Value:  []any{dateTimeRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateTimeRange[1].Add(time.Second - time.Millisecond).Format(pkg.DATE_TIME_MIL_FORMAT)},
			})
		default:
			return dateTimeRange, nil, fmt.Errorf("not expected dateTimeType: %s", q.DateTimeType)
		}
	}

	if _queryParams, err := q.UserInfoRequest.Parse(); err != nil {
		return dateTimeRange, nil, err
	} else {
		queryParams = append(queryParams, _queryParams...)
	}

	if len(q.SourceType) != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "bl.source_type in ?",
			Value:  []any{q.SourceType},
		})
	}

	if q.ItemName != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "(gm.name = ? or cac.name = ?) or (gm.name is null and cac.name is null)", // 兼容还选了除 gacha 外的其他类型的情况
			Value:  []any{q.ItemName, q.ItemName},
		})
	}

	if len(q.ChannelType) != 0 {
		channelTypeList := []string{}
		for _, channelType := range q.ChannelType {
			switch channelType {
			case 1:
				channelTypeList = append(channelTypeList, "alipay")
			case 2:
				channelTypeList = append(channelTypeList, "wechatjs")
				channelTypeList = append(channelTypeList, "wechatapp")
			}
		}

		if len(channelTypeList) != 0 {
			queryParams = append(queryParams, database.QueryWhere{
				Prefix: "ppo.platform_id in ? or ppo.platform_id is null",
				Value:  []any{channelTypeList},
			})
		}
	}

	if len(q.PaySourceType) != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "ppo.pay_source_type in ? or ppo.platform_id is null",
			Value:  []any{q.PaySourceType},
		})
	}

	if q.UpdateAmountRange != nil {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "bl.update_amount between ? and ?",
			Value:  []any{util.ReconvertAmount2Decimal(q.UpdateAmountRange[0]).IntPart(), util.ReconvertAmount2Decimal(q.UpdateAmountRange[1]).IntPart()},
		})
	}

	if len(q.BalanceType) != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "bl.type in ?",
			Value:  []any{q.BalanceType},
		})
	}

	return
}

func (q *AllRequest) Valid() (err error) {
	switch q.DateTimeType {
	case dao.DateTimeType_Created, dao.DateTimeType_Finish:
	default:
		return fmt.Errorf("not expected dateTimeType: %s", q.DateTimeType)
	}

	for _, sourceType := range q.SourceType {
		switch sourceType {
		case 1, 2, 3, 12, 13:
		case 15: // 兑换
		case 101, 102, 103, 104, 105, 106:
		case 201, 202, 203, 204: // 吉市
		case 301, 302, 303, 304:
		case 400:
		case 601:
		case 100004: // 物品置换
		case 100005:
		case 100009: // 积分兑换
		case 100013: // 吉祥值抵扣
		case 999999:
		default:
			return fmt.Errorf("not expected sourceType: %d", q.SourceType)
		}
	}

	for _, paySourceType := range q.PaySourceType {
		switch paySourceType {
		case 100:
		case 201, 202:
		case 301:
		case 601:
		default:
			return fmt.Errorf("not expected paySourceType: %d", q.PaySourceType)
		}
	}

	for _, channelType := range q.ChannelType {
		switch channelType {
		case 1, 2:
		default:
			return fmt.Errorf("not expected channelType: %d", q.ChannelType)
		}
	}

	if q.UpdateAmountRange != nil && q.UpdateAmountRange[1] < q.UpdateAmountRange[0] {
		return fmt.Errorf("invalid UpdateAmountRange: %v", q.UpdateAmountRange)
	}

	return nil
}

type Balance struct {
	ID               int64                                   `json:"id"`
	CreatedAt        string                                  `json:"created_at"`
	FinishAt         string                                  `json:"finish_at"`
	UserID           int64                                   `json:"user_id"`
	UserName         string                                  `json:"user_name"`
	SourceTypeStr    string                                  `json:"source_type_str"`
	ItemName         string                                  `json:"item_name"`
	PlatformOrderId  string                                  `json:"platform_order_id"`
	PaySourceTypeStr string                                  `json:"pay_source_type_str"`
	BeforeBalance    decimal.Decimal                         `json:"before_balance"`
	AfterBalance     decimal.Decimal                         `json:"after_balance"`
	UpdateAmount     decimal.Decimal                         `json:"update_amount"`
	Comment          datatypes.JSONSlice[dao.BalanceComment] `json:"comment"`
	BalanceTypeName  string                                  `json:"balance_type_name"`
}

func Format(ctx context.Context, _summary map[string]any, data []*dao.Balance) (summary map[string]any, result []*Balance) {
	if _summary != nil {
		summary = map[string]any{
			"cnt":           _summary["cnt"],
			"user_cnt":      _summary["user_cnt"],
			"update_amount": util.ConvertAmount2Decimal(_summary["update_amount"]),
		}
	}

	for _, item := range data {
		channelType := ""
		switch item.ChannelType {
		case "":
		case "alipay":
			channelType = global.I18n.T(ctx, "pay.channelType", "1")
		case "wechatjs", "wechatapp":
			channelType = global.I18n.T(ctx, "pay.channelType", "2")
		}

		platformOrderId := ""
		switch item.SourceType {
		case 1, 12:
			platformOrderId = item.PlatformOrderIdPay
		case 2:
			platformOrderId = item.PlatformOrderIdDraw
		default:
		}

		// 特殊处理 ItemName
		if item.SourceType == 15 {
			item.CostAwardName = "吉祥币兑换"
		} else if item.SourceType >= 200 && item.SourceType <= 299 {
			item.CostAwardName = "吉市订单"
		} else if item.SourceType == 100013 {
			item.CostAwardName = "吉祥值抵扣"
		} else if item.SourceType == 100009 {
			if item.CostAwardName == "" {
				item.CostAwardName = "积分兑换"
			} else {
				item.CostAwardName = "积分兑换-" + item.CostAwardName
			}
		} else if item.SourceType == 100004 {
			if item.ItemExchangeName == "" {
				item.ItemExchangeName = "物品置换"
			} else {
				item.ItemExchangeName = "物品置换-" + item.ItemExchangeName
			}
		}

		var paySourceTypeStr string
		if item.PaySourceType != 0 {
			paySourceTypeStr = global.I18n.T(ctx, "source_type", convert.GetString(item.PaySourceType))
		}
		result = append(result, &Balance{
			ID:               item.ID,
			CreatedAt:        item.CreatedAt,
			FinishAt:         item.FinishAt,
			UserID:           item.UserID,
			UserName:         item.UserName,
			SourceTypeStr:    global.I18n.T(ctx, "source_type", convert.GetString(item.SourceType)),
			ItemName:         item.GachaName + item.CostAwardName + item.ItemExchangeName + channelType, // 三者仅一者非空
			PlatformOrderId:  platformOrderId,
			PaySourceTypeStr: paySourceTypeStr,
			BeforeBalance:    util.ConvertAmount2Decimal(item.BeforeBalance),
			AfterBalance:     util.ConvertAmount2Decimal(item.AfterBalance),
			UpdateAmount:     util.ConvertAmount2Decimal(item.UpdateAmount),
			Comment:          item.Comment,
			BalanceTypeName:  global.I18n.T(ctx, "balance_type", fmt.Sprintf("%d", item.BalanceType)),
		})
	}

	return
}

func Format2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.Balance) (excelModel *excel.Excel[*Balance], err error) {
	_, data := Format(ctx, nil, _data)

	reflectMap := map[string]func(*Balance) any{
		"ID":      func(source *Balance) any { return source.ID },
		"创建时间":    func(source *Balance) any { return source.CreatedAt },
		"完成时间":    func(source *Balance) any { return source.FinishAt },
		"用户id":    func(source *Balance) any { return source.UserID },
		"用户昵称":    func(source *Balance) any { return source.UserName },
		"类型":      func(source *Balance) any { return source.SourceTypeStr },
		"项目名称":    func(source *Balance) any { return source.ItemName },
		"第三方订单id": func(source *Balance) any { return source.PlatformOrderId },
		"充值目标":    func(source *Balance) any { return source.PaySourceTypeStr },
		"余额变动前":   func(source *Balance) any { return source.BeforeBalance },
		"余额变动后":   func(source *Balance) any { return source.AfterBalance },
		"余额变动":    func(source *Balance) any { return source.UpdateAmount },
		"备注":      func(source *Balance) any { return source.Comment },
	}

	excelModel = &excel.Excel[*Balance]{
		FileName:   fmt.Sprintf("user_balance_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"用户流水日志"},
		SheetNameWithHead: map[string][]string{
			"用户流水日志": {
				"ID", "创建时间", "完成时间", "用户id", "用户昵称", "类型", "项目名称",
				"第三方订单id", "充值目标",
				"余额变动前", "余额变动后", "余额变动",
				"备注",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*Balance]{
			"用户流水日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*Balance]{
			"用户流水日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

type AddCommentRequest struct {
	ID      int64  `form:"id" binding:"required"`
	Comment string `form:"comment" binding:"required"`
}

type DeleteCommentRequest struct {
	ID        int64 `form:"id" binding:"required"`
	CommentID int64 `form:"comment_id" binding:"required"`
}
