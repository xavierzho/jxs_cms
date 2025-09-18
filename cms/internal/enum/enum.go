package enum

import (
	"context"
	"fmt"

	"data_backend/internal/global"
	"data_backend/pkg/convert"
	"data_backend/pkg/i18n"
)

// 按照 string value 值或 index 来获取 Enum 值
func convertEnum(arr []string, value interface{}) (interface{}, error) {
	switch value := value.(type) {
	case string:
		arrMap := make(map[string]interface{})
		for i := 0; i < len(arr); i++ {
			arrMap[arr[i]] = i
		}
		if index, ok := arrMap[value]; ok {
			return index, nil
		} else {
			return nil, fmt.Errorf("convertEnum: key cannot be found")
		}
	case int:
		if value < 0 {
			return nil, fmt.Errorf("convertEnum: index cannot less than zero")
		}
		if value > len(arr)-1 {
			return "-", nil
		}
		return arr[value], nil
	default:
		return nil, fmt.Errorf("convertEnum: value type not allow")
	}
}

func convertEnumT(ctx context.Context, key string, value interface{}) (interface{}, error) {
	translateContent := key + i18n.NESTED_SEPARATOR + convert.GetString(value)
	result, ok := global.I18n.ShouldT(ctx, "enum", translateContent)
	if !ok {
		return "-", nil // 不返回错误 而是 返回一个占位符
	} else {
		return result, nil
	}

}
