package service

import (
	"context"
	"fmt"
	"time"

	"data_backend/apps/v2/internal/report/cohort/dao"
	"data_backend/apps/v2/internal/report/cohort/form"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

var timeInterval2field map[int]string
var timeInterval2fieldNewUserValidated map[int]string

func init() {
	timeInterval2field = map[int]string{
		0:   "first_day",
		1:   "second_day",
		2:   "third_day",
		3:   "fourth_day",
		4:   "fifth_day",
		5:   "sixth_day",
		6:   "seventh_day",
		13:  "fourteenth_day",
		29:  "thirtieth_day",
		44:  "forty_fifth_day",
		59:  "sixtieth_day",
		89:  "ninety_day",
		179: "no_180_day",
	}
	timeInterval2fieldNewUserValidated = map[int]string{
		6:   "first_day",
		7:   "second_day",
		8:   "third_day",
		9:   "fourth_day",
		10:  "fifth_day",
		11:  "sixth_day",
		12:  "seventh_day",
		13:  "fourteenth_day",
		29:  "thirtieth_day",
		44:  "forty_fifth_day",
		59:  "sixtieth_day",
		89:  "ninety_day",
		179: "no_180_day",
	}
}

type CohortSvc struct {
	logger *logger.Logger
	dao    *dao.CohortDao
}

func NewCohortSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *CohortSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".CohortSvc")))
	return &CohortSvc{
		logger: log,
		dao:    dao.NewCohortDao(engine, center, log),
	}
}

// 统计项为第x天数据
// 所以只查指定日期, 再更新某一个字段即可
func (svc *CohortSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
	// 统计日期范围, 更新日期范围
	dateRange, updateDateList, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	for cDate := dateRange[0]; !cDate.After(dateRange[1]); cDate = cDate.AddDate(0, 0, 1) {
		if e = svc.generate(params.DataTypeList, cDate, updateDateList); e != nil {
			return e
		}
	}

	return nil
}

func (svc *CohortSvc) generate(dataTypeList []string, cDate time.Time, updateDateList []time.Time) (e *errcode.Error) {
	// 需要更新的日期
	var updateDateListOther = updateDateList
	var updateDateListNewUserValidated = updateDateList
	var updateFieldOther []string
	var updateFieldNewUserValidated []string
	if len(updateDateList) == 0 {
		updateDateListOther, updateDateListNewUserValidated = getUpdateDateList(cDate)
	}
	updateFieldOther = getUpdateField(cDate, updateDateListOther, timeInterval2field)
	updateFieldNewUserValidated = getUpdateField(cDate, updateDateListNewUserValidated, timeInterval2fieldNewUserValidated)

	eg := errgroup.Group{}
	for _, dataType := range dataTypeList {
		var _updateDateList []time.Time
		var updateField []string
		_dataType := dataType
		if dataType == form.COHORT_TYPE_NEW_USER_VALIDATED {
			_updateDateList = updateDateListNewUserValidated
			updateField = updateFieldNewUserValidated
		} else {
			_updateDateList = updateDateListOther
			updateField = updateFieldOther
		}
		if len(updateField) == 0 {
			continue
		}

		switch dataType {
		case form.COHORT_TYPE_NEW_USER_ACTIVE:
			eg.Go(func() error {
				return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateNewUserActive)
			})
		case form.COHORT_TYPE_NEW_USER_VALIDATED:
			eg.Go(func() error {
				return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateNewUserValidated)
			})
		case form.COHORT_TYPE_NEW_USER_CONSUME:
			eg.Go(func() error {
				return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateNewUserConsume)
			})
		case form.COHORT_TYPE_PATING_USER_ACTIVE:
			eg.Go(func() error {
				return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GeneratePatingUserActive)
			})
		case form.COHORT_TYPE_CONSUME_USER_ACTIVE:
			eg.Go(func() error {
				return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateConsumeUserActive)
			})
		case form.COHORT_TYPE_INVITED_NEW_USER_ACTIVE: // TODO 暂无
			// eg.Go(func() error {
			// 	return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateInvitedNewUserActive)
			// })
		case form.COHORT_TYPE_INVITED_NEW_USER_CONSUME: // TODO 暂无
			// eg.Go(func() error {
			// 	return svc.generateData(_dataType, cDate, _updateDateList, updateField, svc.dao.GenerateInvitedNewUserConsume)
			// })
		default:
			eg.Go(func() error { return fmt.Errorf("generate, not expected data_type: " + _dataType) })
		}
	}

	if err := eg.Wait(); err != nil {
		return iErrcode.SQLExecFail.WithDetails(err.Error())
	}

	return nil
}

func (svc *CohortSvc) generateData(
	dataType string, cDate time.Time, updateDateList []time.Time, updateField []string,
	generateFunc func(cDate time.Time, startUpdateDate time.Time, lastUpdateDate time.Time) (data *dao.Cohort, err error),
) (err error) {
	if len(updateDateList) == 0 {
		return nil
	}
	var startUpdateDate, lastUpdateDate time.Time = updateDateList[0], updateDateList[len(updateDateList)-1]

	data, err := generateFunc(cDate, startUpdateDate, lastUpdateDate)
	if err != nil {
		return err
	}

	// 填充类型
	data.Date = cDate.Format(pkg.DATE_FORMAT)
	data.DataType = dataType
	// 保存数据
	if err = svc.dao.Save(updateField, data); err != nil {
		return err
	}

	return nil
}

// 获取当前需要执行的日期
func getUpdateDateList(date time.Time) (updateDateListOther, updateDateListNewUserValidated []time.Time) {
	cDate := time.Now()
	cDate = time.Date(cDate.Year(), cDate.Month(), cDate.Day(), 0, 0, 0, 0, pkg.Location)
	day2 := date.AddDate(0, 0, 1)
	day3 := date.AddDate(0, 0, 2)
	day4 := date.AddDate(0, 0, 3)
	day5 := date.AddDate(0, 0, 4)
	day6 := date.AddDate(0, 0, 5)
	day7 := date.AddDate(0, 0, 6)
	day8 := date.AddDate(0, 0, 7)
	day9 := date.AddDate(0, 0, 8)
	day10 := date.AddDate(0, 0, 9)
	day11 := date.AddDate(0, 0, 10)
	day12 := date.AddDate(0, 0, 11)
	day13 := date.AddDate(0, 0, 12)
	day14 := date.AddDate(0, 0, 13)
	day30 := date.AddDate(0, 0, 29)
	day60 := date.AddDate(0, 0, 59)
	day90 := date.AddDate(0, 0, 89)
	day180 := date.AddDate(0, 0, 179)

	updateDateListOther = []time.Time{date, day2, day3, day4, day5, day6, day7, day14, day30, day60, day90, day180}

	// 过滤日期
	for index := range updateDateListOther {
		if updateDateListOther[index].After(cDate) {
			updateDateListOther = updateDateListOther[:index]
			break
		}
	}

	updateDateListNewUserValidated = []time.Time{day7, day8, day9, day10, day11, day12, day13, day14, day30, day60, day90, day180}

	// 过滤日期
	for index := range updateDateListNewUserValidated {
		if updateDateListNewUserValidated[index].After(cDate) {
			updateDateListNewUserValidated = updateDateListNewUserValidated[:index]
			break
		}
	}

	return
}

func getUpdateField(date time.Time, updateDateList []time.Time, timeInterval2fileMap map[int]string) (result []string) {
	for _, updateDate := range updateDateList {
		updateField := timeInterval2fileMap[int(updateDate.Sub(date).Hours()/24)] // 两个时间都是当天零点所以可以直接按24小时计算
		if updateField == "" {
			continue
		}

		result = append(result, updateField)
	}

	return
}

func (svc *CohortSvc) All(params *form.AllRequest) (dataForm []*dao.Cohort, userCnt int64, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, 0, errcode.InvalidParams.WithDetails(err.Error())
	}

	// 有效用户数特殊处理：当前时间最近6天无数据，不进行填充
	if params.DataType == form.COHORT_TYPE_NEW_USER_VALIDATED {
		targetTime := time.Now().AddDate(0, 0, -6)
		if dateRange[0].After(targetTime) {
			return nil, 0, nil
		}
		if dateRange[1].After(targetTime) {
			dateRange[1] = time.Date(targetTime.Year(), targetTime.Month(), targetTime.Day(), 0, 0, 0, 0, pkg.Location)
		}
	}

	data, err := svc.dao.All(dateRange, params.DataType, nil)
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	data, err = form.Format(dateRange, data)
	if err != nil {
		svc.logger.Errorf("All, Format: %v", err)
		return nil, 0, errcode.TransformFail.WithDetails(err.Error())
	}

	userCnt, err = svc.getUserCnt(params.DataType, dateRange, nil)
	if err != nil {
		return nil, 0, errcode.QueryFail.WithDetails(err.Error())
	}

	return data, userCnt, nil
}

func (svc *CohortSvc) getUserCnt(dataType string, dateRange [2]time.Time, queryParams []database.QueryWhere) (int64, error) {
	switch dataType {
	case form.COHORT_TYPE_NEW_USER_ACTIVE, form.COHORT_TYPE_NEW_USER_VALIDATED, form.COHORT_TYPE_NEW_USER_CONSUME:
		return svc.dao.GetNewUserCnt(dateRange, queryParams)
	case form.COHORT_TYPE_PATING_USER_ACTIVE:
		return svc.dao.GetPatingUserCnt(dateRange, queryParams)
	case form.COHORT_TYPE_CONSUME_USER_ACTIVE:
		return svc.dao.GetPayUserCnt(dateRange, queryParams)
	case form.COHORT_TYPE_INVITED_NEW_USER_ACTIVE, form.COHORT_TYPE_INVITED_NEW_USER_CONSUME:
		// return svc.dao.GetInvitedUserCnt(dateRange, queryParams) // TODO 暂无
		return 0, nil
	default:
		err := fmt.Errorf("getUserCnt, not expected data_type: " + dataType)
		return 0, err
	}
}
