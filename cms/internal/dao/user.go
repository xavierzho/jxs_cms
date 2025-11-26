package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

const (
	UserRoleRKey = "UserRole"
	AdminRole    = 1
)

type User struct {
	Model

	UserName      string     `gorm:"column:user_name; type:varchar(100); index:un_username,unique" json:"user_name"`
	Name          string     `gorm:"column:name; type:varchar(50)" json:"name"` // 昵称
	Email         string     `gorm:"column:email; type:varchar(50); index:un_email,unique" json:"email"`
	IsLock        uint8      `gorm:"column:is_lock; type:TINYINT; default:0" json:"is_lock"`
	Password      string     `gorm:"column:password; type:varchar(200)" json:"-"`
	LastLogonTime *time.Time `gorm:"column:last_logon_time; type:datetime" json:"last_logon_time"`

	Role []*Role `gorm:"many2many:user_role" json:"role,omitempty"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) UserRoleRKey() string {
	return fmt.Sprintf("%s:%d", UserRoleRKey, u.ID)
}

func (u *User) IsAdmin() (isAdmin bool) {
	for _, role := range u.Role {
		if role.ID == AdminRole {
			isAdmin = true
			break
		}
	}

	return isAdmin
}

type UserDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewUserDao(engine *gorm.DB, log *logger.Logger) *UserDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".UserDao")))

	return &UserDao{
		engine: engine,
		logger: log,
	}
}

func (d *UserDao) Create(data *User) (err error) {
	errPtr := &err
	tx := d.engine.Begin()
	defer DeferTransaction(tx, d.logger, &errPtr, "Create")

	if err = tx.Create(data).Error; err != nil {
		d.logger.Errorf("Create: %v", err)
	}

	if err = tx.Model(data).Association("Role").Replace(data.Role); err != nil {
		d.logger.Errorf("Create Association.Replace: %v", err)
	}

	return nil
}

func (d *UserDao) First(queryParams database.QueryWhereGroup) (data *User, err error) {
	err = d.engine.Model(data).Preload("Role").Scopes(database.ScopeQuery(queryParams)).First(&data).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		d.logger.Errorf("First: %v", err)
		return nil, err
	}

	return data, err
}

func (d *UserDao) List(queryParams database.QueryWhereGroup, pager app.Pager) (data []*User, count int64, err error) {
	err = d.engine.
		Model(data).
		Preload("Role").
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

func (d *UserDao) All(queryParams database.QueryWhereGroup) (data []*User, err error) {
	err = d.engine.Model(data).Preload("Role").Scopes(database.ScopeQuery(queryParams)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *UserDao) GetRoleNameList(queryParams database.QueryWhereGroup) (roleNameList []string, err error) {
	err = d.getRoleListSql(queryParams).Order("r.name").Pluck("distinct r.name", &roleNameList).Error
	if err != nil {
		d.logger.Errorf("GetRoleNameList: %v", err)
		return nil, err
	}

	return roleNameList, nil
}

func (d *UserDao) GetRoleIDList(queryParams database.QueryWhereGroup) (roleIDList []uint32, err error) {
	err = d.getRoleListSql(queryParams).Order("r.id").Pluck("distinct r.id", &roleIDList).Error
	if err != nil {
		d.logger.Errorf("GetRoleIDList: %v", err)
		return nil, err
	}

	return roleIDList, nil
}

func (d *UserDao) getRoleListSql(queryParams database.QueryWhereGroup) (tx *gorm.DB) {
	return d.engine.
		Table("users u").
		Joins("join user_role ur on ur.user_id = u.id").
		Joins("join roles r on ur.role_id = r.id").
		Scopes(database.ScopeQuery(queryParams))
}

func (d *UserDao) GetPermNameList(queryParams database.QueryWhereGroup) (permNameList []string, err error) {
	err = d.getPermListSql(queryParams).Order("p.name").Pluck("distinct p.name", &permNameList).Error
	if err != nil {
		d.logger.Errorf("GetPermNameList: %v", err)
		return nil, err
	}

	return permNameList, nil
}

func (d *UserDao) GetPermIDList(queryParams database.QueryWhereGroup) (permIDList []uint32, err error) {
	err = d.getPermListSql(queryParams).Order("p.id").Pluck("distinct p.id", &permIDList).Error
	if err != nil {
		d.logger.Errorf("GetPermIDList: %v", err)
		return nil, err
	}

	return permIDList, nil
}

func (d *UserDao) getPermListSql(queryParams database.QueryWhereGroup) (tx *gorm.DB) {
	return d.engine.
		Table("users u").
		Joins("join user_role ur on ur.user_id = u.id").
		Joins("join role_permission rp on rp.role_id = ur.role_id").
		Joins("join permissions p on rp.permission_id = p.id").
		Scopes(database.ScopeQuery(queryParams))
}

func (d *UserDao) Update(data *User) (err error) {
	data.Role = nil // 为避免错误的更新数据, 当需要更新 Role 则调用 UpdateAndAssociationReplace
	if err = d.engine.Updates(data).Error; err != nil {
		d.logger.Errorf("Update: %v", err)
		return err
	}

	return nil
}

func (d *UserDao) UpdateAndAssociationReplace(data *User) (err error) {
	errPtr := &err
	tx := d.engine.Begin()
	defer DeferTransaction(tx, d.logger, &errPtr, "UpdateAndAssociationReplace")

	if err = tx.Updates(data).Error; err != nil {
		d.logger.Errorf("Update: %v", err)
		return err
	}

	if err = tx.Model(data).Association("Role").Replace(data.Role); err != nil {
		d.logger.Errorf("Update Association.Replace: %v", err)
		return err
	}

	return nil
}

func (d *UserDao) Options() (options []map[string]interface{}, err error) {
	err = d.engine.Model(&User{}).Select("id", "name").Order("id").Find(&options).Error
	if err != nil {
		d.logger.Errorf("Options: %v", err)
	}

	return
}
