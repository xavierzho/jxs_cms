package local

import (
	"data_backend/internal/middleware/auth"
	iService "data_backend/internal/service"

	"github.com/gin-gonic/gin"
)

var (
	UserSvc        iService.UserSvc
	PermissionGate auth.PermissionGate
	JWT            auth.JWT
)

func SetupMiddlewareObject() error {
	UserSvc = *iService.NewUserSvc(&gin.Context{}, CMSDB, RedisClient, Logger, NewAlarm)
	PermissionGate = auth.NewPermissionGate(RedisClient, Logger, UserSvc.CheckPerm, UserSvc.CheckPermOr)
	JWT = auth.NewJWT(RedisClient, Logger, UserSvc.CheckPerm)

	return nil
}
