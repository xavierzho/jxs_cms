package requestInfo

import (
	"data_backend/internal/app"

	"github.com/gin-gonic/gin"
)

func RequestInfo() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if ctx.Request != nil && ctx.Request.URL != nil {
			ctx.Set(app.REQUEST_URL_KEY, ctx.Request.URL.Path)
		}

		ctx.Next()
	}
}
