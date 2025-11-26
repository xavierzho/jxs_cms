package dao

import (
	"data_backend/pkg/logger"
	"data_backend/pkg/util"

	"gorm.io/gorm"
)

// InitFirstUser 初始化第一个用户
func InitFirstUser(db *gorm.DB, log *logger.Logger) (err error) {
	errPtr := &err
	tx := db.Begin()
	defer DeferTransaction(tx, log, &errPtr, "InitFirstUser")

	// 先迁移其他表后再进行初始化
	err = tx.AutoMigrate(&Role{}, &User{})
	if err != nil {
		return err
	}

	// 获取全部权限
	var permissionList []*Permission
	if err = tx.Find(&permissionList).Error; err != nil {
		return err
	}

	role := &Role{
		Name:       "Admin",
		Permission: permissionList,
	}
	if err = tx.Create(role).Error; err != nil {
		return err
	}
	if err = tx.Model(role).Association("Permission").Replace(permissionList); err != nil {
		return err
	}

	user := &User{
		UserName: "admin",
		Name:     "Admin",
		Email:    "Admin@demo.com",
		Password: util.GeneratePassword("123123"),
		Role:     []*Role{role},
	}
	if err = tx.Create(user).Error; err != nil {
		return err
	}
	if err = tx.Model(user).Association("Role").Replace(user.Role); err != nil {
		return err
	}

	return nil
}
