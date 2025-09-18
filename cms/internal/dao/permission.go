package dao

import (
	"context"

	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

const (
	PERMISSION_LANG_UPDATE         = "lang_update"
	PERMISSION_SHOW_SENSITIVE_INFO = "show_sensitive_info"
)

type Permission struct {
	ID          uint32 `gorm:"column:id; primary_key" json:"id,omitempty"`
	Name        string `gorm:"column:name; type:varchar(100); index:un_name,unique" json:"name"`
	DisplayName string `gorm:"column:display_name; type:varchar(100)" json:"display_name"`
	Description string `gorm:"column:description; type:varchar(100)" json:"description"`
	Pages       string `gorm:"column:pages" json:"pages"`
}

func (Permission) TableName() string {
	return "permissions"
}

type PermissionDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewPermissionDao(engine *gorm.DB, log *logger.Logger) *PermissionDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".PermissionDao")))

	return &PermissionDao{
		engine: engine,
		logger: log,
	}
}

func (d *PermissionDao) Create(data []*Permission) (err error) {
	if err = d.engine.Create(data).Error; err != nil {
		d.logger.Errorf("Create: %v", err)
		return err
	}

	return nil
}

func (d *PermissionDao) All(queryParams []database.QueryWhere) (data []*Permission, err error) {
	if err = d.engine.Scopes(database.ScopeQuery(queryParams)).Order("id").Find(&data).Error; err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *PermissionDao) Options() (data []*Permission, err error) {
	if err = d.engine.Order("name asc").Find(&data).Error; err != nil {
		d.logger.Errorf("Options: %v", err)
		return nil, err
	}

	return data, nil
}
