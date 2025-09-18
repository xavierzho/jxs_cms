package dao

import (
	cDao "data_backend/apps/v2/internal/common/dao"
	"data_backend/pkg/database"

	"gorm.io/gorm"
)

type GachaDetail struct {
	cDao.Item
	BetNums   int `gorm:"column:bet_nums; type:int" json:"bet_nums"`
	TotalNums int `gorm:"column:total_nums; type:int" json:"total_nums"`
}

func (d *GachaDao) GetDetail(queryParams database.QueryWhereGroup) (data []*GachaDetail, err error) {
	err = d.center.
		Select("*").
		Table("(? union all ?) t", d.getDetailNormalDB(queryParams), d.getDetailExtraDB(queryParams)).
		Order("level_index, inner_price desc, total_nums desc, bet_nums desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("GetDetail: %v", err)
		return nil, err
	}

	return
}

func (d *GachaDao) getDetailNormalDB(queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Select(
			"gba.item_id", "i.name as item_name", "gba.level_index", "gl.level_name", "i.cover_thumb",
			"i.show_price", "i.inner_price", "i.recycling_price",
			"sum(gba.total_nums-gba.left_nums) as bet_nums", "sum(gba.total_nums) as total_nums",
		).
		Table("gacha_box gb, gacha_box_award gba, item i, gacha_level gl").
		Scopes(database.ScopeQuery(queryParams)).
		Where("gb.gacha_id = gba.gacha_id").
		Where("gb.box_index = gba.box_index").
		Where("gba.item_id = i.id").
		Where("gba.level_index = gl.level_index").
		Group("gba.item_id, i.name, gba.level_index, gl.level_name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price")
}

func (d *GachaDao) getDetailExtraDB(queryParams database.QueryWhereGroup) *gorm.DB {
	return d.center.
		Select(
			"ga.item_id", "i.name as item_name", "ga.level_index", "gl.level_name", "i.cover_thumb",
			"i.show_price", "i.inner_price", "i.recycling_price",
			"sum(case gb.state when 2 then ga.total_nums else 0 end) as bet_nums", "sum(ga.total_nums) as total_nums",
		).
		Table("gacha_box gb, gacha_award ga, item i, gacha_level gl").
		Scopes(database.ScopeQuery(queryParams)).
		Where("gb.gacha_id = ga.gacha_id").
		Where("ga.level_type <> 1").
		Where("ga.item_id = i.id").
		Where("ga.level_index = gl.level_index").
		Where("ga.deleted_at is null").
		Group("ga.item_id, i.name, ga.level_index, gl.level_name, i.cover_thumb, i.show_price, i.inner_price, i.recycling_price")
}
