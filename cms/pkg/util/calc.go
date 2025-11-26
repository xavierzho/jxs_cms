package util

import (
	"fmt"

	"data_backend/pkg/convert"

	"github.com/shopspring/decimal"
)

var DECIMAL_THOUSAND = decimal.NewFromInt(1000)
var DECIMAL_HUNDRED = decimal.NewFromInt(100)
var amount_precision = DECIMAL_HUNDRED

func SetPrecision(precision decimal.Decimal) {
	amount_precision = precision
}

// 转换金额 -> 金额(元)
func ConvertAmount2Decimal(amount interface{}) decimal.Decimal {
	return convert.GetDecimal(amount).Div(amount_precision)
}

// 转换金额 -> 金额(元)
func ConvertAmount2Float64(amount interface{}) float64 {
	amountFloat64, _ := convert.GetDecimal(amount).Div(amount_precision).Float64()
	return amountFloat64
}

// 转换金额 -> 金额(元)
func ConvertAmount2Int64(amount interface{}) int64 {
	amountInt := convert.GetDecimal(amount).Div(amount_precision).IntPart()
	return amountInt
}

// 转换金额(元) -> 金额
func ReconvertAmount2Decimal(amount interface{}) decimal.Decimal {
	return convert.GetDecimal(amount).Mul(amount_precision)
}

// 两数相加
func Add2Float64(a, b interface{}) float64 {
	amount, _ := convert.GetDecimal(a).Add(convert.GetDecimal(b)).Float64()
	return amount
}

// 两数相加
func Add2Decimal(a, b interface{}) decimal.Decimal {
	return convert.GetDecimal(a).Add(convert.GetDecimal(b))
}

// 两数相减 a-b
func Sub2Float64(a, b interface{}) float64 {
	amount, _ := convert.GetDecimal(a).Sub(convert.GetDecimal(b)).Float64()
	return amount
}

// 两数相减 a-b
func Sub2Decimal(a, b interface{}) decimal.Decimal {
	return convert.GetDecimal(a).Sub(convert.GetDecimal(b))
}

// 两数相乘
func Multiply2Float64(a, b interface{}) float64 {
	amount, _ := convert.GetDecimal(a).Mul(convert.GetDecimal(b)).Float64()
	return amount
}

// 两数相乘
func Multiply2Decimal(a, b interface{}) decimal.Decimal {
	return convert.GetDecimal(a).Mul(convert.GetDecimal(b))
}

// 两数相除 a/b
func Divide2Float64(a, b interface{}) float64 {
	amount, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Float64()
	return amount
}

// 两数相除 a/b
func Divide2Decimal(a, b interface{}) decimal.Decimal {
	return convert.GetDecimal(a).Div(convert.GetDecimal(b))
}

// 两数相除转为 字符串百分比(2位小数) a/b -> .2f%
func Ratio2String(a, b interface{}) string {
	rate, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2).Float64()
	return fmt.Sprintf("%.2f%%", rate)
}

// 两数相除转为 百分比数值(2位小数) (a/b)*100 -> .2f
func Ratio2Float64(a, b interface{}) float64 {
	rate, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2).Float64()
	return rate
}

// 两数相除转为 百分比数值(2位小数) (a/b)*100 -> .2f
func Ratio2Decimal(a, b interface{}) decimal.Decimal {
	return convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2)
}

// 两数相除 a/b
func SaveDivide2Float64(a, b interface{}) float64 {
	if convert.GetDecimal(b).IsZero() {
		return 0
	}
	amount, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Float64()
	return amount
}

// 两数相除 a/b
func SaveDivide2Decimal(a, b interface{}) decimal.Decimal {
	if convert.GetDecimal(b).IsZero() {
		return decimal.NewFromInt(0)
	}
	return convert.GetDecimal(a).Div(convert.GetDecimal(b))
}

// 两数相除转为 字符串百分比(2位小数) a/b -> .2f%
func SaveRatio2String(a, b interface{}) string {
	if convert.GetDecimal(b).IsZero() {
		return "0%"
	}
	rate, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2).Float64()
	return fmt.Sprintf("%.2f%%", rate)
}

// 两数相除转为 百分比数值(2位小数) (a/b)*100 -> .2f
func SaveRatio2Float64(a, b interface{}) float64 {
	if convert.GetDecimal(b).IsZero() {
		return 0
	}
	rate, _ := convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2).Float64()
	return rate
}

// 两数相除转为 百分比数值(2位小数) (a/b)*100 -> .2f
func SaveRatio2Decimal(a, b interface{}) decimal.Decimal {
	if convert.GetDecimal(b).IsZero() {
		return decimal.NewFromInt(0)
	}
	return convert.GetDecimal(a).Div(convert.GetDecimal(b)).Mul(decimal.NewFromInt(100)).Round(2)
}
