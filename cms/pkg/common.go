package pkg

import "time"

const (
	YEAR_FORMAT           = "2006"
	MONTH_FORMAT          = "2006-01"
	DATE_FORMAT           = "2006-01-02"
	DATE_TIME_FORMAT      = "2006-01-02 15:04:05"
	DATE_TIME_MIL_FORMAT  = "2006-01-02 15:04:05.999999999"
	TIME_FORMAT           = "15:04:05"
	SQL_DATE_FORMAT       = "%Y-%m-%d"
	SQL_DATE_TIME_FORMAT  = "%Y-%m-%d %H:%i:%s"
	FILE_DATE_FORMAT      = "20060102"
	FILE_DATE_TIME_FORMAT = "20060102150405"
)

var TimeZone = time.Local.String()
var Location = time.Local

func SetTimeZone(tz string) (err error) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return err
	}
	TimeZone = tz
	Location = location
	time.Local = Location

	return nil
}
