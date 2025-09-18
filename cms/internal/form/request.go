package form

import (
	"time"

	"data_backend/pkg/util"
)

type DateRangeRequest struct {
	DateRange [2]string `form:"date_range[]" binding:"required"`
	dateRange [2]time.Time
}

func (q *DateRangeRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	return q.dateRange, nil
}

func (q *DateRangeRequest) Valid() (err error) {
	if q.dateRange, err = util.ParseDateParams(q.DateRange); err != nil {
		return err
	}

	return nil
}

type DateTimeRangeRequest struct {
	DateTimeRange [2]string `form:"date_time_range[]" binding:"required"`
	dateTimeRange [2]time.Time
}

func (q *DateTimeRangeRequest) Parse() (dateRange [2]time.Time, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	return q.dateTimeRange, nil
}

func (q *DateTimeRangeRequest) Valid() (err error) {
	if q.dateTimeRange, err = util.ParseDateTimeParams(q.DateTimeRange); err != nil {
		return err
	}

	return nil
}
