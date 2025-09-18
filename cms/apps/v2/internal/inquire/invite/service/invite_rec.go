package service

import (
	"context"

	"data_backend/apps/v2/internal/inquire/invite/dao"
	"data_backend/apps/v2/internal/inquire/invite/form"
	"data_backend/pkg/database"
	"data_backend/pkg/encrypt/md5"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InviteRecSvc struct {
	logger *logger.Logger
	dao    *dao.InviteRecDao
}

func NewInviteRecSvc(ctx *gin.Context, center *gorm.DB, log *logger.Logger) *InviteRecSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".InviteRecSvc")))
	return &InviteRecSvc{
		logger: log,
		dao:    dao.NewInviteRecDao(center, log),
	}
}

// 获取查询条件
func (svc *InviteRecSvc) userCondition(params *form.AllRequest) (condition database.QueryWhereGroup) {
	switch params.UserType {
	case form.UserType_User:
		if params.UserID != 0 {
			condition = append(condition, database.QueryWhere{Prefix: "u.id = ?", Value: []any{params.UserID}})
		}
		if params.UserName != "" {
			condition = append(condition, database.QueryWhere{Prefix: "u.nickname = ?", Value: []any{params.UserName}})
		}
		if params.Tel != "" {
			condition = append(condition, database.QueryWhere{Prefix: "u.phone_num_md5 = ?", Value: []any{md5.EncodeMD5(params.Tel)}})
		}
	case form.UserType_ParentUser:
		if params.UserID != 0 {
			condition = append(condition, database.QueryWhere{Prefix: "up.id = ?", Value: []any{params.UserID}})
		}
		if params.UserName != "" {
			condition = append(condition, database.QueryWhere{Prefix: "up.nickname = ?", Value: []any{params.UserName}})
		}
		if params.Tel != "" {
			condition = append(condition, database.QueryWhere{Prefix: "up.phone_num_md5 = ?", Value: []any{md5.EncodeMD5(params.Tel)}})
		}
	}

	return condition

}

// 获取邀请列表
func (svc *InviteRecSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.InviteRec, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_summary, _data, err := svc.dao.ListAndSummary([]string{"count(0) as total"}, dateRange, svc.userCondition(&params.AllRequest), params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.Format(dateRange, _summary, _data)

	return

}

// 导出邀请列表
func (svc *InviteRecSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.InviteRec], e *errcode.Error) {
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
