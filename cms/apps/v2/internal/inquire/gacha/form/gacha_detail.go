package form

import (
	"fmt"
	"strconv"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/gacha/dao"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"
)

type DetailRequest struct {
	GachaID  int64 `form:"gacha_id" binding:"required"`
	BoxOutNo int64 `form:"box_out_no"`
}

func (q *DetailRequest) Parse() (queryParams database.QueryWhereGroup, err error) {
	queryParams = append(queryParams, database.QueryWhere{
		Prefix: "gb.gacha_id = ?",
		Value:  []any{q.GachaID},
	})

	if q.BoxOutNo != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "gb.box_out_no = ?",
			Value:  []any{q.BoxOutNo},
		})
	}

	return
}

type GachaDetail struct {
	cForm.Item
	BetNums   int     `json:"bet_nums"`
	TotalNums int     `json:"total_nums"`
	BetRate   float64 `json:"bet_rate"`
}

func FormatDetail(data []*dao.GachaDetail) (result []*GachaDetail) {
	for _, item := range data {
		result = append(result, &GachaDetail{
			Item: cForm.Item{
				ItemID:         strconv.FormatInt(item.ItemID, 10),
				ItemName:       item.ItemName,
				LevelName:      item.LevelName,
				CoverThumb:     item.CoverThumb,
				ShowPrice:      util.ConvertAmount2Decimal(item.ShowPrice),
				InnerPrice:     util.ConvertAmount2Decimal(item.InnerPrice),
				RecyclingPrice: util.ConvertAmount2Decimal(item.RecyclingPrice),
			},
			BetNums:   item.BetNums,
			TotalNums: item.TotalNums,
			BetRate:   util.SaveRatio2Float64(item.BetNums, item.TotalNums),
		})
	}

	return
}

func Format2Excel(params *DetailRequest, _data []*dao.GachaDetail) (excelModel *excel.Excel[*GachaDetail], err error) {
	data := FormatDetail(_data)

	reflectMap := map[string]func(source *GachaDetail) any{
		"物品id":  func(source *GachaDetail) any { return source.ItemID },
		"物品名称":  func(source *GachaDetail) any { return source.ItemName },
		"等级":    func(source *GachaDetail) any { return source.LevelName },
		"封面缩略图": func(source *GachaDetail) any { return source.CoverThumb },
		"展示价":   func(source *GachaDetail) any { return source.ShowPrice },
		"成本价":   func(source *GachaDetail) any { return source.InnerPrice },
		"回收价":   func(source *GachaDetail) any { return source.RecyclingPrice },
		"抽数":    func(source *GachaDetail) any { return source.BetNums },
		"总抽数":   func(source *GachaDetail) any { return source.TotalNums },
		"已抽比例":  func(source *GachaDetail) any { return source.BetRate },
	}

	fileName := ""
	if params.BoxOutNo != 0 {
		fileName = fmt.Sprintf("gacha_detail_%d-%d", params.GachaID, params.BoxOutNo)
	} else {
		fileName = fmt.Sprintf("gacha_detail_%d", params.GachaID)
	}

	excelModel = &excel.Excel[*GachaDetail]{
		FileName:   fileName,
		SheetNames: []string{"奖箱详情"},
		SheetNameWithHead: map[string][]string{
			"奖箱详情": {
				"物品id", "物品名称", "等级", "封面缩略图",
				"展示价", "成本价", "回收价", "抽数", "总抽数", "已抽比例",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*GachaDetail]{
			"奖箱详情": data,
		},
		ReflectMap: map[string]excel.RowReflect[*GachaDetail]{
			"奖箱详情": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
