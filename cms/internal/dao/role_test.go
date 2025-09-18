package dao_test

import (
	"fmt"
	"testing"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	"data_backend/pkg/database"
)

func TestRoleList(t *testing.T) {
	d := dao.NewRoleDao(db, log)
	data, count, err := d.List(database.QueryWhereGroup{{Prefix: "id", Value: []interface{}{[]int{1, 2, 3}}}}, app.Pager{Page: 1, PageSize: 10})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(count)
	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestRoleAll(t *testing.T) {
	d := dao.NewRoleDao(db, log)
	data, err := d.All(database.QueryWhereGroup{
		{Prefix: "id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestRoleGetPermList(t *testing.T) {
	d := dao.NewRoleDao(db, log)
	data, err := d.GetPermNameList(database.QueryWhereGroup{
		{Prefix: "r.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	data2, err := d.GetPermIDList(database.QueryWhereGroup{
		{Prefix: "r.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data2)
}

func TestRoleUpdate(t *testing.T) {
	d := dao.NewRoleDao(db, log)

	data := &dao.Role{
		Model: dao.Model{ID: 2},
		Name:  "aaaa",
		Permission: []*dao.Permission{
			{ID: 1},
		},
	}

	if err := d.Update(data); err != nil {
		fmt.Println(err)
	}

}

func TestRoleUpdateAndAssociationReplace(t *testing.T) {
	d := dao.NewRoleDao(db, log)

	data := &dao.Role{
		Model: dao.Model{ID: 2},
		Name:  "aaaa",
		Permission: []*dao.Permission{
			{ID: 1},
		},
	}

	if err := d.UpdateAndAssociationReplace(data); err != nil {
		fmt.Println(err)
	}

}

func TestRoleOption(t *testing.T) {
	d := dao.NewRoleDao(db, log)
	data, err := d.Options()
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}
