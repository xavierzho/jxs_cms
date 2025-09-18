package dao

import (
	"context"
	"fmt"

	"data_backend/internal/app"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

const ROLE_PERM_R_KEY = "RolePerm" // redis 缓存中角色与对应的权限

type Role struct {
	Model

	Name string `gorm:"column:name; type:varchar(100); index:un_name,unique" json:"name"`

	Permission []*Permission `gorm:"many2many:role_permission" json:"permission,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

func (r Role) RolePermRKey() string {
	return fmt.Sprintf("%s:%d", ROLE_PERM_R_KEY, r.ID)
}

type RoleDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewRoleDao(engine *gorm.DB, log *logger.Logger) *RoleDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".RoleDao")))

	return &RoleDao{
		engine: engine,
		logger: log,
	}
}

func (d *RoleDao) Create(data *Role) (err error) {
	errPtr := &err
	tx := d.engine.Begin()
	defer DeferTransaction(tx, d.logger, &errPtr, "Create")

	if err = tx.Create(data).Error; err != nil {
		d.logger.Errorf("Create: %v", err)
		return err
	}

	if err = tx.Model(data).Association("Permission").Replace(data.Permission); err != nil {
		d.logger.Errorf("Create Association.Replace: %v", err)
		return err
	}

	return nil
}

func (d *RoleDao) First(queryParams database.QueryWhereGroup) (data *Role, err error) {
	err = d.engine.Model(data).Preload("Permission").Scopes(database.ScopeQuery(queryParams)).First(&data).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		d.logger.Errorf("First: %v", err)
		return nil, err
	}

	return data, err
}

func (d *RoleDao) List(queryParams database.QueryWhereGroup, pager app.Pager) (data []*Role, count int64, err error) {
	err = d.engine.
		Model(data).
		Preload("Permission").
		Scopes(database.ScopeQuery(queryParams)).
		Count(&count).
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List: %v", err)
		return nil, 0, err
	}

	return data, count, nil
}

func (d *RoleDao) All(queryParams database.QueryWhereGroup) (data []*Role, err error) {
	err = d.engine.Model(data).Preload("Permission").Scopes(database.ScopeQuery(queryParams)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *RoleDao) GetPermNameList(queryParams database.QueryWhereGroup) (permNameList []string, err error) {
	err = d.getPermListSql(queryParams).Order("p.name").Pluck("distinct p.name", &permNameList).Error
	if err != nil {
		d.logger.Errorf("GetPermNameList: %v", err)
		return nil, err
	}

	return permNameList, nil
}

func (d *RoleDao) GetPermIDList(queryParams database.QueryWhereGroup) (permIDList []uint32, err error) {
	err = d.getPermListSql(queryParams).Order("p.id").Pluck("distinct p.id", &permIDList).Error
	if err != nil {
		d.logger.Errorf("GetPermIDList: %v", err)
		return nil, err
	}

	return permIDList, nil
}

func (d *RoleDao) getPermListSql(queryParams database.QueryWhereGroup) (tx *gorm.DB) {
	return d.engine.
		Table("roles r").
		Joins("join role_permission rp on rp.role_id = r.id").
		Joins("join permissions p on rp.permission_id = p.id").
		Scopes(database.ScopeQuery(queryParams))
}

func (d *RoleDao) Update(data *Role) (err error) {
	data.Permission = nil // 为避免错误的更新数据, 当需要更新 Permission 则调用 UpdateAndAssociationReplace
	if err = d.engine.Updates(data).Error; err != nil {
		d.logger.Errorf("Update: %v", err)
		return err
	}

	return nil
}

func (d *RoleDao) UpdateAndAssociationReplace(data *Role) (err error) {
	errPtr := &err
	tx := d.engine.Begin()
	defer DeferTransaction(tx, d.logger, &errPtr, "UpdateAndAssociationReplace")

	if err = tx.Updates(data).Error; err != nil {
		d.logger.Errorf("Update: %v", err)
		return err
	}

	if err = tx.Model(data).Association("Permission").Replace(data.Permission); err != nil {
		d.logger.Errorf("Update Association.Replace: %v", err)
		return err
	}

	return nil
}

func (d *RoleDao) Options() (data []*Role, err error) {
	if err = d.engine.Order("id").Find(&data).Error; err != nil {
		d.logger.Errorf("Options: %v", err)
		return nil, err
	}

	return data, nil
}
