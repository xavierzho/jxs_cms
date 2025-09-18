package dao_test

import (
	"fmt"
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/realtime/dao"
	"data_backend/pkg"
)

func TestRealtimeGet(t *testing.T) {
	startTime, _ := time.Parse(pkg.DATE_FORMAT, "2023-12-01")
	endTime, _ := time.Parse(pkg.DATE_FORMAT, "2023-12-02")
	realTimeDao := dao.NewRealtimeDao(local.CenterDB, local.Logger)
	fmt.Println(realTimeDao.GetActiveUserCnt(startTime, endTime))
	fmt.Println(realTimeDao.GetPatingUserCnt(startTime, endTime))
	fmt.Println(realTimeDao.GetPayData(startTime, endTime))
	fmt.Println(realTimeDao.GetNewUserCnt(startTime, endTime))
}
