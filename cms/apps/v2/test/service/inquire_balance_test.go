package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/balance/form"
	"data_backend/apps/v2/internal/inquire/balance/service"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"

	"github.com/gin-gonic/gin"
)

func TestBalanceList(t *testing.T) {
	isAdmin := false
	svc := service.NewBalanceSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.RedisClient, local.Logger, local.NewAlarm)
	summary, data, err := svc.List(&form.ListRequest{
		Pager: &app.Pager{Page: 1, PageSize: 50},
		AllRequest: form.AllRequest{
			DateTimeRangeRequest: iForm.DateTimeRangeRequest{
				DateTimeRange: [2]string{"2025-03-19 00:00:00", "2024-03-19 23:59:59"},
			},
			UserInfoRequest: cForm.UserInfoRequest{
				IsAdmin: &isAdmin,
			},
			DateTimeType: "finish",
			// UserID:            4,
			// UserName:          "微",
			// Tel:               "17396310621",
			SourceType: []int{1},
			// GachaName:         "洞洞乐预约箱子",
			UpdateAmountRange: &[2]int64{0, 999999999999999999},
		},
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	fmt.Printf("%+v\n", summary)

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestBalanceAddComment(t *testing.T) {
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Set(app.USER_ID_KEY, "2")
	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, local.Logger, local.NewAlarm)
	err := svc.AddComment(&form.AddCommentRequest{
		ID:      34734,
		Comment: "test add 2",
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	<-time.After(time.Second * 2)
}

func TestBalanceDeleteComment(t *testing.T) {
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Set(app.USER_ID_KEY, "2")
	svc := service.NewBalanceSvc(ctx, local.CMSDB, local.CenterDB, local.RedisClient, local.Logger, local.NewAlarm)
	err := svc.DeleteComment(&form.DeleteCommentRequest{
		ID:        10,
		CommentID: 2,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	<-time.After(time.Second * 2)
}
