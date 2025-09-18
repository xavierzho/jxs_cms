package dao_test

import (
	"fmt"
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/inquire/item/dao"
	"data_backend/internal/app"
	"data_backend/pkg/database"
)

func TestItemLog(t *testing.T) {
	itemLogDao := dao.NewItemDao(local.CenterDB, local.Logger)
	cTime := time.Now()
	end := time.Date(cTime.Year(), cTime.Month(), cTime.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, 1)
	start := end.AddDate(0, 0, -10)

	summary, data, err := itemLogDao.GetLog(
		[2]time.Time{start, end}, nil,
		dao.AllRequestParamsGroup{
			BetFlag:    false,
			MarketFlag: false,
			ActivityFlag: dao.ActivityFlag{
				Flag:      true,
				CostAward: true,
				CostRank:  true,
			},
			UsersParams: database.QueryWhereGroup{
				// {Prefix: "u.id = ?", Value: []any{26}},
			},
			GachaParams: database.QueryWhereGroup{
				// {Prefix: "gur.gacha_name = ?", Value: []any{"扭蛋机"}},
			},
			ItemParams: database.QueryWhereGroup{
				// {Prefix: "inner_price between ? and ?", Value: []any{0, 1000000}},
				// {Prefix: "update_amount between ? and ?", Value: []any{-100000, 100000}},
			},
			AmountParams: database.QueryWhereGroup{
				// {Prefix: "update_amount between ? and ?", Value: []any{-100000, 100000}},
			},
		},
		&app.Pager{
			PageSize: 50,
		},
	)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("summary: ", summary)

	fmt.Println(len(data))

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
