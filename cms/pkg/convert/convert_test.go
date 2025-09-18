package convert

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"data_backend/pkg"
)

func TestGetDateTime(t *testing.T) {
	value := "2006-01-02T15:04:05Z07:00"
	fmt.Println(GetDateTime(value))
	fmt.Println(GetDateTimeStr(value))
	fmt.Println(GetDateStr(value))
}

func TestGetTime(t *testing.T) {
	tt, err := time.ParseInLocation(time.RFC3339, "2022-12-19T17:31:14+05:30", pkg.Location)
	fmt.Println(tt)
	fmt.Println(tt.Format(time.RFC3339))
	fmt.Println(err)

	tt, err = time.Parse(time.RFC3339, "2022-12-19T17:31:14+08:00")
	fmt.Println(tt)
	fmt.Println(err)

	fmt.Println(time.Time{}.In(pkg.Location))
}

func TestDecimalUnmarshal(t *testing.T) {
	value := GetDecimal(1.111111111111111)
	_value, err := json.Marshal(value)
	if err != nil {
		t.Log(err)
	}

	fmt.Println(string(_value))
}

func TestFormatMil(t *testing.T) {
	ti, _ := time.Parse(pkg.DATE_TIME_MIL_FORMAT, "2024-06-19 16:21:19.456789")
	fmt.Println(ti)
}
