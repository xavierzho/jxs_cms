package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	iDao "data_backend/internal/dao"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

// TODO 用户名改为 从数据库中查询
type DeliveryOrder struct {
	Date           string    `gorm:"column:date; type:varchar(10); primary_key;" json:"date" form:"date"`
	CreatedAt      time.Time `gorm:"column:created_at; type:datetime; DEFAULT CURRENT_TIMESTAMP" json:"created_at" form:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at; type:datetime; DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at" form:"updated_at"`
	UserID         int64     `gorm:"column:user_id; type:bigint; primary_key;" json:"user_id"`
	UserName       string    `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	ShowPrice      int64     `gorm:"column:show_price; type:bigint" json:"show_price"`
	InnerPrice     int64     `gorm:"column:inner_price; type:bigint" json:"inner_price"`
	RecyclingPrice int64     `gorm:"column:recycling_price; type:bigint" json:"recycling_price"`
}

func (DeliveryOrder) TableName() string {
	return "delivery_order"
}

type DeliveryOrderDao struct {
	*iDao.DailyModelDao[*DeliveryOrder]
	engine *gorm.DB
	center *gorm.DB
	logger *logger.Logger
}

func NewDeliveryOrderDao(engine, center *gorm.DB, log *logger.Logger) *DeliveryOrderDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".DeliveryOrderDao")))
	return &DeliveryOrderDao{
		DailyModelDao: iDao.NewDailyModelDao[*DeliveryOrder](engine, log),
		engine:        engine,
		center:        center,
		logger:        log,
	}
}

func (d *DeliveryOrderDao) Generate(cDate time.Time, queryParams database.QueryWhereGroup) (data []*DeliveryOrder, err error) {
	err = d.center.
		Select(
			fmt.Sprintf("date_format(FROM_UNIXTIME(o.delivery_time / 1000), '%s') as date", pkg.SQL_DATE_FORMAT),
			"o.user_id", "u.nickname as user_name",
			"sum(i.show_price) as show_price",
			"sum(i.inner_price) as inner_price",
			"sum(i.recycling_price) as recycling_price",
		).
		Table("`order` o, item i, users u, ? as j", d.getOrderItemJsonDB()).
		Where("j.item_id = i.id").
		Where("u.id = o.user_id").
		Where("o.state = 4").
		Where("u.is_admin = 0").
		Where("o.delivery_time between ? and ?", cDate.UnixMilli(), cDate.Add(24*time.Hour-time.Millisecond).UnixMilli()).
		Scopes(database.ScopeQuery(queryParams)).
		Group(fmt.Sprintf("date_format(FROM_UNIXTIME(o.delivery_time / 1000), '%s'), o.user_id, u.nickname", pkg.SQL_DATE_FORMAT)).
		Order("`date` DESC").
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("Generate: %v", err)
	}

	return
}

func (d *DeliveryOrderDao) getOrderItemJsonDB() *gorm.DB {
	return d.center.Raw(strings.ReplaceAll(strings.ReplaceAll(`
	JSON_TABLE(
		JSON_UNQUOTE(o.order_items), 
		'$[*]' COLUMNS(
				name varchar(255) path '$.name',
				state int path '$.state',
				item_id bigint path '$.item_id',
				stock_id int path '$.stock_id'
			)
	)
	`, "\t", " "), "\n", " "))
}
