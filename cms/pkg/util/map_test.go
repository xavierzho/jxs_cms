package util

import (
	"fmt"
	"testing"
	"time"
)

func TestInsertEmptyDataAndSort(t *testing.T) {
	startTime := time.Now()
	endTime := time.Now().AddDate(0, 0, 5)
	dateRange := [2]time.Time{startTime, endTime}

	type TestStruct struct {
		Date string `json:"date"`
		Val  int
	}

	data := []TestStruct{
		{Date: "2023-05-05", Val: 1},
		// {Date: "2023-05-06", Val: 2},
		{Date: "2023-05-07", Val: 3},
	}

	dataList, err := insertEmptyDataAndSort(dateRange, data, "date", "", "day")
	if err != nil {
		fmt.Println("err: ", err)
		return
	}

	for _, item := range dataList {
		fmt.Println(item)
	}

}

func TestMergeSliceMapByThemKey(t *testing.T) {
	a := []map[string]interface{}{
		{"a": 1, "b": 2, "c": 3},
		{"a": 2},
	}

	b := []map[string]interface{}{
		{"a": 1, "b": 20, "d": 40},
		{"a": 3},
	}

	c := MergeSliceMapByThemKey(a, b, "a", "a")

	fmt.Printf("%+v\n", a)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", c)

	c[0]["e"] = 11
	c[1]["e"] = 11
	c[2]["e"] = 11

	fmt.Printf("%+v\n", a)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", c)
}

func TestMergeSliceMapByThemKeyList(t *testing.T) {
	a := []map[string]interface{}{
		{"a": 1, "b": 2, "c": 3},
		{"a": 2},
	}

	b := []map[string]interface{}{
		{"a": 1, "b": 20, "d": 40},
		{"a": 3},
	}

	c, err := MergeSliceMapByThemKeyList(a, b, []string{"a", "a"}, []string{"a", "a"})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", a)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", c)

	c[0]["e"] = 11
	c[1]["e"] = 11
	c[2]["e"] = 11

	fmt.Printf("%+v\n", a)
	fmt.Printf("%+v\n", b)
	fmt.Printf("%+v\n", c)
}
