package util

import (
	"sort"
	"time"

	"data_backend/pkg"
)

// 将获得的日期 格式化为顺序参数
// 例子: ["2022-02-10","2022-02-11","2022-02-13"] => [["2022-02-10","2022-02-11"],["2022-02-13","2022-02-13"]]
func SliceDateToSliceDateRange(dates []string) [][2]string {
	if len(dates) == 0 {
		return nil
	}
	res := make([][2]string, 0, len(dates))
	temp := ""
	firstParam := ""
	sliceSortByDate(dates, "asc")
	for i := 0; i < len(dates); i++ {
		if temp == "" {
			temp = dates[i]
			firstParam = dates[i]
			continue
		}
		if handle, err := IsSequences(temp, dates[i]); err == nil && handle {
			temp = dates[i]
			continue
		} else {
			res = append(res, [2]string{firstParam, temp})
			firstParam = dates[i]
			temp = dates[i]
			continue
		}
	}
	if firstParam != "" {
		res = append(res, [2]string{firstParam, temp})
	}
	return res
}

// 判断两个参数是否为连续的日期
func IsSequences(first, second string) (bool, error) {
	firstTime, err := time.Parse(pkg.DATE_FORMAT, first)
	if err != nil {
		return false, err
	}
	secondTime, err := time.Parse(pkg.DATE_FORMAT, second)
	if err != nil {
		return false, err
	}
	return firstTime.AddDate(0, 0, 1).Format(pkg.DATE_FORMAT) == secondTime.Format(pkg.DATE_FORMAT), nil
}

func sliceSortByDate(dateArr []string, sortBy string) {
	if len(dateArr) == 0 {
		return
	}

	if _, err := time.Parse(pkg.DATE_FORMAT, dateArr[0]); err != nil {
		return
	}

	sort.Slice(dateArr, func(i, j int) bool {
		iTime, _ := time.Parse(pkg.DATE_FORMAT, dateArr[i])
		jTime, _ := time.Parse(pkg.DATE_FORMAT, dateArr[j])
		if sortBy == "desc" {
			return iTime.After(jTime)
		} else {
			return iTime.Before(jTime)
		}
	})
}

// 并集
func IntSliceUnionSet(slice1, slice2 []int) []int {
	intMap := make(map[int]struct{})
	for _, item := range slice1 {
		intMap[item] = struct{}{}
	}

	for _, item := range slice2 {
		intMap[item] = struct{}{}
	}

	var result []int
	for key := range intMap {
		result = append(result, key)
	}

	return result
}

// 并集
func Uint32SliceUnionSet(slice1, slice2 []uint32) []uint32 {
	intMap := make(map[uint32]struct{})
	for _, item := range slice1 {
		intMap[item] = struct{}{}
	}

	for _, item := range slice2 {
		intMap[item] = struct{}{}
	}

	var result []uint32
	for key := range intMap {
		result = append(result, key)
	}

	return result
}

// 并集
func StringSliceUnionSet(slice1, slice2 []string) []string {
	sliceMap := make(map[string]struct{})
	for _, item := range slice1 {
		sliceMap[item] = struct{}{}
	}

	for _, item := range slice2 {
		sliceMap[item] = struct{}{}
	}

	var result []string
	for key := range sliceMap {
		result = append(result, key)
	}

	return result
}

// 交集
func IntSliceIntersectionSet(slice1, slice2 []int) []int {
	sliceMap1 := make(map[int]struct{})
	sliceMap2 := make(map[int]struct{})
	var result []int

	for _, item := range slice1 {
		sliceMap1[item] = struct{}{}
	}

	for _, item := range slice2 {
		sliceMap2[item] = struct{}{}
	}

	for key := range sliceMap1 {
		if _, ok := sliceMap2[key]; ok {
			result = append(result, key)
		}
	}

	return result
}

// 交集
func Uint32SliceIntersectionSet(slice1, slice2 []uint32) []uint32 {
	sliceMap1 := make(map[uint32]struct{})
	sliceMap2 := make(map[uint32]struct{})
	var result []uint32

	for _, item := range slice1 {
		sliceMap1[item] = struct{}{}
	}

	for _, item := range slice2 {
		sliceMap2[item] = struct{}{}
	}

	for key := range sliceMap1 {
		if _, ok := sliceMap2[key]; ok {
			result = append(result, key)
		}
	}

	return result
}

// 交集
func StringSliceIntersectionSet(slice1, slice2 []string) []string {
	sliceMap1 := make(map[string]struct{})
	sliceMap2 := make(map[string]struct{})
	var result []string

	for _, item := range slice1 {
		sliceMap1[item] = struct{}{}
	}

	for _, item := range slice2 {
		sliceMap2[item] = struct{}{}
	}

	for key := range sliceMap1 {
		if _, ok := sliceMap2[key]; ok {
			result = append(result, key)
		}
	}

	return result
}

// 差集 slice1-slice2
func IntSliceDifferenceSet(slice1, slice2 []int) []int {
	intMap := make(map[int]struct{})
	var result []int

	for _, item := range slice2 {
		intMap[item] = struct{}{}
	}

	for _, item := range slice1 {
		if _, ok := intMap[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// 差集 slice1-slice2
func Uint32SliceDifferenceSet(slice1, slice2 []uint32) []uint32 {
	intMap := make(map[uint32]struct{})
	var result []uint32

	for _, item := range slice2 {
		intMap[item] = struct{}{}
	}

	for _, item := range slice1 {
		if _, ok := intMap[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// 差集 slice1-slice2
func StringSliceDifferenceSet(slice1, slice2 []string) []string {
	slice2Set := make(map[string]struct{})
	var result []string

	for _, item := range slice2 {
		slice2Set[item] = struct{}{}
	}

	for _, item := range slice1 {
		if _, ok := slice2Set[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// 对称差集 (slice1-slice2) ∪ (slice2-slice1)
func Uint32SliceSymmetricDifferenceSet(slice1, slice2 []uint32) []uint32 {
	diff1 := Uint32SliceDifferenceSet(slice1, slice2)
	diff2 := Uint32SliceDifferenceSet(slice2, slice1)
	result := Uint32SliceUnionSet(diff1, diff2)

	return result
}

// Contain slice in mainSlice
func IsContainAllStringSlice(slice, mainSlice []string) bool {
	if len(slice) == 0 {
		return true
	}
	return len(StringSliceDifferenceSet(slice, mainSlice)) == 0
}

// Contain slice any in mainSlice
func IsContainAnyStringSlice(slice, mainSlice []string) bool {
	if len(slice) == 0 {
		return true
	}
	return len(StringSliceIntersectionSet(slice, mainSlice)) > 0
}

// 确认是否权限都通过
func PermissionCheckAll(confirmPermission []string, permissions []string) bool {
	permMap := make(map[string]struct{})

	for i := 0; i < len(permissions); i++ {
		permMap[permissions[i]] = struct{}{}
	}
	// 循环判断该确认的权限
	for _, p := range confirmPermission {
		if p == "" {
			continue
		}
		if _, ok := permMap[p]; !ok {
			return false
		}
	}
	return true
}

// 确认是否权限有至少一个通过校验
func PermissionCheckOr(confirmPermission []string, permissions []string) bool {
	permMap := make(map[string]struct{})

	for i := 0; i < len(permissions); i++ {
		permMap[permissions[i]] = struct{}{}
	}
	// 循环判断该确认的权限
	for _, p := range confirmPermission {
		if p == "" {
			continue
		}
		if _, ok := permMap[p]; ok {
			return true
		}
	}
	return false
}

// distinct
func DistinctStrings(slice []string) (result []string) {
	sliceSet := make(map[string]struct{})

	for _, item := range slice {
		if _, ok := sliceSet[item]; !ok {
			sliceSet[item] = struct{}{}
			result = append(result, item)
		}
	}

	return
}

// 按照给定字段是否参与"group by"的组合情况 添加到groupByFieldsList中
// 例如：有两个字段a, b 那么可能的组合有(下面的other 表示实际存在的其他字段)：
// 1. group by a, other 2. group by a, b, other 3. group by b, other 4. group by other
// 则对应的列表元素为: [a], [a, b], [b], []
func PerCombination(fields []string) [][]string {
	if len(fields) == 0 {
		return [][]string{}
	}
	groupByFieldsList := [][]string{{fields[0]}, {}} // 需要添加一个起始值，不然后续的for ind := range group_2 无法开始

	for _, field := range fields[1:] {
		groupByFieldsList = addCombinationType(field, groupByFieldsList)
	}

	return groupByFieldsList
}

// 用于给多组数据的sql 生成group by 字段列表
// 将新的字段 field 按照使用和不使用的情况添加到列表中的每一项中去
// 即 复制一份列表， 其中一份保留原样，另一份末尾添加该 field
func addCombinationType(field string, groupByFieldsList [][]string) [][]string {
	group_1 := make([][]string, 0, len(groupByFieldsList))
	group_2 := make([][]string, 0, len(groupByFieldsList))
	for _, item := range groupByFieldsList {
		elem_1 := make([]string, len(item))
		elem_2 := make([]string, len(item))
		copy(elem_1, item)
		copy(elem_2, item)
		group_1 = append(group_1, elem_1)
		group_2 = append(group_2, elem_2)
	}

	for ind := range group_2 {
		group_2[ind] = append(group_2[ind], field)
	}

	return append(group_1, group_2...)
}

// 返回时间范围的连续日期序列
func TimeRange2StrDateRangeSlice(dateRange [2]time.Time, asc bool) []string {
	var result []string
	startTime := dateRange[0]
	endTime := dateRange[1]

	if startTime.After(endTime) {
		return nil
	}

	if asc {
		for cDate := startTime; cDate.Before(endTime) || cDate.Equal(endTime); cDate = cDate.AddDate(0, 0, 1) {
			result = append(result, cDate.Format(pkg.DATE_FORMAT))
		}
	} else {
		for cDate := endTime; cDate.After(startTime) || cDate.Equal(startTime); cDate = cDate.AddDate(0, 0, -1) {
			result = append(result, cDate.Format(pkg.DATE_FORMAT))
		}
	}

	return result
}

// 对于切分进行分段操作
func SplitToOperate(totalLen int, batchSize int, batchFunc func(int, int) error) (err error) {
	for i := 0; i < (totalLen/batchSize)+1; i++ {
		first := i * batchSize
		end := (i + 1) * batchSize
		if end > totalLen {
			end = totalLen
		}
		if first == end {
			continue
		}
		err = batchFunc(first, end)
		if err != nil {
			return err
		}
	}
	return nil
}
