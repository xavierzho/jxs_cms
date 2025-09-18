package dao

import (
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

// ? 考虑 是否限制 T 必须有主键
type Dao[T any] struct {
	engine *gorm.DB
	logger *logger.Logger
}

func NewDao[T any](engine *gorm.DB, log *logger.Logger) *Dao[T] {
	return &Dao[T]{engine: engine, logger: log}
}

func (d *Dao[T]) Create(data T) (err error) {
	if err = d.engine.Create(data).Error; err != nil {
		d.logger.Errorf("Create: %v", err)
		return err
	}

	return nil
}

func (d *Dao[T]) Save(data ...T) (err error) {
	if len(data) == 0 {
		return nil
	}
	if err = d.engine.Save(data).Error; err != nil {
		d.logger.Errorf("Save: %v", err)
		return err
	}

	return nil
}

func (d *Dao[T]) List(queryParams database.QueryWhereGroup, pager app.Pager) (data []T, count int64, err error) {
	err = d.engine.
		Model(data).
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

func (d *Dao[T]) All(queryParams database.QueryWhereGroup) (data []T, err error) {
	err = d.engine.
		Model(data).
		Scopes(database.ScopeQuery(queryParams)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *Dao[T]) Update(selectField, omitField []string, data T) (err error) {
	if err = d.engine.Select(selectField).Omit(omitField...).Updates(data).Error; err != nil {
		d.logger.Errorf("Update: %v", err)
		return err
	}

	return nil
}

type Model struct {
	ID        uint32    `gorm:"column:id; primary_key" json:"id" form:"id"`
	CreatedAt time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
}

type ModelDao[T any] struct {
	*Dao[T]
}

func NewModelDao[T any](engine *gorm.DB, log *logger.Logger) *ModelDao[T] {
	return &ModelDao[T]{Dao: NewDao[T](engine, log)}
}

func (d *ModelDao[T]) First(queryParams database.QueryWhereGroup) (data T, err error) {
	err = d.engine.
		Model(data).
		Scopes(database.ScopeQuery(queryParams)).
		Order("id desc").
		First(&data).Error
	if err != nil {
		d.logger.Errorf("First: %v", err)
		return data, err
	}

	return data, nil
}

func (d *ModelDao[T]) List(queryParams database.QueryWhereGroup, pager app.Pager) (data []T, count int64, err error) {
	err = d.engine.
		Model(data).
		Scopes(database.ScopeQuery(queryParams)).
		Count(&count).
		Order("id desc").
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List: %v", err)
		return nil, 0, err
	}

	return data, count, nil
}

func (d *ModelDao[T]) All(queryParams database.QueryWhereGroup) (data []T, err error) {
	err = d.engine.
		Model(data).
		Scopes(database.ScopeQuery(queryParams)).
		Order("id desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *ModelDao[T]) Update(data T) (err error) {
	return d.Dao.Update([]string{"*"}, []string{"date"}, data)
}

type DailyModel struct {
	Date      string    `gorm:"column:date; type:varchar(10); primary_key" json:"date" form:"date"`
	CreatedAt time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
}

type DailyModelDao[T any] struct {
	*Dao[T]
}

func NewDailyModelDao[T any](engine *gorm.DB, log *logger.Logger) *DailyModelDao[T] {
	return &DailyModelDao[T]{Dao: NewDao[T](engine, log)}
}

func (d *DailyModelDao[T]) First(queryParams database.QueryWhereGroup) (data T, err error) {
	err = d.engine.
		Model(data).
		Scopes(database.ScopeQuery(queryParams)).
		Order("date desc").
		First(&data).Error
	if err != nil {
		d.logger.Errorf("First: %v", err)
		return data, err
	}

	return data, nil
}

func (d *DailyModelDao[T]) List(dateRange [2]time.Time, queryParams database.QueryWhereGroup, pager app.Pager) (data []T, count int64, err error) {
	err = d.all(dateRange, queryParams).
		Model(data).
		Count(&count).
		Order("date desc").
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List: %v", err)
		return nil, 0, err
	}

	return data, count, nil
}

func (d *DailyModelDao[T]) ListAndSummary(summaryField []string, dateRange [2]time.Time, queryParams database.QueryWhereGroup, pager app.Pager) (summary map[string]any, data []T, err error) {
	summary = make(map[string]any)

	err = d.all(dateRange, queryParams).Model(data).Order("date desc").Scopes(database.Paginate(pager.Page, pager.PageSize)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("ListAndSummary list: %v", err)
		return
	}

	err = d.all(dateRange, queryParams).Model(data).Select(summaryField).Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("ListAndSummary summary: %v", err)
		return
	}

	return
}

func (d *DailyModelDao[T]) All(dateRange [2]time.Time, queryParams database.QueryWhereGroup) (data []T, err error) {
	err = d.all(dateRange, queryParams).
		Model(data).
		Order("date desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *DailyModelDao[T]) all(dateRange [2]time.Time, queryParams database.QueryWhereGroup) *gorm.DB {
	return d.engine.
		Where("date between ? and ?", dateRange[0].Format(pkg.DATE_FORMAT), dateRange[1].Format(pkg.DATE_FORMAT)).
		Scopes(database.ScopeQuery(queryParams))
}

func (d *DailyModelDao[T]) Update(data T) (err error) {
	return d.Dao.Update([]string{"*"}, []string{"date"}, data)
}

type DailyTypeModel struct {
	Date      string    `gorm:"column:date; type:varchar(10); primary_key" json:"date" form:"form"`
	DataType  string    `gorm:"column:data_type; type:varchar(32); primary_key" json:"data_type" form:"data_type"`
	CreatedAt time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
}

type DailyTypeModelDao[T any] struct {
	*Dao[T]
}

func NewDailyTypeModelDao[T any](engine *gorm.DB, log *logger.Logger) *DailyTypeModelDao[T] {
	return &DailyTypeModelDao[T]{Dao: NewDao[T](engine, log)}
}

func (d *DailyTypeModelDao[T]) First(dataType string, queryParams database.QueryWhereGroup) (data T, err error) {
	err = d.engine.
		Model(data).
		Where("data_type = ?", dataType).
		Scopes(database.ScopeQuery(queryParams)).
		Order("date desc").
		First(&data).Error
	if err != nil {
		d.logger.Errorf("First: %v", err)
		return data, err
	}

	return data, nil
}

func (d *DailyTypeModelDao[T]) List(dateRange [2]time.Time, dataType string, queryParams database.QueryWhereGroup, pager app.Pager) (data []T, count int64, err error) {
	err = d.engine.
		Model(data).
		Where("date between ? and ?", dateRange[0].Format(pkg.DATE_FORMAT), dateRange[1].Format(pkg.DATE_FORMAT)).
		Where("data_type = ?", dataType).
		Scopes(database.ScopeQuery(queryParams)).
		Count(&count).
		Order("date desc").
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List: %v", err)
		return nil, 0, err
	}

	return data, count, nil
}

func (d *DailyTypeModelDao[T]) All(dateRange [2]time.Time, dataType string, queryParams database.QueryWhereGroup) (data []T, err error) {
	err = d.engine.
		Model(data).
		Where("date between ? and ?", dateRange[0].Format(pkg.DATE_FORMAT), dateRange[1].Format(pkg.DATE_FORMAT)).
		Where("data_type = ?", dataType).
		Scopes(database.ScopeQuery(queryParams)).
		Order("date desc").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return nil, err
	}

	return data, nil
}

func (d *DailyTypeModelDao[T]) Update(data T) (err error) {
	return d.Dao.Update([]string{"*"}, []string{"date", "data_type"}, data)
}
