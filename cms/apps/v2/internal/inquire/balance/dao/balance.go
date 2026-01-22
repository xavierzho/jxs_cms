package dao

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DateTimeType string

const (
	DateTimeType_Created DateTimeType = "created" // 时间精度不够 额外增加id排序
	DateTimeType_Finish  DateTimeType = "finish"  // 时间精度不够 额外增加id排序
)

type BalanceComment struct {
	ID        int64
	CreatedAt string
	UserID    uint32
	Comment   string
}

func (b *BalanceComment) Value() (value driver.Value, err error) {
	commentByte, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}

	return string(commentByte), nil
}

type Balance struct {
	ID                  int64                               `gorm:"column:id; type:bigint;" json:"id"`
	CreatedAt           string                              `gorm:"column:created_at; type:datetime(3)" json:"created_at"`
	FinishAt            string                              `gorm:"column:finish_at; type:datetime(3)" json:"finish_at"`
	UserID              int64                               `gorm:"column:user_id; type:bigint" json:"user_id"`
	UserName            string                              `gorm:"column:user_name; type:varchar(64)" json:"user_name"`
	SourceType          int                                 `gorm:"column:source_type; type:int" json:"source_type"`
	GachaName           string                              `gorm:"column:gacha_name; type:longtext" json:"gacha_name"`
	CostAwardName       string                              `gorm:"column:cost_award_name; type:longtext" json:"cost_award_name"`
	ChannelType         string                              `gorm:"column:channel_type; type:longtext" json:"channel_type"`
	PlatformOrderIdPay  string                              `gorm:"column:platform_order_id_pay; type:longtext" json:"platform_order_id_pay"`
	PlatformOrderIdDraw string                              `gorm:"column:platform_order_id_draw; type:longtext" json:"platform_order_id_draw"`
	PaySourceType       int                                 `gorm:"column:pay_source_type; type:int" json:"pay_source_type"`
	BeforeBalance       int64                               `gorm:"column:before_balance; type:bigint" json:"before_balance"`
	AfterBalance        int64                               `gorm:"column:after_balance; type:bigint" json:"after_balance"`
	UpdateAmount        int64                               `gorm:"column:update_amount; type:bigint" json:"update_amount"`
	Comment             datatypes.JSONSlice[BalanceComment] `gorm:"column:comment; type:json" json:"comment"`
	BalanceType         int                                 `gorm:"column:balance_type; type:int" json:"balance_type"`
	ItemExchangeName    string                              `gorm:"column:item_exchange_name; type:longtext" json:"item_exchange_name"`
}

var selectField = []string{
	"bl.id",
	fmt.Sprintf("date_format(bl.created_at, '%s') as created_at", pkg.SQL_DATE_TIME_FORMAT),
	fmt.Sprintf("date_format(bl.finish_at, '%s') as finish_at", pkg.SQL_DATE_TIME_FORMAT),
	"bl.user_id",
	"u.nickname as user_name",
	"bl.source_type",
	"gm.name as gacha_name",
	"cac.name as cost_award_name",
	"ppo.platform_id as channel_type",
	"ppo.platform_order_id as platform_order_id_pay",
	"pdo.remark as platform_order_id_draw",
	"ppo.pay_source_type",
	"bl.before_balance",
	"bl.after_balance",
	"bl.update_amount",
	"bl.comment",
	"bl.type as balance_type",
	"i.name as item_exchange_name",
}

type BalanceDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewBalanceDao(center *gorm.DB, log *logger.Logger) *BalanceDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".BalanceDao")))
	return &BalanceDao{
		center: center,
		logger: log,
	}
}

func (d *BalanceDao) First(queryParams database.QueryWhereGroup) (data *Balance, err error) {
	err = d.allDB(d.center, queryParams).
		Select(selectField).
		Order("bl.id desc").
		Limit(1).
		Scan(&data).Error
	if err != nil {
		d.logger.Errorf("First: %v", err)
		return
	}

	return
}

// 允许 不传时间
func (d *BalanceDao) List(dateTimeType DateTimeType, queryParams database.QueryWhereGroup, pager *app.Pager) (summary map[string]any, data []*Balance, err error) {
	summary = make(map[string]any)

	err = d.allDB(d.center, queryParams).
		Select("count(distinct bl.id) as cnt, count(distinct bl.user_id) as user_cnt, sum(bl.update_amount) as update_amount").
		Scan(&summary).Error
	if err != nil {
		d.logger.Errorf("List summary: %v", err)
		return
	}

	err = d.allDB(d.center, queryParams).
		Select(selectField).
		Scopes(func(d *gorm.DB) *gorm.DB {
			switch dateTimeType {
			case DateTimeType_Created:
				return d.Order("bl.created_at desc, bl.id desc")
			case DateTimeType_Finish:
				return d.Order("bl.finish_at desc, bl.id desc")
			default:
				return d
			}
		}).
		Scopes(database.Paginate(pager.Page, pager.PageSize)).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("List log: %v", err)
		return
	}

	return
}

func (d *BalanceDao) All(dateTimeType DateTimeType, queryParams database.QueryWhereGroup) (data []*Balance, err error) {
	err = d.allDB(d.center, queryParams).
		Select(selectField).
		Scopes(func(d *gorm.DB) *gorm.DB {
			switch dateTimeType {
			case DateTimeType_Created:
				return d.Order("bl.created_at desc, bl.id desc")
			case DateTimeType_Finish:
				return d.Order("bl.finish_at desc, bl.id desc")
			default:
				return d
			}
		}).
		Find(&data).Error
	if err != nil {
		d.logger.Errorf("All: %v", err)
		return
	}

	return
}

func (d *BalanceDao) allDB(tx *gorm.DB, queryParams database.QueryWhereGroup) *gorm.DB {
	return tx.
		Table("users u, balance_log bl").
		Joins("left join gacha_machine gm on bl.source_type between 100 and 199 and bl.source_id = gm.id").
		Joins("left join (select distinct config_id, name from activity_cost_award_config) cac on bl.source_type in (601, 100009) and bl.source_id = cac.config_id").
		Joins("left join activity_item_exchange_config aiec on bl.source_type = 100004 and bl.source_id = aiec.id").
		Joins("left join item i on aiec.item_id = i.id").
		Joins("left join pay_payment_order ppo on bl.source_type in (1, 12) and bl.source_id = ppo.id"). // 金币退款 source_id 是退款订单id 没有第三方id
		Joins("left join pay_payout_order pdo on bl.source_type in (2) and bl.source_id = pdo.id").
		Where("bl.user_id = u.id").
		Scopes(database.ScopeQuery(queryParams))
}

// AddComment 使用 GORM 的 gorm.Expr 进行安全的参数绑定
func (d *BalanceDao) AddComment(id int64, comment *BalanceComment) (err error) {

	jsonExpr := gorm.Expr(
		"JSON_ARRAY_INSERT(IF(comment IS NULL OR comment = '', JSON_ARRAY(), comment), '$[0]', ?)",
		comment,
	)

	err = d.center.Table("balance_log").
		Where("id = ?", id).
		Update("comment", jsonExpr).
		Error

	if err != nil {
		d.logger.Errorf("AddComment: %v", err)
		return err
	}
	return nil
}

// DeleteComment 使用 gorm.Expr 避免字符串拼接
func (d *BalanceDao) DeleteComment(id int64, index int) (err error) {

	jsonPath := fmt.Sprintf("$[%d]", index)

	jsonExpr := gorm.Expr(
		"JSON_REMOVE(comment, ?)",
		jsonPath,
	)
	err = d.center.Table("balance_log").
		Where("id = ? AND comment IS NOT NULL", id).
		Update("comment", jsonExpr).
		Error

	if err != nil {
		d.logger.Errorf("DeleteComment: %v", err)
		return err
	}
	return nil
}
