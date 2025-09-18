package dao

import (
	"context"
	"fmt"

	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type User struct {
	UserID   int64  `gorm:"column:id; type:bigint" json:"user_id"`
	UserName string `gorm:"column:nickname; type:varchar(64)" json:"user_name"`
	Tel      string `gorm:"column:phone_num_md5; type:varchar(64)" json:"tel"`
}

func (User) TableName() string {
	return "users"
}

type UserDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewUserDao(center *gorm.DB, log *logger.Logger) *UserDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".UserDao")))
	return &UserDao{
		center: center,
		logger: log,
	}
}

func (d *UserDao) First(queryParams database.QueryWhereGroup) (users *User, err error) {
	err = d.center.
		Table(fmt.Sprintf("%s as u", User{}.TableName())).
		Scopes(database.ScopeQuery(queryParams)).
		First(&users).Error
	if err != nil {
		d.logger.Errorf("First err: %v", err)
		return nil, err
	}

	return
}
