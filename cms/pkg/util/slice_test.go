package util

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"data_backend/pkg"

	"github.com/shopspring/decimal"
)

func TestDemo(t *testing.T) {
	d, e := decimal.NewFromInt(100).Div(decimal.NewFromInt(3)).Float64()
	if e {
		fmt.Println(e, d)
	} else {
		fmt.Println(e, d)
	}
}

func TestPerCombination(t *testing.T) {
	groupByFieldsList := PerCombination([]string{"room_type", "user_type"})
	fmt.Println(groupByFieldsList)
	for ind, groupByFields := range groupByFieldsList {
		var groupByStr string = ""
		var groupByUserStrList []string
		var groupByUserStr string = ""
		if len(groupByFields) > 0 {
			groupByStr = strings.Join(groupByFields, ",") + ","
			for _, groupByField := range groupByFields {
				groupByUserStrList = append(groupByUserStrList, "u."+groupByField)
			}
			groupByUserStr = strings.Join(groupByUserStrList, ",") + ","
		}
		fmt.Println(ind, groupByStr)
		fmt.Println(ind, groupByUserStr)
	}
}

func TestPerCombination2(t *testing.T) {
	groupByFieldsList := PerCombination([]string{"room_type", "user_type"})
	fmt.Println(groupByFieldsList)
	for ind, groupByFields := range groupByFieldsList {
		var joinStrList []string
		var joinStr string = ""
		var groupByUserStrList []string
		var groupByUserStr string = ""
		var groupByStrList []string
		var selectStrList []string
		var groupByStr string = ""
		var selectStr string = ""
		if len(groupByFields) > 0 {
			for _, groupByField := range groupByFields {
				groupByUserStrList = append(groupByUserStrList, "u."+groupByField)
				groupByStrList = append(groupByStrList, fmt.Sprintf("ifnull(recharge.%s, withdraw.%s)", groupByField, groupByField))
				selectStrList = append(selectStrList, fmt.Sprintf("ifnull(recharge.%s, withdraw.%s) as %s", groupByField, groupByField, groupByField))
				joinStrList = append(joinStrList, fmt.Sprintf("%s.%s = %s.%s", "withdraw", groupByField, "recharge", groupByField))
			}
			groupByUserStr = strings.Join(groupByUserStrList, ",") + ","
			groupByStr = strings.Join(groupByStrList, ",") + ","
			selectStr = strings.Join(selectStrList, ",") + ","
			joinStr = strings.Join(joinStrList, " AND ") + " AND "
		}
		fmt.Println(ind, groupByStr)
		fmt.Println(ind, groupByUserStr)
		fmt.Println(ind, joinStr)
		fmt.Println(ind, selectStr)
	}
}

func TestTimeRange2StrDateRangeSlice(t *testing.T) {
	sTime, _ := time.Parse(pkg.DATE_FORMAT, "2022-12-01")
	eTime, _ := time.Parse(pkg.DATE_FORMAT, "2022-12-04")
	fmt.Println(TimeRange2StrDateRangeSlice([2]time.Time{sTime, eTime}, false))
	fmt.Println(TimeRange2StrDateRangeSlice([2]time.Time{sTime, eTime}, true))
}

func TestIsContainAllStringSlice(t *testing.T) {
	fmt.Println(IsContainAllStringSlice(
		[]string{"0"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAllStringSlice(
		[]string{"0", "1"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAllStringSlice(
		[]string{"1"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAllStringSlice(
		[]string{},
		[]string{"1", "2"},
	))
}

func TestIsContainAnyStringSlice(t *testing.T) {
	fmt.Println(IsContainAnyStringSlice(
		[]string{"0"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAnyStringSlice(
		[]string{"0", "1"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAnyStringSlice(
		[]string{"1"},
		[]string{"1", "2"},
	))
	fmt.Println(IsContainAnyStringSlice(
		[]string{},
		[]string{"1", "2"},
	))
}
