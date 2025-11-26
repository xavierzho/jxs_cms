package service

import (
	"context"

	"data_backend/apps/v2/internal/activity/sign_in/dao"
	"data_backend/apps/v2/internal/activity/sign_in/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SignInSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	dao    *dao.SignInDao
}

func NewSignInSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *SignInSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".SignInSvc")))
	return &SignInSvc{
		ctx:    ctx,
		logger: log,
		dao:    dao.NewSignInDao(center, log),
	}
}

func (svc *SignInSvc) List(params *form.Request) (data []*form.SignIn, total int64, e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, 0, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, total, err := svc.dao.List(dateTimeRange, queryParams, params.Pager)
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	data = form.Format(_data)

	return
}

func (svc *SignInSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.SignIn], e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(dateTimeRange, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("Export Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}
