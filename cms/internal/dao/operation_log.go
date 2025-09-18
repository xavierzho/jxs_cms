package dao

import (
	"context"
	"net/url"

	"data_backend/internal/app"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type OperationLog struct {
	Model

	UserID     uint32 `gorm:"column:user_id; type:int" json:"user_id"`
	ModuleID   string `gorm:"column:module_id; type:varchar(50)" json:"module_id"`
	ModuleName string `gorm:"column:module_name; type:varchar(50); index:ind_modulename" json:"module_name"`
	Operation  string `gorm:"column:operation; type:varchar(100)" json:"operation"`
	Request    string `gorm:"column:request; type:varchar(1024)" json:"request"`

	User User `gorm:"foreignKey:user_id"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type OperationLogDao struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewOperationLogDao(engine *gorm.DB, log *logger.Logger) *OperationLogDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".OperationLogDao")))

	return &OperationLogDao{
		engine: engine,
		logger: log,
	}
}

func (d *OperationLogDao) Create(form url.Values, data *OperationLog) (err error) {
	if err = d.engine.Create(data).Error; err != nil {
		d.logger.Errorf("Create: %+v, params: %+v: %v", data, form, err)
		return err
	}

	return nil
}

func (d *OperationLogDao) List(queryParams []database.QueryWhere, pager app.Pager) (logModel []*OperationLog, count int64, err error) {
	err = d.engine.
		Model(logModel).
		Preload("User").
		Scopes(database.ScopeQuery(queryParams)).
		Count(&count).
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Order("id desc").
		Find(&logModel).Error
	if err != nil {
		d.logger.Errorf("List: %v", err)
		return nil, 0, err
	}

	return logModel, count, err
}
