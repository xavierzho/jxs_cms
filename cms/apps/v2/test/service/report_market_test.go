package service_test

import (
	"fmt"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/market/form"
	"data_backend/apps/v2/internal/report/market/service"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"

	"github.com/gin-gonic/gin"
)

func TestMarketGenerate(t *testing.T) {
	svc := service.NewMarketSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	svc.Generate(&form.GenerateRequest{
		DateRangeRequest: iForm.DateRangeRequest{
			DateRange: [2]string{"2024-04-02", "2024-04-02"},
		},
	})
}

func TestMarketList(t *testing.T) {
	svc := service.NewMarketSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	_, data, err := svc.List(&form.ListRequest{
		Pager: app.Pager{Page: 1, PageSize: 10},
		AllRequest: form.AllRequest{
			DateRangeRequest: iForm.DateRangeRequest{
				DateRange: [2]string{"2024-07-01", "2024-07-18"},
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
