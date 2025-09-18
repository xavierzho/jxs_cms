package dao_test

import (
	"fmt"
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/cohort/dao"
	"data_backend/pkg"
)

func TestCohortGet(t *testing.T) {
	startTime, _ := time.ParseInLocation(pkg.DATE_FORMAT, "2023-12-01", pkg.Location)
	endTime, _ := time.ParseInLocation(pkg.DATE_FORMAT, "2023-12-02", pkg.Location)
	cohortDao := dao.NewCohortDao(local.CMSDB, local.CenterDB, local.Logger)
	fmt.Println(cohortDao.GetNewUserCnt([2]time.Time{startTime, endTime}, nil))
	fmt.Println(cohortDao.GetNewUserCntList([2]time.Time{startTime, endTime}, nil))
	fmt.Println(cohortDao.GetPatingUserCnt([2]time.Time{startTime, endTime}, nil))
	fmt.Println(cohortDao.GetPayUserCnt([2]time.Time{startTime, endTime}, nil))
	fmt.Println(cohortDao.GetInvitedUserCnt([2]time.Time{startTime, endTime}, nil))
}
