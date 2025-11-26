package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	"data_backend/internal/form"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OperationLogSvc struct {
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
	dao    *dao.OperationLogDao
}

func NewOperationLogSvc(ctx *gin.Context, engine *gorm.DB, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *OperationLogSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".OperationLogSvc")))
	alarm := newAlarm(log)
	return &OperationLogSvc{
		ctx:    ctx,
		logger: log,
		alarm:  alarm,
		dao:    dao.NewOperationLogDao(engine, log),
	}
}

func (svc *OperationLogSvc) initModel(data *dao.OperationLog) (form url.Values, err error) {
	// form = make(url.Values)
	// 屏蔽敏感信息
	if svc.ctx.Request != nil {
		form = svc.ctx.Request.Form
	}

	if form.Has("password") {
		form.Set("password", "*")
	}
	if form.Has("old_password") {
		form.Set("old_password", "*")
	}
	if form.Has("new_password") {
		form.Set("new_password", "*")
	}

	jsonStr, err := json.Marshal(form)
	if err != nil {
		svc.logger.Errorf("initModel data: %+v, params: %+v: %v", data, form, err)
		return form, err
	}
	if data.ModuleID == "" && svc.ctx.Param("id") != "" {
		data.ModuleID = svc.ctx.Param("id")
	}
	data.Request = string(jsonStr)
	if v, ok := svc.ctx.Get(app.USER_ID_KEY); ok {
		userID, err := strconv.Atoi(v.(string))
		if err != nil {
			svc.logger.Errorf("initModel data: %+v, params: %+v: %v", data, form, err)
			return form, err
		}
		data.UserID = uint32(userID)
	}

	return form, nil
}

// 创建操作日志为独立行为, 不返回错误, 直接告警
func (svc *OperationLogSvc) Create(data *dao.OperationLog) {
	formVal, err := svc.initModel(data)
	if err != nil {
		svc.alarm.AlertErrorMsg(fmt.Sprintf("OperationLogSvc.initModel\n%#v\n%v", data, err), message.CmsId)
	}

	if err = svc.dao.Create(formVal, data); err != nil {
		svc.alarm.AlertErrorMsg(fmt.Sprintf("OperationLogSvc.dao.Create\n%#v\n%#v\n%v", formVal, data, err), message.CmsId)
	}
}

func (svc *OperationLogSvc) List(params *form.OperationLogListRequest) ([]*dao.OperationLog, int64, *errcode.Error) {
	queryParams, err := params.Parse()
	if err != nil {
		return nil, 0, errcode.InvalidParams.WithDetails(err.Error())
	}

	data, count, err := svc.dao.List(queryParams, app.GetPager(svc.ctx))
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, count, nil
}
