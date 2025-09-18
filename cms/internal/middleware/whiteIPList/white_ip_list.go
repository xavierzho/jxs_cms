package whiteIPList

import (
	"data_backend/internal/app"
	"data_backend/internal/global"

	"github.com/gin-gonic/gin"
)

func ServerWhiteIPList() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		flag := false
		if global.ServerSetting.RunMode == global.RUN_MODE_DEBUG {
			flag = true
		} else {
			for _, host := range global.ServerSetting.WhiteList {
				if ip == host {
					flag = true
					break
				}
			}
		}

		if !flag {
			response := app.NewResponse(ctx)
			response.ToResponseDetail(0, gin.H{})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
