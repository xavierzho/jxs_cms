package dao

import (
	"fmt"
	"testing"
	"time"

	"data_backend/apps/v2/internal/common/local"
	"data_backend/apps/v2/internal/report/revenue/dao"
	"data_backend/pkg"
)

func TestRevenueBalanceGenerate(t *testing.T) {
	dao := dao.NewBalanceDao(local.CMSDB, local.CenterDB, local.Logger)
	data, err := dao.Generate()
	if err != nil {
		return
	}

	fmt.Printf("%v\n", data)
}

func TestRevenuePatingGenerate(t *testing.T) {
	dao := dao.NewPatingDao(local.CMSDB, local.CenterDB, local.Logger)
	cTime, _ := time.Parse(pkg.DATE_FORMAT, "2023-12-01")
	data, err := dao.Generate(cTime)
	if err != nil {
		return
	}

	fmt.Printf("%v\n", data)
}

func TestRevenueDrawGenerate(t *testing.T) {
	dao := dao.NewDrawDao(local.CMSDB, local.CenterDB, local.Logger)
	cTime, _ := time.Parse(pkg.DATE_FORMAT, "2023-12-01")
	data, err := dao.Generate(cTime)
	if err != nil {
		return
	}

	fmt.Printf("%v\n", data)
}

func TestRevenueWastageGenerate(t *testing.T) {
	dao := dao.NewWastageDao(local.CMSDB, local.CenterDB, local.Logger)
	cTime, _ := time.Parse(pkg.DATE_FORMAT, "2024-10-28")
	data, err := dao.Generate(cTime)
	if err != nil {
		return
	}

	fmt.Printf("%v\n", data)
}
