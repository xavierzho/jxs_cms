package service

import (
	"context"

	"data_backend/apps/v2/internal/inquire/item/dao"
	"data_backend/apps/v2/internal/inquire/item/form"
	"data_backend/internal/global"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.ItemDao
}

func NewItemSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *ItemSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".ItemSvc")))
	return &ItemSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewItemDao(center, log),
	}
}

func (svc *ItemSvc) OptionsLogType() []map[string]string {
	return []map[string]string{
		{"value": "101", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "101")},
		{"value": "102", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "102")},
		{"value": "103", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "103")},
		{"value": "104", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "104")},
		{"value": "200", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "200")},
		{"value": "300", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "300")},
		{"value": "100002", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100002")},
		{"value": "100003", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100003")},
		{"value": "100004", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100004")},
		{"value": "100005", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100005")}, //转盘抽奖
		{"value": "999999", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "999999")},
	}
}

func (svc *ItemSvc) GetLog(params *form.LogRequest) (summary map[string]any, data []*form.ItemLog, e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	summary, _data, err := svc.dao.GetLog(dateTimeRange, params.LogTypeList, paramsGroup, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data, err = form.FormatLog(svc.ctx.Request.Context(), summary, _data)
	if err != nil {
		return nil, nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *ItemSvc) ExportLog(params *form.LogAllRequest) (data *excel.Excel[*form.ItemLog], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_, _data, err := svc.dao.GetLog(dateTimeRange, params.LogTypeList, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatLog2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportLog FormatLog2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *ItemSvc) GetDetail(params *form.DetailRequest) (data any, e *errcode.Error) {
	queryParams, betQueryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	if params.LogType == 200 {
		dataCreator, dataOfferer, err := svc.dao.GetDetailMarket(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return []any{form.FormatMarketItemDetail(dataCreator), form.FormatMarketItemDetail(dataOfferer)}, nil
	} else if params.LogType == 300 {
		data, err := svc.dao.GetDetailOrder(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatItem(data), nil
	} else if params.LogType == 100002 { // 欧气值 存在多份 在前端 进行乘法
		data, err := svc.dao.GetDetailActivityCostAwardConfig(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatItem(data), nil
	} else if params.LogType == 100003 {
		data, err := svc.dao.GetDetailActivityCostRankConfig(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatItem(data), nil
	} else if params.LogType == 100004 {
		dataCreator, dataOfferer, err := svc.dao.GetDetailActivityItemExchange(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return []any{form.FormatMarketItemDetail(dataCreator), form.FormatMarketItemDetail(dataOfferer)}, nil
	} else if params.LogType == 100005 { //转盘抽奖
		data, err := svc.dao.GetDetailActivityPrizeWheel(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatItem(data), nil

	} else if params.LogType == 999999 {
		data, err := svc.dao.GetDetailAdmin(queryParams)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatItem(data), nil
	} else {
		data, err := svc.dao.GetDetailBet(append(queryParams, betQueryParams...))
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}
		return form.FormatBetItemDetail(data), nil
	}
}

// TODO 未使用 update_amount 条件
func (svc *ItemSvc) ExportDetail(params *form.DetailAllRequest) (data *excel.Excel[*form.ItemDetail], e *errcode.Error) {
	dateTimeRange, paramsGroup, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}
	_data, err := svc.dao.GetDetail(params.LogTypeList, paramsGroup, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.FormatItemDetail2Excel(svc.ctx.Request.Context(), dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("ExportDetail FormatItemDetail2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
