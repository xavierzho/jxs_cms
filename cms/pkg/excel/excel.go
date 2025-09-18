package excel

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"data_backend/pkg"

	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
)

func NewFile() *excelize.File {
	return excelize.NewFile()
}

type Excel[T any] struct {
	FileName string
	// 用于确定生成的sheet表的顺序
	SheetNames []string
	// map[sheetName][] 各个表的表头名称
	SheetNameWithHead map[string][]string
	DefaultColWidth   float64
	DefaultRowHeight  float64
	Excelize          *excelize.File
	Data              map[string]SheetData[T]
	ReflectMap        map[string]RowReflect[T]
}

type SheetName string
type SheetData[T any] []T

// 反射行的值
type RowReflect[T any] map[string]func(source T) interface{}

func (excelFile *Excel[T]) InitExcelFile() error {
	excel := excelize.NewFile()
	excelFile.Excelize = excel
	letterA := 'A'
	// 表
	for _, sheetName := range excelFile.SheetNames {
		// 设置工作表名或者创建新的工作表
		headNames := excelFile.SheetNameWithHead[sheetName]
		if excel.GetSheetName(0) == "Sheet1" {
			excel.SetSheetName("Sheet1", sheetName)
		} else {
			excel.NewSheet(sheetName)
		}

		// 设置长度和宽度
		err := excel.SetColWidth(sheetName, fmt.Sprintf("%c", letterA), fmt.Sprintf("%c", letterA+int32(len(headNames))-1), excelFile.DefaultColWidth)
		if err != nil {
			return err
		}
		err = excel.SetRowHeight(sheetName, 1, excelFile.DefaultRowHeight)
		if err != nil {
			return err
		}
		// 设置表头字段
		for j := 0; j < len(headNames); j++ {
			err = excel.SetCellValue(sheetName, fmt.Sprintf("%s1", GetLetterByIndex(j)), headNames[j])
			if err != nil {
				return err
			}
		}
		// 准备插入数据
		if _, ok := excelFile.Data[sheetName]; !ok {
			continue
		}
		if _, ok := excelFile.ReflectMap[sheetName]; !ok {
			continue
		}
		sheetData := excelFile.Data[sheetName]
		reflectMap := excelFile.ReflectMap[sheetName]
		decimalStyle, err := excel.NewStyle(&excelize.Style{NumFmt: 2})
		if err != nil {
			return err
		}
		// 行
		for i := 0; i < len(sheetData); i++ {
			rowNum := i + 2
			// 列
			for j := 0; j < len(headNames); j++ {
				if _, ok := reflectMap[headNames[j]]; !ok {
					continue
				}
				// 获取该row的值
				cellTag := fmt.Sprintf("%s%d", GetLetterByIndex(j), rowNum)
				value := reflectMap[headNames[j]](sheetData[i])
				switch value.(type) {
				case decimal.Decimal:
					f, _ := value.(decimal.Decimal).Float64()
					err = excel.SetCellValue(sheetName, cellTag, f)
					if err != nil {
						return err
					}
					err = excel.SetCellStyle(sheetName, cellTag, cellTag, decimalStyle)
					if err != nil {
						return err
					}
				case string:
					if v, err := decimal.NewFromString(value.(string)); err == nil && len(value.(string)) < 12 {
						if !v.IsInteger() {
							f, _ := v.Float64()
							err = excel.SetCellValue(sheetName, cellTag, f)
							if err != nil {
								return err
							}
						} else {
							err = excel.SetCellValue(sheetName, cellTag, v.IntPart())
							if err != nil {
								return err
							}
						}
					} else {
						err = excel.SetCellStr(sheetName, cellTag, value.(string))
						if err != nil {
							return err
						}
					}
				case int64:
					// 防止大的int64转科学计数法
					err = excel.SetCellStr(sheetName, cellTag, strconv.Itoa(int(value.(int64))))
					if err != nil {
						return err
					}
				case time.Time:
					err = excel.SetCellValue(sheetName, cellTag, value.(time.Time).Format(pkg.DATE_TIME_FORMAT))
					if err != nil {
						return err
					}
				default:
					err = excel.SetCellValue(sheetName, cellTag, value)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (excelFile *Excel[T]) Save(excelName string) (string, error) {
	if excelFile.Excelize == nil {
		return "", errors.New("excel uninitialized")
	}
	now := time.Now()
	path := fmt.Sprintf("./tmp/%s", now.Format(pkg.FILE_DATE_FORMAT))
	filePath := fmt.Sprintf("%s/%s_%d.xlsx", path, excelName, now.Unix())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}
	err := excelFile.Excelize.SaveAs(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func GetLetterByIndex(i int) string {
	if i > 676 {
		return "A"
	}
	letterA := 'A'
	letterStr := ""
	if (letterA+int32(i)-int32(65))/26 == int32(0) {
		letterStr = fmt.Sprintf("%c", letterA+int32(i))
	} else {
		letterStr = fmt.Sprintf("%c%c", (letterA+int32(i)-int32(65))/26-1+65, (letterA+int32(i)-int32(65))%26+65)
	}
	return letterStr
}
