package service

import (
	"fmt"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/revenue/form"
	"data_backend/apps/v2/internal/report/revenue/service"

	"github.com/gin-gonic/gin"
)

func TestRevenueGenerate(t *testing.T) {
	svc := service.NewRevenueSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{"2024-05-09", "2024-05-09"},
		// DataTypeList: []string{"draw"},
	})
}

func TestRevenueAll(t *testing.T) {
	svc := service.NewRevenueSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	data, err := svc.All(&form.AllRequest{
		DateRange: [2]string{"2024-04-12", "2024-04-12"},
		DataType:  form.REVENUE_REPORT_TYPE_SUMMARY,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", data)
}
