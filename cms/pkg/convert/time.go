package convert

import (
	"fmt"
	"time"

	"data_backend/pkg"

	"database/sql/driver"
)

type Time struct {
	time.Time
}

func Now() Time {
	return Time{time.Now().In(pkg.Location)}
}

func CurrentDate() Time {
	currentDateStr := time.Now().In(pkg.Location).Format(pkg.DATE_FORMAT)
	currentDate, _ := time.ParseInLocation(pkg.DATE_FORMAT, currentDateStr, pkg.Location)
	return Time{currentDate}
}

func Date(year int, month time.Month, day, hour, min, sec int) Time {
	return Time{time.Date(year, month, day, hour, min, sec, 0, pkg.Location)}
}

func Parse(format, dateTime string) (Time, error) {
	result, err := time.ParseInLocation(format, dateTime, pkg.Location)
	if err != nil {
		return Time{}, err
	}

	return Time{result}, err
}

// Value 转换类型 使数据库能插入Time类型数据
func (t *Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(v interface{}) (err error) {
	dbTime, ok := v.(time.Time)
	if ok {
		now, err := time.ParseInLocation(pkg.DATE_TIME_FORMAT, dbTime.Format(pkg.DATE_TIME_FORMAT), pkg.Location)
		if err != nil {
			return err
		}
		t.Time = now
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+pkg.DATE_TIME_FORMAT+`"`, string(data), pkg.Location)
	if err != nil {
		return err
	}
	t.Time = now
	return
}

func (t *Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(pkg.DATE_TIME_FORMAT)+2)
	b = append(b, '"')
	b = t.Time.In(pkg.Location).AppendFormat(b, pkg.DATE_TIME_FORMAT)
	b = append(b, '"')
	return b, nil
}

func (t *Time) String() string {
	return t.Time.Format(pkg.DATE_TIME_FORMAT)
}

func (t *Time) Format(layout string) string {
	return t.Time.Format(layout)
}

func (t *Time) FormatInTimeZone(layout string) string {
	return t.Time.In(pkg.Location).Format(layout)
}
