package auth

import (
	"data_backend/internal/app"
	"data_backend/pkg/convert"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/redisdb"

	"github.com/gin-gonic/gin"
)

type PermissionGate struct {
	rdb         *redisdb.RedisClient
	logger      *logger.Logger
	checkPerm   func(ctx *gin.Context, userID uint32, permList []string) (bool, error)
	checkPermOr func(ctx *gin.Context, userID uint32, permList []string) (bool, error)
}

func NewPermissionGate(
	rdb *redisdb.RedisClient, logger *logger.Logger,
	checkPerm, checkPermOr func(ctx *gin.Context, userID uint32, permList []string) (bool, error),
) PermissionGate {
	return PermissionGate{
		rdb:         rdb,
		logger:      logger,
		checkPerm:   checkPerm,
		checkPermOr: checkPermOr,
	}
}

func getUserInfo(ctx *gin.Context) (uint32, *errcode.Error) {
	// 查询是否存在权限
	userID, ok := ctx.Get(app.USER_ID_KEY)
	if !ok {
		return 0, errcode.PermissionDenied.WithDetails("user info not exist")
	}
	intUserID, err := convert.StrTo(userID.(string)).UInt32()
	if err != nil {
		return 0, errcode.PermissionDenied.WithDetails("user info error")
	}

	return intUserID, nil
}

func (p PermissionGate) CheckPerm(perms ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 查询是否存在权限
		var eCode *errcode.Error
		p.logger = p.logger.WithContext(ctx)
		defer app.DeferResponse(ctx, &eCode)

		intUserID, eCode := getUserInfo(ctx)
		if eCode != nil {
			return
		}
		if flag, checkErr := p.checkPerm(ctx, intUserID, perms); checkErr != nil {
			eCode = errcode.ServerError.WithDetails(checkErr.Error())
			return
		} else if !flag {
			eCode = errcode.PermissionDenied
			return
		}
	}
}

func (p PermissionGate) CheckPermOr(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var eCode *errcode.Error
		p.logger = p.logger.WithContext(ctx)
		defer app.DeferResponse(ctx, &eCode)

		intUserID, eCode := getUserInfo(ctx)
		if eCode != nil {
			return
		}
		if flag, checkErr := p.checkPermOr(ctx, intUserID, permissions); checkErr != nil {
			eCode = errcode.ServerError.WithDetails(checkErr.Error())
			return
		} else if !flag {
			eCode = errcode.PermissionDenied
			return
		}
	}
}
