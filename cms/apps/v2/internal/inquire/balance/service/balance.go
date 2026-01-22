package service

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/inquire/balance/dao"
	"data_backend/apps/v2/internal/inquire/balance/form"
	"data_backend/internal/global"
	iService "data_backend/internal/service"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/excel"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"
	"data_backend/pkg/redisdb"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const BALANCE_COMMENT = "balanceComment"

func balanceCommentRKey(id int64) string {
	return fmt.Sprintf("%s:%d", BALANCE_COMMENT, id)
}

type BalanceSvc struct {
	ctx      *gin.Context
	engine   *gorm.DB
	rdb      *redisdb.RedisClient
	logger   *logger.Logger
	userSvc  *iService.UserSvc
	dao      *dao.BalanceDao
	newAlarm func(log *logger.Logger) message.Alarm
}

func NewBalanceSvc(ctx *gin.Context, engine *gorm.DB, center *gorm.DB, rdb *redisdb.RedisClient, log *logger.Logger, newAlarm func(log *logger.Logger) message.Alarm) *BalanceSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".BalanceSvc")))
	return &BalanceSvc{
		ctx:      ctx,
		engine:   engine,
		rdb:      rdb,
		logger:   log,
		dao:      dao.NewBalanceDao(center, log),
		newAlarm: newAlarm,
	}
}

func (svc *BalanceSvc) OptionsSourceType() []map[string]string {
	return []map[string]string{
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "2")},
		{"value": "3", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "3")},
		{"value": "12", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "12")},
		{"value": "13", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "13")},
		{"value": "101", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "101")},
		{"value": "102", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "102")},
		{"value": "103", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "103")},
		{"value": "104", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "104")},
		{"value": "105", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "105")},
		{"value": "106", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "106")},
		{"value": "201", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "201")},
		{"value": "202", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "202")},
		{"value": "203", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "203")},
		{"value": "204", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "204")},
		{"value": "301", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "301")},
		{"value": "400", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "400")},
		{"value": "601", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "601")},
		{"value": "100005", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100005")},
		{"value": "999999", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "999999")},
	}
}

func (svc *BalanceSvc) OptionsChannelType() []map[string]string {
	return []map[string]string{
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "pay.channelType", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "pay.channelType", "2")},
		{"value": "3", "label": global.I18n.T(svc.ctx.Request.Context(), "pay.channelType", "3")},
	}
}

func (svc *BalanceSvc) OptionsPaySourceType() []map[string]string {
	return []map[string]string{
		{"value": "100", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "100")},
		{"value": "201", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "201")},
		{"value": "202", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "202")},
		{"value": "301", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "301")},
		{"value": "601", "label": global.I18n.T(svc.ctx.Request.Context(), "source_type", "601")},
	}
}

func (svc *BalanceSvc) List(params *form.ListRequest) (summary map[string]any, data []*form.Balance, e *errcode.Error) {
	_, queryParams, err := params.Parse()
	if err != nil {
		return nil, nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_summary, _data, err := svc.dao.List(params.DateTimeType, queryParams, params.Pager)
	if err != nil {
		return nil, nil, errcode.QueryFail.WithDetails(err.Error())
	}

	summary, data = form.Format(svc.ctx, _summary, _data)

	return
}

// 记录操作日志
func (svc *BalanceSvc) AddComment(params *form.AddCommentRequest) (e *errcode.Error) {
	// 加锁
	if err := svc.rdb.Lock(svc.ctx, balanceCommentRKey(params.ID)); err != nil {
		svc.logger.Errorf("AddComment Lock: %v", err)
		return errcode.UpdateFail.WithDetails(e.Error())
	}
	defer func() {
		if unlockErr := svc.rdb.Unlock(svc.ctx, balanceCommentRKey(params.ID)); unlockErr != nil {
			svc.logger.Errorf("AddComment Unlock: %v", unlockErr)
		}
	}()

	svc.userSvc = iService.NewUserSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)

	// 获取用户
	operator, e := svc.userSvc.CurrentUser()
	if e.Is(errcode.UnauthorizedTokenError) {
		return e
	} else if e != nil {
		return errcode.UpdateFail.WithDetails(e.Error())
	}

	// 获取记录
	data, err := svc.dao.First(database.QueryWhereGroup{
		{Prefix: "bl.id", Value: []any{params.ID}},
	})
	if err != nil {
		return errcode.QueryFail.WithDetails(err.Error())
	}
	if data == nil {
		return errcode.QueryFail.WithDetails("Invalid record")
	}

	// add
	var maxID int64
	for _, item := range data.Comment {
		if maxID < item.ID {
			maxID = item.ID
		}
	}
	err = svc.dao.AddComment(params.ID, &dao.BalanceComment{
		ID:        maxID + 1,
		CreatedAt: time.Now().Format(pkg.DATE_TIME_FORMAT),
		UserID:    operator.ID,
		Comment:   params.Comment,
	})
	if err != nil {
		return errcode.UpdateFail.WithDetails(err.Error())
	}

	return nil
}

// 需要校验 是否是该用户添加的 comment。若不是添加用户也不是admin则不能删除
// 记录操作日志
func (svc *BalanceSvc) DeleteComment(params *form.DeleteCommentRequest) (e *errcode.Error) {
	// 加锁
	if err := svc.rdb.Lock(svc.ctx, balanceCommentRKey(params.ID)); err != nil {
		svc.logger.Errorf("DeleteComment Lock: %v", err)
		return errcode.UpdateFail.WithDetails(e.Error())
	}
	defer func() {
		if unlockErr := svc.rdb.Unlock(svc.ctx, balanceCommentRKey(params.ID)); unlockErr != nil {
			svc.logger.Errorf("DeleteComment Unlock: %v", unlockErr)
		}
	}()

	svc.userSvc = iService.NewUserSvc(svc.ctx, svc.engine, svc.rdb, svc.logger, svc.newAlarm)

	// 获取用户
	operator, e := svc.userSvc.CurrentUser()
	if e.Is(errcode.UnauthorizedTokenError) {
		return e
	} else if e != nil {
		return errcode.UpdateFail.WithDetails(e.Error())
	}

	// 获取记录
	data, err := svc.dao.First(database.QueryWhereGroup{
		{Prefix: "bl.id", Value: []any{params.ID}},
	})
	if err != nil {
		return errcode.QueryFail.WithDetails(err.Error())
	}
	if data == nil {
		return errcode.QueryFail.WithDetails("invalid record")
	}

	// 校验
	var commentIndex int = -1
	for ind, item := range data.Comment {
		if item.ID == params.CommentID {
			commentIndex = ind
			break
		}
	}
	if commentIndex == -1 {
		return errcode.ErrorParam.WithDetails("invalid comment id")
	}
	if !operator.IsAdmin() && data.Comment[commentIndex].UserID != operator.ID {
		return errcode.PermissionDenied.WithDetails("you can only modify records you create")
	}

	// delete
	err = svc.dao.DeleteComment(params.ID, commentIndex)
	if err != nil {
		return errcode.UpdateFail.WithDetails(err.Error())
	}

	return nil
}

func (svc *BalanceSvc) Export(params *form.AllRequest) (data *excel.Excel[*form.Balance], e *errcode.Error) {
	dateTimeRange, queryParams, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	_data, err := svc.dao.All(params.DateTimeType, queryParams)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format2Excel(svc.ctx, dateTimeRange, _data)
	if err != nil {
		svc.logger.Errorf("Export Format2Excel err: %v", err.Error())
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return
}

func (svc *BalanceSvc) OptionsBalanceType() []map[string]string {
	return []map[string]string{
		{"value": "0", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "0")},
		{"value": "1", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "1")},
		{"value": "2", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "2")},
		{"value": "3", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "3")},
		{"value": "10", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "10")},
		{"value": "11", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "11")},
		{"value": "20", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "20")},
		{"value": "100", "label": global.I18n.T(svc.ctx.Request.Context(), "balance_type", "100")},
	}
}
