# Dao 方法命名规范
## 创建
* (d \*ModelDao) Create(data \*Model) (err error)
  * 创建单个数据, 例如: 创建角色, 用户等

## 储存
* (d \*ModelDao) Save(data []\*Model) (err error)
  * 保存多条数据, 例如: 储存数据到缓存表; 无则创建, 有则更新

## 查找
* (d \*ModelDao) First(queryParams database.QueryWhereGroup) (data \*Model, err error)
  * 按条件查找第一条数据
* (d \*ModelDao) List(queryParams database.QueryWhereGroup[, orderBy app.OrderBy] ,pager app.Pager) (data []\*Model, count int64, err error)
  * 按条件查找并limit, 返回一批数据及总数
  ```go
  (d *ModelDao) List(queryParams database.QueryWhereGroup[, orderBy app.OrderBy] ,pager app.Pager) (data []*Model, count int64, err error){
      err = d.engine.Model(data).Scopes(database.ScopeQuery(queryParams)).Count(&count).
      Scopes(database.Paginate(pager.Page, pager.PageSize)).Find(&data).Error
    if err != nil {
      d.logger.Errorf("List: %v", err)
      return nil, 0, err
    }

    return data, count, nil
  }
  ```
* (d \*ModelDao) All(queryParams database.QueryWhereGroup) (data []\*Model, err error)
  * 按条件查找全部数据

## 更新
* (d \*ModelDao) Update(data \*Model) (err error)
  * 更改数据
* (d \*ModelDao) UpdateAndAssociationReplace(data \*Model) (err error)
  * 更改数据, 并更新关联数据

## 通用 Gorm.DB 对象
* (d \*ModelDao) xxxSql(queryParams database.QueryWhereGroup) (tx \*gorm.DB)

## 生成数据有关的查询
* (d \*ModelDao) Generate(dateRange [2]time.Time) (data []\*Model, err error)

## Options
* (d \*ModelDao) Options() (options []map[string]interface{}, err error)

-------------------------------------------

# Service 方法命名规范
## 创建
* (svc \*ModelSvc) Create(params \*form.ModelCreateRequest) (err error)

## 查找
* (svc \*ModelSvc) List(params \*form.ModelListRequest) ([]form.Model, int64, error)

## 更新
* (svc \*ModelSvc) Update(id uint32, params \*form.ModelUpdateRequest) (err error)

## 生成
## redis 相关
* (svc \*ModelSvc) cachedXxx(xxxRKey string, valueList []string)
* (svc \*ModelSvc) delXxxCache(xxxRKey string) (err error)

## Options
* (svc \*ModelSvc) Options() ([]map[string]interface{}, error)
