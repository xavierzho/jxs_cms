package service_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/bet/form"
	"data_backend/apps/v2/internal/report/bet/service"

	"github.com/gin-gonic/gin"
)

func TestBetGenerate(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := service.NewBetSvc(ctx, local.CMSDB, local.CenterDB, local.Logger)
	err := svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{"2024-04-02", "2024-04-02"},
	})
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestBetAll(t *testing.T) {
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	svc := service.NewBetSvc(ctx, local.CMSDB, local.CenterDB, local.Logger)
	data, err := svc.All(&form.AllRequest{
		DateRange: [2]string{"2024-05-05", "2024-05-06"},
		DataType:  "Gashapon",
	})
	if err != nil {
		fmt.Println(err)
		fmt.Println(err.Details())
		return
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
