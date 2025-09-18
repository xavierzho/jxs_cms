package local

import (
	"context"
	"net/http"

	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// 服务自身context, 用于保存如版本, 服务名等信息
// 内容不可修改
var Ctx context.Context
var GinCtx *gin.Context // TODO 所有需要 local.Ctx 的地方都换成 GinCtx

var Module = logger.NewCustomInfo("v2")

func init() {
	Ctx = context.WithValue(context.Background(), logger.ModuleKey, Module)
	GinCtx = &gin.Context{Request: &http.Request{}}
	GinCtx.Request = GinCtx.Request.WithContext(Ctx)
}
