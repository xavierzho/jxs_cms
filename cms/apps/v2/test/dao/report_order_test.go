package dao_test

import (
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/order/dao"
	"data_backend/pkg"
)

func TestDaoGenerate(t *testing.T) {
	cTime, _ := time.Parse(pkg.DATE_TIME_FORMAT, "2025-01-17 00:00:00")

	orderDao := dao.NewDeliveryOrderDao(local.CMSDB, local.CenterDB, local.Logger)
	orderDao.Generate(cTime, nil)
}
