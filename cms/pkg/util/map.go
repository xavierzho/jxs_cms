package util

import (
	"encoding/json"
	"sort"
	"time"

	"data_backend/pkg"
	"data_backend/pkg/convert"
)

type QuerySliceMap struct {
	Field string
	Value interface{}
}

// 按照给出的条件获取[]map[string]interface{} 的某些字段值
// 类似left join 但只能取第一个值
func GetFirstMapData(data []map[string]interface{}, keyList []QuerySliceMap, valueKeyList []string) []interface{} {
	if len(data) == 0 {
		return nil
	}
	if len(keyList) == 0 {
		var result []interface{}
		for _, valueKey := range valueKeyList {
			result = append(result, data[0][valueKey])
		}
		return result
	}

	for _, item := range data {
		flag := true
		for _, key := range keyList {
			if item[key.Field] != key.Value {
				flag = false
				break
			}
		}
		if flag {
			var result []interface{}
			for _, valueKey := range valueKeyList {
				result = append(result, item[valueKey])
			}
			return result
		}
	}
	return nil
}

// 用于SliceMap多条件合并前 生成合并key
func GenMergeKeyMap(data map[string]interface{}, keyList []string) (string, error) {
	mergeKeyValues := make(map[string]interface{})
	for _, key := range keyList {
		if value, ok := data[key]; ok {
			mergeKeyValues[key] = value
		}
	}
	result, err := json.Marshal(mergeKeyValues)
	return string(result), err

}

func GenMergeKeyList(data map[string]interface{}, keyList []string) (string, error) {
	mergeKeyValues := []interface{}{}
	for _, key := range keyList {
		if value, ok := data[key]; ok {
			mergeKeyValues = append(mergeKeyValues, value)
		}
	}
	result, err := json.Marshal(mergeKeyValues)
	return string(result), err

}

// 按照 key 将两个map的数据合在一起
func MergeSliceMapByKey(sourceData, mergeData []map[string]interface{}, key string) []map[string]interface{} {
	return MergeSliceMapByThemKey(sourceData, mergeData, key, key)
}

func MergeSliceMapByThemKey(sourceData, mergeData []map[string]interface{}, sourceDataKey, mergeDataKey string) (result []map[string]interface{}) {
	if len(mergeData) == 0 {
		return sourceData
	}

	mapData := make(map[string]map[string]interface{})
	for i := 0; i < len(mergeData); i++ {
		keyValue := convert.GetString(mergeData[i][mergeDataKey])
		mapData[keyValue] = make(map[string]interface{})
		for key, value := range mergeData[i] {
			mapData[keyValue][key] = value
		}
	}

	for i := 0; i < len(sourceData); i++ {
		sourceItem := make(map[string]interface{})
		for key, val := range sourceData[i] {
			sourceItem[key] = val
		}

		keyValue := convert.GetString(sourceData[i][sourceDataKey])
		if mergeItem, ok := mapData[keyValue]; ok {
			for key, val := range mergeItem {
				sourceItem[key] = val
			}
			delete(mapData, keyValue)
		}

		result = append(result, sourceItem)

	}
	for _, item := range mapData {
		result = append(result, item)
	}

	return result
}

func MergeSliceMapByKeyList(sourceData, mergeData []map[string]interface{}, keyList []string) ([]map[string]interface{}, error) {
	return MergeSliceMapByThemKeyList(sourceData, mergeData, keyList, keyList)
}

func MergeSliceMapByThemKeyList(sourceData, mergeData []map[string]interface{}, sourceDataKey, mergeDataKey []string) (result []map[string]interface{}, err error) {
	if len(mergeData) == 0 {
		return sourceData, nil
	}

	mapData := make(map[string]map[string]interface{})
	for i := 0; i < len(mergeData); i++ {
		keyValue, err := GenMergeKeyList(mergeData[i], mergeDataKey)
		if err != nil {
			return nil, err
		}

		mapData[keyValue] = make(map[string]interface{})
		for key, value := range mergeData[i] {
			mapData[keyValue][key] = value
		}
	}

	for i := 0; i < len(sourceData); i++ {
		sourceItem := make(map[string]interface{})
		for key, val := range sourceData[i] {
			sourceItem[key] = val
		}

		keyValue, err := GenMergeKeyList(sourceData[i], sourceDataKey)
		if err != nil {
			return nil, err
		}

		if mergeItem, ok := mapData[keyValue]; ok {
			for key, val := range mergeItem {
				sourceItem[key] = val
			}
			delete(mapData, keyValue)
		}

		result = append(result, sourceItem)
	}
	for _, item := range mapData {
		result = append(result, item)
	}

	return result, nil
}

func PluckSliceFields(data []map[string]interface{}, field string) []string {
	arr := make([]string, 0, len(data))
	for i := 0; i < len(data); i++ {
		fieldVal := convert.GetString(data[i][field])
		if fieldVal == "" {
			continue
		}
		arr = append(arr, fieldVal)
	}
	return arr
}

// 根据时间语义 key 进行时间排序
func SliceMapSortByTimeField(mapData []map[string]interface{}, field, sortBy, key string) {
	if len(mapData) == 0 {
		return
	}
	if _, ok := mapData[0][field]; !ok {
		return
	}
	layout := pkg.DATE_TIME_FORMAT
	switch key {
	case "day":
		layout = pkg.DATE_FORMAT
	case "time":
		layout = pkg.TIME_FORMAT
	case "month":
		layout = pkg.MONTH_FORMAT
	case "year":
		layout = pkg.YEAR_FORMAT
	}

	sort.Slice(mapData, func(i, j int) bool {
		iTime, _ := time.Parse(layout, convert.GetString(mapData[i][field]))
		jTime, _ := time.Parse(layout, convert.GetString(mapData[j][field]))
		if sortBy == "desc" {
			return iTime.After(jTime)
		} else {
			return iTime.Before(jTime)
		}
	})
}

// 根据时间语义 key 插入空白时间数据
func InsertEmptyDataByTimeKey(dateRange [2]time.Time, mapData []map[string]interface{}, field, key string) []map[string]interface{} {
	addDateArgs := []int{0, 0, 1}
	layout := pkg.DATE_FORMAT
	switch key {
	case "month":
		addDateArgs = []int{0, 1, 0}
		layout = pkg.MONTH_FORMAT
	case "year":
		addDateArgs = []int{1, 0, 0}
		layout = pkg.YEAR_FORMAT
	}

	existDate := make(map[string]struct{})
	for i := 0; i < len(mapData); i++ {
		existDate[convert.GetDateTimeMust(mapData[i][field]).Format(layout)] = struct{}{}
	}
	notExistMap := []map[string]interface{}{}

	for cDate := dateRange[0]; cDate.Format(layout) <= dateRange[1].Format(layout); cDate = cDate.AddDate(addDateArgs[0], addDateArgs[1], addDateArgs[2]) {
		if _, ok := existDate[cDate.Format(layout)]; !ok {
			notExistMap = append(notExistMap, map[string]interface{}{
				field: cDate.Format(layout),
			})
		}
	}

	mapData = append(mapData, notExistMap...)
	return mapData
}

func InsertEmptyYearDataAndSort(dateRange [2]time.Time, data interface{}, field, sortBy string) (dataList []map[string]interface{}, err error) {
	return insertEmptyDataAndSort(dateRange, data, field, sortBy, "year")
}

func InsertEmptyMonthDataAndSort(dateRange [2]time.Time, data interface{}, field, sortBy string) (dataList []map[string]interface{}, err error) {
	return insertEmptyDataAndSort(dateRange, data, field, sortBy, "month")
}

func InsertEmptyDayDataAndSort(dateRange [2]time.Time, data interface{}, field, sortBy string) (dataList []map[string]interface{}, err error) {
	return insertEmptyDataAndSort(dateRange, data, field, sortBy, "day")
}

func insertEmptyDataAndSort(dateRange [2]time.Time, data interface{}, field, sortBy, key string) (dataList []map[string]interface{}, err error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return
	}
	err = json.Unmarshal(dataJson, &dataList)
	if err != nil {
		return
	}

	dataList = InsertEmptyDataByTimeKey(dateRange, dataList, field, key)
	SliceMapSortByTimeField(dataList, field, sortBy, key)

	return
}
