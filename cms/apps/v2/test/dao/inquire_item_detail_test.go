package dao_test

import (
	"testing"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/item/dao"
	"data_backend/pkg/database"
)

func TestItemDetailBet(t *testing.T) {
	itemLogDao := dao.NewItemDao(local.CenterDB, local.Logger)

	itemLogDao.GetDetailBet(
		database.QueryWhereGroup{
			{Prefix: "t.id = ?", Value: []any{30}},
			{Prefix: "j.level_type = ?", Value: []any{1}},
		},
	)

}

func TestItemDetailMarket(t *testing.T) {
	itemLogDao := dao.NewItemDao(local.CenterDB, local.Logger)

	itemLogDao.GetDetailMarket(
		database.QueryWhereGroup{{Prefix: "t.id = ?", Value: []any{30}}},
	)

}
