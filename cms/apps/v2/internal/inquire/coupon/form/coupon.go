package form

import (
	"context"
	"fmt"
	"time"

	cForm "data_backend/apps/v2/internal/common/form"
	"data_backend/apps/v2/internal/inquire/coupon/dao"
	"data_backend/internal/app"
	iForm "data_backend/internal/form"
	"data_backend/internal/global"
	"data_backend/pkg"
	"data_backend/pkg/convert"
	"data_backend/pkg/database"
	"data_backend/pkg/excel"
	"data_backend/pkg/util"

	"github.com/shopspring/decimal"
)

type ListRequest struct {
	*app.Pager
	AllRequest
}

func (q *ListRequest) Parse() (dateRange [2]time.Time, explain dao.Explain, queryParams database.QueryWhereGroup, err error) {
	q.Pager.Parse()
	return q.AllRequest.Parse()
}

type AllRequest struct {
	iForm.DateTimeRangeRequest
	cForm.UserInfoRequest
	CouponID   int64              `form:"coupon_id"`
	CouponType int8               `form:"coupon_type"`
	CouponName string             `form:"coupon_name"`
	Action     []dao.CouponAction `form:"action[]"`
	explain    dao.Explain
}

func (q *AllRequest) Parse() (dateRange [2]time.Time, explain dao.Explain, queryParams database.QueryWhereGroup, err error) {
	if err = q.Valid(); err != nil {
		return
	}

	explain = q.explain

	if dateRange, err = q.DateTimeRangeRequest.Parse(); err != nil {
		return
	}

	if queryParams, err = q.UserInfoRequest.Parse(); err != nil {
		return
	}

	if q.CouponID != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "c.id = ?",
			Value:  []any{q.CouponID},
		})
	}

	if q.CouponType != 0 {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "c.`type` = ?",
			Value:  []any{q.CouponType},
		})
	}

	if q.CouponName != "" {
		queryParams = append(queryParams, database.QueryWhere{
			Prefix: "c.name = ?",
			Value:  []any{q.CouponName},
		})
	}

	return
}

func (q *AllRequest) Valid() error {
	if len(q.Action) == 0 {
		q.explain = dao.Explain{
			Gain:   true,
			Used:   true,
			Expire: true,
		}
	}
	for _, action := range q.Action {
		switch action {
		case dao.CouponAction_Gain:
			q.explain.Gain = true
		case dao.CouponAction_Used:
			q.explain.Used = true
		case dao.CouponAction_Expire:
			q.explain.Expire = true
		default:
			return fmt.Errorf("not expected action: %d", action)
		}
	}

	return nil
}

type Coupon struct {
	DateTime      string          `json:"date_time"`
	CouponID      int64           `json:"coupon_id"`
	CouponTypeStr string          `json:"coupon_type_str"`
	CouponName    string          `json:"coupon_name"`
	UserID        int64           `json:"user_id"`
	UserName      string          `json:"user_name"`
	ActionStr     string          `json:"action_str"`
	Amount        decimal.Decimal `json:"amount"`
}

func Format(ctx context.Context, data []*dao.Coupon) (result []*Coupon) {
	for _, item := range data {
		result = append(result, &Coupon{
			DateTime:      item.DateTime,
			CouponID:      item.CouponID,
			CouponTypeStr: global.I18n.T(ctx, "coupon.type", convert.GetString(item.CouponType)),
			CouponName:    item.CouponName,
			UserID:        item.UserID,
			UserName:      item.UserName,
			ActionStr:     global.I18n.T(ctx, "coupon.action", convert.GetString(uint8(item.Action))),
			Amount:        util.ConvertAmount2Decimal(item.Amount),
		})
	}

	return
}

func Format2Excel(ctx context.Context, dateTimeRange [2]time.Time, _data []*dao.Coupon) (excelModel *excel.Excel[*Coupon], err error) {
	data := Format(ctx, _data)

	reflectMap := map[string]func(source *Coupon) any{
		"时间":    func(source *Coupon) any { return source.DateTime },
		"优惠券ID": func(source *Coupon) any { return source.CouponID },
		"优惠券类型": func(source *Coupon) any { return source.CouponTypeStr },
		"优惠券名称": func(source *Coupon) any { return source.CouponName },
		"用户ID":  func(source *Coupon) any { return source.UserID },
		"用户昵称":  func(source *Coupon) any { return source.UserName },
		"行为":    func(source *Coupon) any { return source.ActionStr },
		"抵扣金额":  func(source *Coupon) any { return source.Amount },
	}

	excelModel = &excel.Excel[*Coupon]{
		FileName:   fmt.Sprintf("user_coupon_log_%s-%s", dateTimeRange[0].Format(pkg.FILE_DATE_TIME_FORMAT), dateTimeRange[1].Format(pkg.FILE_DATE_TIME_FORMAT)),
		SheetNames: []string{"用户优惠券日志"},
		SheetNameWithHead: map[string][]string{
			"用户优惠券日志": {
				"时间",
				"优惠券ID", "优惠券类型", "优惠券名称",
				"用户ID", "用户昵称", "行为",
				"抵扣金额",
			},
		},
		DefaultColWidth:  20,
		DefaultRowHeight: 15,
		Excelize:         nil,
		Data: map[string]excel.SheetData[*Coupon]{
			"用户优惠券日志": data,
		},
		ReflectMap: map[string]excel.RowReflect[*Coupon]{
			"用户优惠券日志": reflectMap,
		},
	}
	err = excelModel.InitExcelFile()
	if err != nil {
		return nil, err
	}

	return
}
