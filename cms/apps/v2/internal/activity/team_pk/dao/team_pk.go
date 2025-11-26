package dao

import (
	"context"
	"fmt"
	"time"

	"data_backend/internal/app"
	"data_backend/pkg"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

var sql = `
select
	t.*
from
	(
	select
		t.*,
		DENSE_RANK() over(order by t.team_point desc, t.team_amount desc, t.team_id) as team_no
	FROM 
		(
		select
			t.*,
			sum(t.user_amount) over(PARTITION by t.team_id) as team_amount,
			sum(t.user_rate) over(PARTITION by t.team_id) as team_rate,
			FLOOR(sum(t.user_amount) over(PARTITION by t.team_id) * (100+sum(t.user_rate) over(PARTITION by t.team_id))/100) as team_point
		FROM 
			(
			SELECT
				t.team_id, t.team_name, t.user_id, u.nickname as user_name, u.created_at as created_at, t.user_amount, t.user_rate+ifnull(atpc.point_rate, 0) as user_rate
			FROM 
				(
				select
					t.team_id, t.team_name, t.user_id,
					sum(t.user_amount) as user_amount,
					sum(t.user_rate) as user_rate
				from
					(
					select
						t.id,
						tu.team_id, tu.team_name,
						tu.user_id,
						ifnull(t.amount, 0) as user_amount,
						ifnull(t.rate, 0) as user_rate
					from
						(
							select
								atpt.id as team_id, atpt.name as team_name,
								atptu.user_id,
								atpt.created_at
							FROM 
								activity_team_pk_team atpt
								left join activity_team_pk_team_user atptu on atpt.id = atptu.team_id and atptu.deleted_at is null
							WHERE 
								atpt.deleted_at is null
						) tu
						left join (
							select
							    /**这里可以通过配置rate_type来实现按流水金额或固定比例算积分加成**/
								gur.id, gur.created_at, gur.user_id, gur.amount, SUM(IF(IFNULL(atpc.rate_type, 0) = 0,IFNULL(atpc.point_rate, 0) * t.nums,IFNULL(atpc.point_rate, 0) * gm.price * t.nums / 100)) AS rate
							from
								gacha_user_record gur LEFT JOIN gacha_machine gm ON gm.id = gur.gacha_id,
								item i,
								JSON_TABLE(
									JSON_UNQUOTE(gur.items), 
									'$[*]' COLUMNS(
											nums int path '$.Nums',
											item_id bigint path '$.ItemID',
											level_index int path '$.LevelIndex'
										)
								) t
								left join activity_team_pk_config atpc on atpc.deleted_at is null and atpc.type = 'betAwardLevel' and cast(atpc.condition as signed) = t.level_index
							where
								gur.created_at between '%s' and '%s' AND
								gur.deleted_at is null and
								t.item_id = i.id
							group by gur.id, gur.created_at, gur.user_id, gur.amount
						) t on t.user_id = tu.user_id AND t.created_at >= tu.created_at
					) t
				group by
					t.team_id, t.team_name, t.user_id
				) t
				left join users u on t.user_id = u.id
				left join activity_team_pk_config atpc on atpc.deleted_at is null and atpc.type = 'createTime' and atpc.condition <= u.created_at
				
			)t
		order by
			team_point desc, team_amount desc, team_id, user_amount desc, user_rate desc, user_id
		) t
	) t
`

type Team struct {
	TeamID     int64   `gorm:"column:team_id"`
	TeamName   string  `gorm:"column:team_name"`
	UserID     int64   `gorm:"column:user_id"`
	UserName   string  `gorm:"column:user_name"`
	UserAmount int64   `gorm:"column:user_amount"`
	UserRate   float64 `gorm:"column:user_rate"`
	TeamAmount int64   `gorm:"column:team_amount"`
	TeamRate   float64 `gorm:"column:team_rate"`
	TeamPoint  float64 `gorm:"column:team_point"`
	TeamNO     int64   `gorm:"column:team_no"`
}

type TeamPKDao struct {
	center *gorm.DB
	logger *logger.Logger
}

func NewTeamPKDao(center *gorm.DB, log *logger.Logger) *TeamPKDao {
	log = log.WithContext(context.WithValue(log.Context, logger.ModuleKey, log.ModuleKey().Add(".TeamPKDao")))
	return &TeamPKDao{
		center: center,
		logger: log,
	}
}

func (d *TeamPKDao) List(dateRange [2]time.Time, paper app.Pager) (data []*Team, count int64, err error) {
	err = d.center.Table("activity_team_pk_team").Where("deleted_at is null").Count(&count).Error
	if err != nil {
		d.logger.Errorf("List Count err: %v", err)
		return
	}

	var listSql = sql + "where team_no > %d and team_no <= %d"
	err = d.center.Raw(fmt.Sprintf(
		listSql,
		dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Format(pkg.DATE_TIME_MIL_FORMAT),
		paper.GetPageOffset(), paper.GetPageOffset()+paper.PageSize,
	)).Find(&data).Error
	if err != nil {
		d.logger.Errorf("List err: %v", err)
		return
	}

	return
}

func (d *TeamPKDao) All(dateRange [2]time.Time) (data []*Team, err error) {
	err = d.center.Raw(fmt.Sprintf(sql, dateRange[0].Format(pkg.DATE_TIME_MIL_FORMAT), dateRange[1].Format(pkg.DATE_TIME_MIL_FORMAT))).Find(&data).Error
	if err != nil {
		d.logger.Errorf("All err: %v", err)
		return
	}

	return
}
