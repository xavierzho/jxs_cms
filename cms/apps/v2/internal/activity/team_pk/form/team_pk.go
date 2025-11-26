package form

import (
	"fmt"
	"time"

	"data_backend/apps/v2/internal/activity/team_pk/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/pkg"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type ListRequest struct {
	app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateTimeRange [2]time.Time, err error) {
	q.Pager.Parse()

	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
}

func (q *AllRequest) Parse() (dateTimeRange [2]time.Time, err error) {
	if dateTimeRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	return
}

type Team struct {
	TeamID     int64           `json:"team_id"`
	TeamName   string          `json:"team_name"`
	UserID     int64           `json:"user_id"`
	UserName   string          `json:"user_name"`
	UserAmount decimal.Decimal `json:"user_amount"`
	UserRate   string          `json:"user_rate"`
	TeamAmount decimal.Decimal `json:"team_amount"`
	TeamRate   string          `json:"team_rate"`
	TeamPoint  decimal.Decimal `json:"team_point"`
	TeamNO     int64           `json:"team_no"`
	TeamGap    decimal.Decimal `json:"team_gap"`

	User []*Team `json:"user"`
}

// format
// 转成嵌套结构
func Format(_data []*dao.Team) (teamData []*Team) {
	data := format(_data)
	teamMap := make(map[int64]*Team)
	for _, item := range data {
		if _, ok := teamMap[item.TeamNO]; !ok {
			teamMap[item.TeamNO] = &Team{
				TeamID:     item.TeamID,
				TeamName:   item.TeamName,
				TeamAmount: item.TeamAmount,
				TeamRate:   item.TeamRate,
				TeamPoint:  item.TeamPoint,
				TeamNO:     item.TeamNO,
				TeamGap:    item.TeamGap,
				User:       []*Team{},
			}
			teamData = append(teamData, teamMap[item.TeamNO])
		}

		if item.UserID == 0 {
			continue
		}

		teamMap[item.TeamNO].User = append(teamMap[item.TeamNO].User, &Team{
			UserID:     item.UserID,
			UserName:   item.UserName,
			UserAmount: item.UserAmount,
			UserRate:   item.UserRate,
		})
	}

	return
}

func Format2Excel(dateRange [2]time.Time, _data []*dao.Team) (excelModel *excel.Excel[*Team], err error) {
	data := format(_data)

	reflectMap := map[string]func(source *Team) any{
		"队伍ID":      func(source *Team) any { return source.TeamID },
		"队伍名":       func(source *Team) any { return source.TeamName },
		"用户ID":      func(source *Team) any { return source.UserID },
		"用户名":       func(source *Team) any { return source.UserName },
		"用户流水":      func(source *Team) any { return source.UserAmount },
		"用户积分加成":    func(source *Team) any { return source.UserRate },
		"队伍流水":      func(source *Team) any { return source.TeamAmount },
		"队伍积分加成":    func(source *Team) any { return source.TeamRate },
		"队伍积分":      func(source *Team) any { return source.TeamPoint },
		"队伍排名":      func(source *Team) any { return source.TeamNO },
		"队伍与上一名积分差": func(source *Team) any { return source.TeamGap },
	}

	excelModel = &excel.Excel[*Team]{
		FileName:   fmt.Sprintf("team_pk_%s-%s", dateRange[0].Format(pkg.FILE_DATE_FORMAT), dateRange[1].Format(pkg.FILE_DATE_FORMAT)),
		SheetNames: []string{"队伍PK"},
		SheetNameWithHead: map[string][]string{
			"队伍PK": {
				"队伍ID", "队伍名",
				"用户ID", "用户名", "用户流水", "用户积分加成",
				"队伍流水", "队伍积分加成", "队伍积分", "队伍排名",
				"队伍与上一名积分差",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*Team]{
			"队伍PK": data,
		},
		ReflectMap: map[string]excel.RowReflect[*Team]{
			"队伍PK": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}

func format(data []*dao.Team) (result []*Team) {
	if len(data) == 0 {
		return
	}
	var prevTeamNo = data[0].TeamNO
	var prevTeamPoint = data[0].TeamPoint
	for _, _item := range data {
		item := &Team{
			TeamID:     _item.TeamID,
			TeamName:   _item.TeamName,
			UserID:     _item.UserID,
			UserName:   _item.UserName,
			UserAmount: util.ConvertAmount2Decimal(_item.UserAmount),
			UserRate:   fmt.Sprintf("%.2f%%", _item.UserRate),
			TeamAmount: util.ConvertAmount2Decimal(_item.TeamAmount),
			TeamRate:   fmt.Sprintf("%.2f%%", _item.TeamRate),
			TeamPoint:  util.ConvertAmount2Decimal(_item.TeamPoint),
			TeamNO:     _item.TeamNO,
		}

		if item.TeamNO != prevTeamNo {
			item.TeamGap = util.ConvertAmount2Decimal(prevTeamPoint - _item.TeamPoint)
			prevTeamNo = _item.TeamNO
			prevTeamPoint = _item.TeamPoint
		}

		result = append(result, item)
	}

	return
}
