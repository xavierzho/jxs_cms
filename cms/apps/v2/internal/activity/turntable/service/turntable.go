package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/turntable/dao"
	"data_backend/apps/v2/internal/activity/turntable/form"
	"data_backend/pkg/database"
	"data_backend/pkg/encrypt/md5"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TurntableSvc struct {
	logger *logger.Logger
	dao    *dao.TurntableDao
}

func NewTurntableSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *TurntableSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".TurntableSvc")))
	return &TurntableSvc{
		logger: log,
		dao:    dao.NewTurntableDao(center, log),
	}
}

// 获取查询条件
func (svc *TurntableSvc) userCondition(params *form.AllRequest) (condition database.QueryWhereGroup) {
	if params.UserID != 0 {
		condition = append(condition, database.QueryWhere{Prefix: "u.id = ?", Value: []any{params.UserID}})
	}
	if params.UserName != "" {
		condition = append(condition, database.QueryWhere{Prefix: "u.nickname = ?", Value: []any{params.UserName}})
	}
	if params.Tel != "" {
		condition = append(condition, database.QueryWhere{Prefix: "u.phone_num_md5 = ?", Value: []any{md5.EncodeMD5(params.Tel)}})
	}
	if params.Type != 0 {
		if params.Type == form.PrizeType_Coin {
			condition = append(condition, database.QueryWhere{Prefix: "apwac.type = ?", Value: []any{0}})
		} else {
			condition = append(condition, database.QueryWhere{Prefix: "apwac.type = ?", Value: []any{params.Type}})
		}
	}
	if params.Name != "" {
		condition = append(condition, database.QueryWhere{Prefix: "apwc.name = ?", Value: []any{params.Name}})
	}
	return condition

}

// 获取转盘抽奖列表
func (svc *TurntableSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.Turntable, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_summary, _data, err := svc.dao.ListAndSummary([]string{"count(0) as total,sum(i.inner_price) as total_amount"}, dateRange, svc.userCondition(&params.AllRequest), params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.Format(dateRange, _summary, _data)

	return

}

// 导出转盘抽奖列表
func (svc *TurntableSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Turntable], e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateRange, svc.userCondition(params))
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(dateRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
