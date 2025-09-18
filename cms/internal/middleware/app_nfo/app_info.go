package app_info

import (
	"data_backend/internal/app"
	"data_backend/internal/global"

	"github.com/gin-gonic/gin"
)

func APPInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(app.APP_NAME_KEY, global.ServerSetting.ServerName)
		ctx.Set(app.APP_VERSION_KEY, global.ServerSetting.Version)
		ctx.Next()
	}
}
