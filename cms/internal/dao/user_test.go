package dao_test

import (
	"errors"
	"fmt"
	"testing"

	"data_backend/internal/app"
	"data_backend/internal/dao"
	"data_backend/pkg/database"

	"gorm.io/gorm"
)

func TestUserList(t *testing.T) {
	d := dao.NewUserDao(db, log)
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

func TestUserAll(t *testing.T) {
	d := dao.NewUserDao(db, log)
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

func TestUserGetRoleList(t *testing.T) {
	d := dao.NewUserDao(db, log)
	data, err := d.GetRoleNameList(database.QueryWhereGroup{
		{Prefix: "u.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	data2, err := d.GetRoleIDList(database.QueryWhereGroup{
		{Prefix: "u.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data2)
}

func TestUserGetPermList(t *testing.T) {
	d := dao.NewUserDao(db, log)
	data, err := d.GetPermNameList(database.QueryWhereGroup{
		{Prefix: "u.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data)

	data2, err := d.GetPermIDList(database.QueryWhereGroup{
		{Prefix: "u.id", Value: []interface{}{[]string{"1", "2", "3"}}},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(data2)
}

func TestUserUpdate(t *testing.T) {
	d := dao.NewUserDao(db, log)

	data := &dao.User{
		Model: dao.Model{ID: 2},
		Name:  "aaaa",
		Role: []*dao.Role{
			{Model: dao.Model{ID: 1}},
		},
	}

	if err := d.Update(data); err != nil {
		fmt.Println(err)
	}

}

func TestUserUpdateAndAssociationReplace(t *testing.T) {
	d := dao.NewUserDao(db, log)

	data := &dao.User{
		Model: dao.Model{ID: 2},
		Name:  "aaaa",
		Role: []*dao.Role{
			{Model: dao.Model{ID: 1}},
		},
	}

	if err := d.UpdateAndAssociationReplace(data); err != nil {
		fmt.Println(err)
	}

}

func TestUserOption(t *testing.T) {
	d := dao.NewUserDao(db, log)
	data, err := d.Options()
	if err != nil {
		fmt.Println(err)
	}

	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestFirst(t *testing.T) {
	d := dao.NewUserDao(db, log)
	user, err := d.First(database.QueryWhereGroup{
		{Prefix: "user_name", Value: []any{"hjw"}},
	})

	fmt.Println(user.IsAdmin())
	fmt.Println(err != gorm.ErrRecordNotFound)
	fmt.Println(errors.Is(err, gorm.ErrRecordNotFound))
}
