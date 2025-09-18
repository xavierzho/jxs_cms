package service_test

import (
	"fmt"
	"net/http"
	"testing"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/item/form"
	"data_backend/apps/v2/internal/inquire/item/service"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"

	"github.com/gin-gonic/gin"
)

func TestItemOption(t *testing.T) {
	svc := service.NewItemSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)
	fmt.Println(svc.OptionsLogType())
}

func TestItemLog(t *testing.T) {
	svc := service.NewItemSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)

	summary, data, err := svc.GetLog(&form.LogRequest{
		Pager: &app.Pager{Page: 1, PageSize: 30, TotalRows: 0},
		LogAllRequest: form.LogAllRequest{
			DateTimeRangeRequest: iForm.DateTimeRangeRequest{
				DateTimeRange: [2]string{"2024-05-10 18:00:00", "2024-05-11 18:00:00"},
			},
			// UserID:            26,
			// UserName:          "哈哈哈",
			// Tel:               "17396310621",
			// LogTypeList:       []int{101, 200},
			// GachaName:         "扭蛋机",
			// UpdateAmountRange: &[2]int64{-10000, 10000},
			// InnerPriceRange:   &[2]int64{0, 10000},
			// ShowPriceRange:    &[2]int64{-10000, 10000},
		},
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

func TestItemDetail(t *testing.T) {
	svc := service.NewItemSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)

	params := &form.DetailRequest{
		ID:        "39",
		LogType:   200,
		LevelType: 0,
	}

	data, err := svc.GetDetail(params)
	if err != nil {
		t.Error(err)
		return
	}

	if params.LogType == 200 {
		for _, i := range data.([]any)[0].([]*form.MarketItemDetail) {
			fmt.Printf("%+v\n", i)
		}

		for _, i := range data.([]any)[1].([]*form.MarketItemDetail) {
			fmt.Printf("%+v\n", i)
		}
	} else {
		for _, i := range data.([]*form.BetItemDetail) {
			fmt.Printf("%+v\n", i)
		}
	}

}

func TestItemGetDetail(t *testing.T) {
	svc := service.NewItemSvc(&gin.Context{Request: &http.Request{}}, local.CenterDB, local.Logger)

	params := &form.DetailAllRequest{
		LogAllRequest: form.LogAllRequest{
			DateTimeRangeRequest: iForm.DateTimeRangeRequest{
				DateTimeRange: [2]string{"2024-07-01 00:00:00", "2024-07-02 00:00:00"},
			},
			UserInfoRequest: cForm.UserInfoRequest{
				UserID:   123,
				UserName: "456",
				Tel:      "789",
			},
			LogTypeList:       nil,
			GachaName:         "test",
			UpdateAmountRange: &[2]int64{1, 2},
			ShowPriceRange:    &[2]int64{3, 4},
			InnerPriceRange:   &[2]int64{5, 6},
		},
	}

	svc.ExportDetail(params)
}
