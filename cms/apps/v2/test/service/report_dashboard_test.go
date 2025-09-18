package service_test

import (
	"fmt"
	"net/http"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/dashboard/form"
	"data_backend/apps/v2/internal/report/dashboard/service"

	"github.com/gin-gonic/gin"
)

func TestDashboardGenerate(t *testing.T) {
	svc := service.NewDashboardSvc(&gin.Context{Request: &http.Request{}}, local.CMSDB, local.CenterDB, local.Logger)
	svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{"2024-05-20", "2024-05-20"},
	})
}

func TestDashboardList(t *testing.T) {
	svc := service.NewDashboardSvc(&gin.Context{Request: &http.Request{}}, local.CMSDB, local.CenterDB, local.Logger)
	data, summary, err := svc.List()
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Printf("%+v\n", summary)

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
