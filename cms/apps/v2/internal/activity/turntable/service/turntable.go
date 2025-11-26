package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/turntable/dao"
	"data_backend/apps/v2/internal/activity/turntable/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TurntableSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.TurntableDao
}

func NewTurntableSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *TurntableSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".TurntableSvc")))
	return &TurntableSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewTurntableDao(center, log),
	}
}

// 获取转盘抽奖列表
func (svc *TurntableSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.Turntable, e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_summary, _data, err := svc.dao.ListAndSummary(
		[]string{"count(0) as total,sum(i.inner_price) as total_amount"},
		dateRange, queryParams, params.Pager,
	)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.Format(svc.ctx, dateRange, _summary, _data)

	return

}

// 导出转盘抽奖列表
func (svc *TurntableSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Turntable], e *errcode.Error) {
	dateRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateRange, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(svc.ctx, dateRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
