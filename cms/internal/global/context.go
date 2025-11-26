package global

import (
	"context"

	"data_backend/pkg/logger"
)

// Ctx 服务自生context, 用于保存如版本, 服务名等信息
// 内容不可修改
var Ctx context.Context // TODO 改为 gin.Context

var Module = logger.NewCustomInfo("global")

func init() {
	Ctx = context.WithValue(context.Background(), logger.ModuleKey, Module)
}
