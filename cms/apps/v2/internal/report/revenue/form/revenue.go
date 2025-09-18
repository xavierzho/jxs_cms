package form

import (
	"fmt"
	"time"

	"data_backend/pkg/util"
)

const (
	REVENUE_DATA_TYPE_PAY     = "pay"     // 支付数据
	REVENUE_DATA_TYPE_DRAW    = "draw"    // 退款(￥)数据
	REVENUE_DATA_TYPE_BALANCE = "balance" // 钱包数据
	REVENUE_DATA_TYPE_ACTIVE  = "active"  // 活跃数据
	REVENUE_DATA_TYPE_PATING  = "pating"  // 参与数据
	REVENUE_DATA_TYPE_WASTAGE = "wastage" // 流失数据

	REVENUE_REPORT_TYPE_ACTIVE  = "active"  // 活跃用户
	REVENUE_REPORT_TYPE_PATING  = "pating"  // 参与用户
	REVENUE_REPORT_TYPE_PAY     = "pay"     // 付费数据
	REVENUE_REPORT_TYPE_DRAW    = "draw"    // 体现数据
	REVENUE_REPORT_TYPE_SUMMARY = "summary" // 营收数据
)

var revenue_data_type_list = []string{
	REVENUE_DATA_TYPE_PAY,
	REVENUE_DATA_TYPE_DRAW,
	REVENUE_DATA_TYPE_BALANCE,
	REVENUE_DATA_TYPE_ACTIVE,
	REVENUE_DATA_TYPE_PATING,
	REVENUE_DATA_TYPE_WASTAGE,
}

type GenerateRequest struct {
	DateRange    [2]string `form:"date_range[]" binding:"required"`
	DataTypeList []string  `form:"data_type[]"`
}

func (q *GenerateRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return dateRange, err
	}

	dateRange, err = util.ParseDateParams(q.DateRange)
	if err != nil {
		return dateRange, err
	}

	if len(q.DataTypeList) == 0 {
		q.DataTypeList = []string{
			REVENUE_DATA_TYPE_PAY,
			REVENUE_DATA_TYPE_DRAW,
			// REVENUE_DATA_TYPE_BALANCE, -- 钱包不更新
			REVENUE_DATA_TYPE_ACTIVE,
			REVENUE_DATA_TYPE_PATING,
			REVENUE_DATA_TYPE_WASTAGE,
		}
	}

	return dateRange, nil
}

func (q *GenerateRequest) Valid() (err error) {
	for _, dataType := range q.DataTypeList {
		switch dataType {
		case REVENUE_DATA_TYPE_PAY:
		case REVENUE_DATA_TYPE_DRAW:
		case REVENUE_DATA_TYPE_BALANCE:
		case REVENUE_DATA_TYPE_ACTIVE:
		case REVENUE_DATA_TYPE_PATING:
		case REVENUE_DATA_TYPE_WASTAGE:
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

func (q *AllRequest) Valid() (err error) {
	switch q.DataType {
	case REVENUE_REPORT_TYPE_ACTIVE:
	case REVENUE_REPORT_TYPE_PATING:
	case REVENUE_REPORT_TYPE_PAY:
	case REVENUE_REPORT_TYPE_DRAW:
	case REVENUE_REPORT_TYPE_SUMMARY:
	default:
		return fmt.Errorf("not expected data_type: " + q.DataType)
	}

	return nil
}
