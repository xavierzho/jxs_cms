package service_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"data_backend/internal/app"
	"data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/internal/service"

	"github.com/gin-gonic/gin"
)

func TestUserLogin(t *testing.T) {
	svc := service.NewUserSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	data, err := svc.Login(&form.LoginRequest{
		UserName: "hjw",
		Password: "hjw123456",
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v", data)

	<-time.After(time.Second * 10)
}

func TestUserCreate(t *testing.T) {
	svc := service.NewUserSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	svc.Create(&form.UserCreateRequest{
		UserName:   "test2",
		Name:       "test2",
		Email:      "test2@demo.com",
		Password:   "123123",
		RoleIDList: []uint32{4, 5},
	})

}

func TestUserList(t *testing.T) {
	req, _ := http.NewRequest("GET", "test?order_field=role&order=desc", nil)
	ctx := &gin.Context{
		Request: req,
	}
	fmt.Println(app.GetOrderBy(ctx))
	svc := service.NewUserSvc(ctx, db, rdb, log, global.NewAlarm)
	data, count, err := svc.List(&form.UserListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(count)
	for _, item := range data {
		fmt.Printf("%+v\n", item)
	}
}

func TestUserUpdate(t *testing.T) {
	svc := service.NewUserSvc(&gin.Context{}, db, rdb, log, global.NewAlarm)
	err := svc.Update(3, &form.UserUpdateRequest{
		Name:       "tttt",
		Email:      "t23@demo.com",
		Password:   "111",
		IsLock:     false,
		RoleIDList: []uint32{4, 5},
	})

	if err != nil {
		fmt.Printf("%#v", err)
	}

	<-time.After(time.Second * 10)
}

func TestUserUpdateSelf(t *testing.T) {
	ctx := &gin.Context{}
	svc := service.NewUserSvc(ctx, db, rdb, log, global.NewAlarm)
	data, err := svc.UpdateSelf(&form.UserUpdateSelfRequest{
		Name:        "etwsatsfda",
		Email:       "sadf@qq.com",
		NewPassword: "123",
	})
	if err != nil {
		fmt.Printf("%#v", err)
		return
	}
	fmt.Println(data)
}

func TestCanShowPhoneNum(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set(app.USER_ID_KEY, "1")
	svc := service.NewUserSvc(ctx, db, rdb, log, global.NewAlarm)
	fmt.Println(svc.CanShowPhoneNum())
}

func TestPagePermission(t *testing.T) {
	ctx := &gin.Context{}
	ctx.Set(app.USER_ID_KEY, "1")
	svc := service.NewUserSvc(ctx, db, rdb, log, global.NewAlarm)
	fmt.Println(svc.PagePermission())
}
