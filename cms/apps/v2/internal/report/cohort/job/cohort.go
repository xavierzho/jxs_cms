package job

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/cohort/form"
	"data_backend/apps/v2/internal/report/cohort/service"
	"data_backend/pkg"
	"data_backend/pkg/logger"
	"data_backend/pkg/message"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// 更新当天的数据
type CohortJob struct {
	now    string
	ctx    *gin.Context
	logger *logger.Logger
	alarm  message.Alarm
}

func NewCohortJob() *CohortJob {
	log := local.JobWorkerLogger.WithContext(context.WithValue(local.JobWorkerLogger.Context, logger.ModuleKey, local.JobWorkerLogger.ModuleKey().Add(".CohortJob")))
	ctx := &gin.Context{
		Request: &http.Request{},
	}
	ctx.Request = ctx.Request.WithContext(log.Context)
	return &CohortJob{
		ctx:    ctx,
		logger: log,
		alarm:  local.NewAlarm(log),
	}
}

func (*CohortJob) Name() string {
	return "CohortJob"
}

func (j *CohortJob) Run() {
	local.JobWorker.AddJobToQueue(j.Name())
}

func (j *CohortJob) Work() {
	now := time.Now()
	if j.now != "" {
		_now, err := time.ParseInLocation(pkg.DATE_FORMAT, j.now, pkg.Location)
		if err == nil {
			now = _now
		}
	}

	eg := errgroup.Group{}
	svc := service.NewCohortSvc(j.ctx, local.CMSDB, local.CenterDB, j.logger)
	// 更新历史七天前数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange:      [2]string{now.AddDate(0, 0, -6).Format(pkg.DATE_FORMAT), now.Format(pkg.DATE_FORMAT)},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST_DAY_1_TO_180,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 1~7th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新历史第七~第十三天的数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -12).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -6).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   []string{form.COHORT_TYPE_NEW_USER_VALIDATED},
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 7~13th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新14天前的cohort数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -13).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -13).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 14th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新30天前的cohort数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -29).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -29).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 30th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新60天前的cohort数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -59).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -59).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 60th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新90天前的cohort数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -89).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -89).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 90th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	// 更新180天前的cohort数据
	eg.Go(func() error {
		err := svc.Generate(&form.GenerateRequest{
			DateRange: [2]string{
				now.AddDate(0, 0, -179).Format(pkg.DATE_FORMAT),
				now.AddDate(0, 0, -179).Format(pkg.DATE_FORMAT),
			},
			UpdateDateList: []string{now.Format(pkg.DATE_FORMAT)},
			DataTypeList:   form.COHORT_TYPE_LIST,
		})
		if err != nil {
			j.alarm.AlertErrorMsg(fmt.Sprintf("CohortSvc.Generate %s 180th: %v", now.Format(pkg.DATE_FORMAT), err), message.CmsId)
			return err
		}

		return nil
	})

	eg.Wait()
}
