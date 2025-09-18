package app

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"data_backend/pkg"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Test struct {
	Name  string `form:"name"`
	Value string `form:"value"`
}

func TestBindAndValid(t *testing.T) {
	tt := &Test{}
	jackLog := &lumberjack.Logger{
		Filename:  "test-0000-00-00.log",
		MaxSize:   500,
		LocalTime: true,
	}
	req, _ := http.NewRequest("GET", "test?name=abc&value=123", nil)
	ctx := &gin.Context{
		Request: req,
	}
	response := NewResponse(ctx)
	fmt.Println(response.BindAndValid(ctx, tt, logger.NewLogger(ctx, jackLog)))
	fmt.Println(tt)
}

func TestPaginationDateRange(t *testing.T) {
	pager := Pager{Page: 2, PageSize: 5}
	startDate, _ := time.Parse(pkg.DATE_FORMAT, "2024-07-01")
	endDate := startDate.AddDate(0, 0, 19)

	dateRange := [2]time.Time{startDate, endDate}
	fmt.Printf("dateRange: %+v\n", dateRange)
	dateRange = pager.PaginationDateRange(dateRange)
	fmt.Printf("dateRange: %+v\n", dateRange)
	fmt.Printf("pager: %+v\n", pager)
}
