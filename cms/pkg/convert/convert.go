package convert

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"data_backend/pkg"

	"github.com/shopspring/decimal"
)

type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func CamelCaseToSnakeCase(camelCase string) string {
	if camelCase == "" {
		return ""
	}
	strLen := len(camelCase)
	result := make([]byte, 0, strLen*2)
	j := false
	for i := 0; i < strLen; i++ {
		char := camelCase[i]
		if i > 0 && char >= 'A' && char <= 'Z' && j {
			result = append(result, '_')
		}
		if char != '_' {
			j = true
		}
		result = append(result, char)
	}
	camelCase = strings.ToLower(string(result))
	return camelCase
}

func GetInt(value interface{}) int {
	return int(GetInt64(value))
}

func GetUint(value interface{}) uint {
	return uint(GetInt64(value))
}

func GetInt32(value interface{}) int32 {
	return int32(GetInt64(value))
}

func GetUint32(value interface{}) uint32 {
	return uint32(GetInt64(value))
}

func GetInt64(value interface{}) int64 {
	if value == nil {
		return int64(0)
	}
	switch value := value.(type) {
	case uint:
		return int64(value)
	case uint8:
		return int64(value)
	case uint16:
		return int64(value)
	case uint32:
		return int64(value)
	case uint64:
		return int64(value)
	case int:
		return int64(value)
	case int8:
		return int64(value)
	case int16:
		return int64(value)
	case int32:
		return int64(value)
	case int64:
		return value
	case *uint:
		return int64(*value)
	case *uint8:
		return int64(*value)
	case *uint16:
		return int64(*value)
	case *uint32:
		return int64(*value)
	case *uint64:
		return int64(*value)
	case *int:
		return int64(*value)
	case *int8:
		return int64(*value)
	case *int16:
		return int64(*value)
	case *int32:
		return int64(*value)
	case *int64:
		return *value
	case float32:
		return int64(value)
	case float64:
		return int64(value)
	case *float32:
		return int64(*value)
	case *float64:
		return int64(*value)
	case string:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return 0
		}
		return int64(intVal)
	case *string:
		intVal, err := strconv.Atoi(*value)
		if err != nil {
			return 0
		}
		return int64(intVal)
	}
	return int64(0)
}

func GetFloat64(value interface{}) float64 {
	if value == nil {
		return float64(0)
	}
	switch value := value.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *int, *int8, *int16, *int32, *int64:
		v, _ := decimal.NewFromInt(GetInt64(value)).Float64()
		return v
	case float32:
		return float64(value)
	case float64:
		return value
	case *float32:
		return float64(*value)
	case *float64:
		return *value
	case string:
		intVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return float64(0)
		}
		return intVal
	case *string:
		intVal, err := strconv.ParseFloat(*value, 64)
		if err != nil {
			return float64(0)
		}
		return intVal
	default:
		return float64(0)
	}
}

func GetDecimal(value interface{}) decimal.Decimal {
	if value == nil {
		return decimal.NewFromInt(0)
	}
	switch value := value.(type) {
	case decimal.Decimal:
		return value
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *int, *int8, *int16, *int32, *int64:
		return decimal.NewFromInt(GetInt64(value))
	case float32:
		return decimal.NewFromFloat32(value)
	case float64:
		return decimal.NewFromFloat(value)
	case *decimal.Decimal:
		return *value
	case *float32:
		return decimal.NewFromFloat32(*value)
	case *float64:
		return decimal.NewFromFloat(*value)
	case string:
		intVal, err := decimal.NewFromString(value)
		if err != nil {
			return decimal.NewFromInt(0)
		}
		return intVal
	case *string:
		intVal, err := decimal.NewFromString(*value)
		if err != nil {
			return decimal.NewFromInt(0)
		}
		return intVal
	default:
		return decimal.NewFromInt(0)
	}
}

// GetDateTime 将字符串转为 time.Time
// time.Time{} 的值会受当前 Location 在该时间戳的历史时间中采用的 TimeZone 影响
// 一个 Location 在历史上可能采用了多个不同的 TimeZone
// TimeZone 有 DST(夏令时) LMT(地方平时) UTC 等几种
// 所以 time.Time{} 结果不一定是 0000-00-00 00:00:00
func GetDateTime(value interface{}) (time.Time, error) {
	if value == nil {
		return time.Time{}, fmt.Errorf("GetDateTime value is nil")
	}
	switch value := value.(type) {
	case time.Time:
		return value, nil
	case *time.Time:
		return *value, nil
	case string:
		if t, err := time.ParseInLocation(pkg.YEAR_FORMAT, value, pkg.Location); err == nil {
			return t, nil
		}
		if t, err := time.ParseInLocation(pkg.MONTH_FORMAT, value, pkg.Location); err == nil {
			return t, nil
		}
		if t, err := time.ParseInLocation(pkg.DATE_FORMAT, value, pkg.Location); err == nil {
			return t, nil
		}
		if t, err := time.ParseInLocation(pkg.DATE_TIME_FORMAT, value, pkg.Location); err == nil {
			return t, nil
		}
		if t, err := time.Parse(time.RFC3339, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("GetDateTime value is illegality")
}

func GetDateTimeMust(value interface{}) time.Time {
	t, _ := GetDateTime(value)
	return t
}

func GetDateStr(value interface{}) string {
	result, err := GetDateTime(value)
	if err != nil {
		return "-"
	}
	return result.Format(pkg.DATE_FORMAT)
}

func GetDateTimeStr(value interface{}) string {
	result, err := GetDateTime(value)
	if err != nil {
		return "-"
	}
	return result.Format(pkg.DATE_TIME_FORMAT)
}

func GetString(value interface{}) string {
	if value == nil {
		return ""
	}
	res := ""
	switch value := value.(type) {
	case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64,
		*uint, *uint8, *uint16, *uint32, *uint64, *int, *int8, *int16, *int32, *int64:
		res = strconv.FormatInt(GetInt64(value), 10)
	case float32:
		res = decimal.NewFromFloat32(value).String()
	case float64:
		res = decimal.NewFromFloat(value).String()
	case *float32:
		res = decimal.NewFromFloat32(*value).String()
	case *float64:
		res = decimal.NewFromFloat(*value).String()
	case time.Time:
		res = GetDateTimeStr(value)
	case string:
		res = value
	}
	return res
}
