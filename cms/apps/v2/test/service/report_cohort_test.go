package service

import (
	"fmt"
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/cohort/form"
	"data_backend/apps/v2/internal/report/cohort/service"

	"github.com/gin-gonic/gin"
)

func TestCohortGenerate(t *testing.T) {
	cohortGenerate()
	// go generate()
	// go generate()
	// <-time.After(time.Minute)
}

func cohortGenerate() {
	svc := service.NewCohortSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	svc.Generate(&form.GenerateRequest{
		DateRange: [2]string{"2024-04-02", "2024-04-02"},
		// DataTypeList: []string{form.COHORT_TYPE_PATING_USER_ACTIVE},
	})
}

func TestCohortGet(t *testing.T) {
	svc := service.NewCohortSvc(&gin.Context{}, local.CMSDB, local.CenterDB, local.Logger)
	dataForm, count, err := svc.All(&form.AllRequest{
		DateRange: [2]string{"2024-01-01", "2024-01-02"},
		DataType:  form.COHORT_TYPE_INVITED_NEW_USER_ACTIVE,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(count)

	for _, item := range dataForm {
		fmt.Println(item)
	}
}
