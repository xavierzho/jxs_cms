/*
输出注册的变量
*/
package app

import (
	"expvar"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Expvar(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	first := true
	report := func(key string, value interface{}) {
		if !first {
			fmt.Fprintf(ctx.Writer, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			fmt.Fprintf(ctx.Writer, "%q: %q", key, str)
		} else {
			fmt.Fprintf(ctx.Writer, "%q: %v", key, value)
		}
	}

	fmt.Fprintf(ctx.Writer, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	fmt.Fprintf(ctx.Writer, "\n}\n")
}
