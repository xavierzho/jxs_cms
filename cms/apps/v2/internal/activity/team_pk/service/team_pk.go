package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/team_pk/dao"
	"data_backend/apps/v2/internal/activity/team_pk/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TeamPKSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.TeamPKDao
}

func NewTeamPKSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *TeamPKSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".TeamPKSvc")))
	return &TeamPKSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewTeamPKDao(center, log),
	}
}

// list
func (svc *TeamPKSvc) List(params *form.ListRequest) (data []*form.Team, total int64, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, 0, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, total, err := svc.dao.List(dateRange, params.Pager)
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	data = form.Format(_data)

	return data, total, nil
}

// export
func (svc *TeamPKSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Team], e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateRange)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(dateRange, _data)
	if err != nil {
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return data, nil
}
