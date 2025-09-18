package contextTimeout

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func ContextTimeout(t time.Duration) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		_ctx, cancel := context.WithTimeout(ctx.Request.Context(), t)
		defer cancel()

		ctx.Request = ctx.Request.WithContext(_ctx)
		ctx.Next()
	}
}
