package recovery

import (
	"fmt"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
)

type Recovery struct {
	logger *logger.Logger
	alarm  message.Alarm
}

func NewRecovery(logger *logger.Logger, alarm message.Alarm) Recovery {
	return Recovery{
		logger: logger,
		alarm:  alarm,
	}
}

func (r Recovery) Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				r.logger.WithContext(ctx).WithCallersFrames().Errorf("panic recover: %v", err)
				r.alarm.SendMsg(fmt.Sprintf("[%s] %v", time.Now().Format(pkg.DATE_TIME_FORMAT), err), message.CmsId)

				app.NewResponse(ctx).ToErrorResponse(errcode.ServerError)
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
