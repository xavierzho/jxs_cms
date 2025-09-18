/*
用于记录响应信息
*/
package accessLog

import (
	"io"

	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body io.Writer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func AccessLog(body io.Writer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyWriter := &AccessLogWriter{body: body, ResponseWriter: ctx.Writer}
		ctx.Writer = bodyWriter

		ctx.Next()
	}
}
