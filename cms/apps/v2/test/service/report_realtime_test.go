package service_test

import (
	"fmt"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/realtime/form"
	"data_backend/apps/v2/internal/report/realtime/service"

	"github.com/gin-gonic/gin"
)

func TestRealtimeCached(t *testing.T) {
	svc := service.NewRealtimeSvc(&gin.Context{}, local.CenterDB, local.RedisClient, local.Logger)
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:00:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:10:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:20:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:30:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:40:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-01 12:50:00"})

	svc.Cached(&form.CachedRequest{DateTime: "2024-05-20 09:10:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-02 12:10:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-02 12:20:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-02 12:30:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-02 12:40:00"})
	// go svc.Cached(&form.CachedRequest{DateTime: "2024-04-02 12:50:00"})

	// <-time.After(time.Minute)
}

func TestRealtimeList(t *testing.T) {
	svc := service.NewRealtimeSvc(&gin.Context{}, local.CenterDB, local.RedisClient, local.Logger)
	data, yData, summaryData, err := svc.All()
	fmt.Println(err)

	for key, item := range data {
		fmt.Println(key)
		fmt.Println(item)
	}

	fmt.Println("==============")

	for key, item := range yData {
		fmt.Println(key)
		fmt.Println(item)
	}

	fmt.Println("==============")
	fmt.Println(summaryData)
}
