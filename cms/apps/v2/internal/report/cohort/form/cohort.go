package form

import (
	"fmt"
	"sort"
	"time"

	"data_backend/apps/v2/internal/report/cohort/dao"
	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/util"
)

// TODO consume 统一 改为 pay
const (
	COHORT_TYPE_NEW_USER_ACTIVE          = "new_user_active"          // 新增留存
	COHORT_TYPE_NEW_USER_VALIDATED       = "new_user_validated"       // 有效用户
	COHORT_TYPE_NEW_USER_CONSUME         = "new_user_consume"         // 新增消费
	COHORT_TYPE_PATING_USER_ACTIVE       = "pating_user_active"       // 参与留存
	COHORT_TYPE_CONSUME_USER_ACTIVE      = "consume_user_active"      // 消费留存
	COHORT_TYPE_INVITED_NEW_USER_ACTIVE  = "invited_new_user_active"  // 受邀留存
	COHORT_TYPE_INVITED_NEW_USER_CONSUME = "invited_new_user_consume" // 受邀新增消费
)

var COHORT_TYPE_LIST_DAY_1_TO_180 = []string{
	COHORT_TYPE_NEW_USER_ACTIVE,
	// COHORT_TYPE_NEW_USER_VALIDATED,
	COHORT_TYPE_NEW_USER_CONSUME,
	COHORT_TYPE_PATING_USER_ACTIVE,
	COHORT_TYPE_CONSUME_USER_ACTIVE,
	COHORT_TYPE_INVITED_NEW_USER_ACTIVE,
	COHORT_TYPE_INVITED_NEW_USER_CONSUME,
}

var COHORT_TYPE_LIST = append(COHORT_TYPE_LIST_DAY_1_TO_180, COHORT_TYPE_NEW_USER_VALIDATED)

// cohort 报表维度为用户创建日期/消费日期, 统计项为第X天的数据
type GenerateRequest struct {
	DateRange      [2]string `form:"date_range[]" binding:"required"` // 统计日期
	UpdateDateList []string  `form:"update_date_list[]"`              // 更新数据的日期(单个日期列表); 为空表示全部
	DataTypeList   []string  `form:"data_type_list[]"`                // 为空表示全部
}

func (q *GenerateRequest) Parse() (dateRange [2]time.Time, updateDateList []time.Time, err error) {
	if err = q.Valid(); err != nil {
		return dateRange, nil, err
	}

	dateRange, err = util.ParseDateParams(q.DateRange)
	if err != nil {
		return dateRange, nil, err
	}

	for _, dateStr := range q.UpdateDateList {
		updateDate, err := time.ParseInLocation(pkg.DATE_FORMAT, dateStr, pkg.Location)
		if err != nil {
			return dateRange, nil, err
		}
		updateDateList = append(updateDateList, updateDate)
	}

	sort.Slice(updateDateList, func(i, j int) bool {
		return updateDateList[i].Before(updateDateList[j])
	})

	// 需执行的报表类型
	if len(q.DataTypeList) == 0 {
		q.DataTypeList = COHORT_TYPE_LIST
	}

	return dateRange, updateDateList, nil
}

func (q *GenerateRequest) Valid() error {
	for _, dataType := range q.DataTypeList {
		switch dataType {
		case COHORT_TYPE_NEW_USER_ACTIVE, COHORT_TYPE_NEW_USER_VALIDATED, COHORT_TYPE_NEW_USER_CONSUME:
		case COHORT_TYPE_PATING_USER_ACTIVE:
		case COHORT_TYPE_CONSUME_USER_ACTIVE:
		case COHORT_TYPE_INVITED_NEW_USER_ACTIVE, COHORT_TYPE_INVITED_NEW_USER_CONSUME:
		default:
			return fmt.Errorf("not expected data_type: " + dataType)
		}
	}

	return nil
}

type AllRequest struct {
	DateRange [2]string `form:"date_range[]" binding:"required"`
	DataType  string    `form:"data_type" binding:"required"`
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return dateRange, err
	}

	dateRange, err = util.ParseDateParams(q.DateRange)
	if err != nil {
		return dateRange, err
	}

	return dateRange, nil
}

func (q *AllRequest) Valid() error {
	switch q.DataType {
	case COHORT_TYPE_NEW_USER_ACTIVE, COHORT_TYPE_NEW_USER_VALIDATED, COHORT_TYPE_NEW_USER_CONSUME:
	case COHORT_TYPE_PATING_USER_ACTIVE:
	case COHORT_TYPE_CONSUME_USER_ACTIVE:
	case COHORT_TYPE_INVITED_NEW_USER_ACTIVE, COHORT_TYPE_INVITED_NEW_USER_CONSUME:
	default:
		return fmt.Errorf("not expected data_type: " + q.DataType)
	}

	return nil
}

func Format(dateRange [2]time.Time, data []*dao.Cohort) (result []*dao.Cohort, err error) {
	var dataMap = make(map[string]dao.Cohort)
	for _, item := range data {
		dataMap[item.Date] = *item
	}

	for cDate := dateRange[1]; !dateRange[0].After(cDate); cDate = cDate.AddDate(0, 0, -1) {
		cDateStr := cDate.Format(pkg.DATE_FORMAT)
		item := &dao.Cohort{
			DailyTypeModel: iDao.DailyTypeModel{
				Date: cDateStr,
			},
			Total:             dataMap[cDateStr].Total,
			FirstDay:          dataMap[cDateStr].FirstDay,
			SecondDay:         dataMap[cDateStr].SecondDay,
			ThirdDay:          dataMap[cDateStr].ThirdDay,
			FourthDay:         dataMap[cDateStr].FourthDay,
			FifthDay:          dataMap[cDateStr].FifthDay,
			SixthDay:          dataMap[cDateStr].SixthDay,
			SeventhDay:        dataMap[cDateStr].SeventhDay,
			FourteenthDay:     dataMap[cDateStr].FourteenthDay,
			ThirtiethDay:      dataMap[cDateStr].ThirtiethDay,
			SixtiethDay:       dataMap[cDateStr].SixtiethDay,
			NinetyDay:         dataMap[cDateStr].NinetyDay,
			No180Day:          dataMap[cDateStr].No180Day,
			FirstDayRate:      util.SaveRatio2Decimal(dataMap[cDateStr].FirstDay, dataMap[cDateStr].Total),
			SecondDayRate:     util.SaveRatio2Decimal(dataMap[cDateStr].SecondDay, dataMap[cDateStr].Total),
			ThirdDayRate:      util.SaveRatio2Decimal(dataMap[cDateStr].ThirdDay, dataMap[cDateStr].Total),
			FourthDayRate:     util.SaveRatio2Decimal(dataMap[cDateStr].FourthDay, dataMap[cDateStr].Total),
			FifthDayRate:      util.SaveRatio2Decimal(dataMap[cDateStr].FifthDay, dataMap[cDateStr].Total),
			SixthDayRate:      util.SaveRatio2Decimal(dataMap[cDateStr].SixthDay, dataMap[cDateStr].Total),
			SeventhDayRate:    util.SaveRatio2Decimal(dataMap[cDateStr].SeventhDay, dataMap[cDateStr].Total),
			FourteenthDayRate: util.SaveRatio2Decimal(dataMap[cDateStr].FourteenthDay, dataMap[cDateStr].Total),
			ThirtiethDayRate:  util.SaveRatio2Decimal(dataMap[cDateStr].ThirtiethDay, dataMap[cDateStr].Total),
			SixtiethDayRate:   util.SaveRatio2Decimal(dataMap[cDateStr].SixtiethDay, dataMap[cDateStr].Total),
			NinetyDayRate:     util.SaveRatio2Decimal(dataMap[cDateStr].NinetyDay, dataMap[cDateStr].Total),
			No180DayRate:      util.SaveRatio2Decimal(dataMap[cDateStr].No180Day, dataMap[cDateStr].Total),
		}

		result = append(result, item)
	}

	return result, nil
}
