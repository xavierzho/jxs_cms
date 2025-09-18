package dao

import (
	"context"

	"data_backend/pkg/database"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type AllRequestParamsGroup struct {
	BetFlag                bool
	MarketFlag             bool
	ActivityFlag           ActivityFlag
	AdminFlag              bool
	OrderFlag              bool
	DateTimeParams         database.QueryWhereGroup
	ActivityDateTimeParams database.QueryWhereGroup
	UsersParams            database.QueryWhereGroup
	GachaParams            database.QueryWhereGroup
	ItemParams             database.QueryWhereGroup
	AmountParams           database.QueryWhereGroup
	OtherParams            database.QueryWhereGroup
}

type ActivityFlag struct {
	Flag         bool
	CostAward    bool
	CostRank     bool
	ItemExchange bool
	PrizeWheel   bool //转盘抽奖
}

type ItemDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewItemDao(center *gorm.DB, log *logger.Logger) *ItemDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".ItemDao")))
	return &ItemDao{
		center: center,
		logger: log,
	}
}
