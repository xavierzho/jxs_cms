package app

import (
	"github.com/gin-gonic/gin"
)

type OrderBy struct {
	Order string `form:"order"`
	Field string `form:"order_field"`
}

func GetOrder(ctx *gin.Context) string {
	return ctx.Query(ORDER_KEY)
}

func GetOrderField(ctx *gin.Context) string {
	return ctx.Query(ORDER_FIELD_KEY)
}

func GetOrderBy(ctx *gin.Context) OrderBy {
	return OrderBy{
		Order: GetOrder(ctx),
		Field: GetOrderField(ctx),
	}
}
