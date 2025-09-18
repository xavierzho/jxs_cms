package util

import (
	"fmt"
	"time"

	"data_backend/pkg"
)

func ParseDateParams(TimeParams [2]string) ([2]time.Time, error) {
	var dateRange [2]time.Time
	startTime, err1 := time.ParseInLocation(pkg.DATE_FORMAT, TimeParams[0], pkg.Location)
	endTime, err2 := time.ParseInLocation(pkg.DATE_FORMAT, TimeParams[1], pkg.Location)
	if err1 == nil && err2 == nil {
		if startTime.After(endTime) {
			return dateRange, fmt.Errorf("time range error")
		}
		dateRange[0] = startTime
		dateRange[1] = endTime
	} else {
		return dateRange, fmt.Errorf("time format error")
	}

	return dateRange, nil
}

func ParseDateParams2String(TimeParams [2]string) ([2]string, error) {
	dateRange, err := ParseDateParams(TimeParams)
	if err != nil {
		return [2]string{}, err
	}

	return [2]string{dateRange[0].Format(pkg.DATE_FORMAT), dateRange[1].Format(pkg.DATE_FORMAT)}, nil
}

func ParseDateTimeParams(TimeParams [2]string) ([2]time.Time, error) {
	var dateRange [2]time.Time
	startTime, err1 := time.ParseInLocation(pkg.DATE_TIME_FORMAT, TimeParams[0], pkg.Location)
	endTime, err2 := time.ParseInLocation(pkg.DATE_TIME_FORMAT, TimeParams[1], pkg.Location)
	if err1 == nil && err2 == nil {
		if startTime.After(endTime) {
			return dateRange, fmt.Errorf("time range error")
		}
		dateRange[0] = startTime
		dateRange[1] = endTime
	} else {
		return dateRange, fmt.Errorf("time format error")
	}

	return dateRange, nil
}

func ParseDateTimeParams2String(TimeParams [2]string) ([2]string, error) {
	dateRange, err := ParseDateTimeParams(TimeParams)
	if err != nil {
		return [2]string{}, err
	}

	return [2]string{dateRange[0].Format(pkg.DATE_TIME_FORMAT), dateRange[1].Format(pkg.DATE_TIME_FORMAT)}, nil
}
