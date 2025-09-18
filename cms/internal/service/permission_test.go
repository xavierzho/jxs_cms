package service_test

import (
	"fmt"
	"testing"

	"data_backend/internal/global"
	"data_backend/internal/service"
	"data_backend/pkg/errcode"

	"github.com/gin-gonic/gin"
)

func TestRefresh(t *testing.T) {
	var e *errcode.Error
	svc := service.NewPermissionSvc(&gin.Context{}, db, log, global.NewAlarm)
	e = svc.Refresh()
	fmt.Println(e)
}
