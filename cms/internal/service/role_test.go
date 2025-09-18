package service_test

import (
	"fmt"
	"net/http"
	"testing"

	"data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/internal/service"
	"data_backend/pkg/database"

	"github.com/gin-gonic/gin"
)

func TestRoleCreate(t *testing.T) {
	svc := service.NewRoleSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	err := svc.Create(&form.RoleCreateRequest{
		Name:             "Test2",
		PermissionIDList: []uint32{1, 3, 5, 6},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func TestRoleList(t *testing.T) {
	req, _ := http.NewRequest("GET", "test?page_size=10", nil)
	ctx := &gin.Context{
		Request: req,
	}
	svc := service.NewRoleSvc(ctx, db, rdb, log, global.NewAlarm)
	data, count, err := svc.List(&form.RoleListRequest{
		Name: "es",
	})
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
	svc := service.NewRoleSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	data, err := svc.All(database.QueryWhereGroup{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestRoleUpdate(t *testing.T) {
	ctx := &gin.Context{}
	svc := service.NewRoleSvc(ctx, db, rdb, log, global.NewAlarm)
	permList, err := svc.Update(5, &form.RoleUpdateRequest{Name: "Test2", PermissionIDList: []uint32{1, 2, 3}})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(permList)

}

func TestRoleGetPermNameList(t *testing.T) {
	svc := service.NewRoleSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	permList, err := svc.GetPermNameList([]uint32{1})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(permList)
}
