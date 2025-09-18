package dao_test

import (
	"testing"

	"data_backend/internal/app"
	"data_backend/internal/dao"
)

type Test struct {
	Name   string `gorm:"primary_key"`
	Value  uint
	Value2 uint
	Value3 uint
}

func (Test) TableName() string {
	return "test"
}

func TestDao(t *testing.T) {
	dao := dao.NewDao[*Test](db, log)

	a := &Test{Name: "test"}
	dao.Create(a)
	dao.Save(a)
	dao.List(nil, app.Pager{})
	dao.All(nil)
	dao.Update([]string{"value", "value2"}, []string{"value", "value3"}, a)
	dao.Update([]string{"*"}, []string{"value2"}, a)
}
