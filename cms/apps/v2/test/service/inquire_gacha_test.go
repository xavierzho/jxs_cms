package service_test

import (
	"fmt"
	"net/http"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/gacha/form"
	"data_backend/apps/v2/internal/inquire/gacha/service"
	"data_backend/internal/app"

	"github.com/gin-gonic/gin"
)

func TestGachaOption(t *testing.T) {
	svc := service.NewGachaSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)
	fmt.Println(svc.OptionsGachaType())
}

func TestGachaRevenue(t *testing.T) {
	svc := service.NewGachaSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)

	summary, data, err := svc.GetRevenue(&form.RevenueRequest{
		Pager: &app.Pager{Page: 1, PageSize: 50, TotalRows: 0},
		// IsBoxDim: true,
		// TypeList:         []int{101, 102},
		// GachaName:        "魔神欢乐池",
		// BetRate:          &[2]int{50, 100},
		// RevenueRange:     &[2]int64{-2 << 32, 0},
		// RevenueRateRange: &[2]int{-400, 100},
	})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("summary: ", summary)

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}

}

func TestGachaDetail(t *testing.T) {
	svc := service.NewGachaSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)

	data, err := svc.GetDetail(&form.DetailRequest{
		GachaID:  229753949740077990,
		BoxOutNo: 20240510853001,
	})
	if err != nil {
		t.Error(err)
		return
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
