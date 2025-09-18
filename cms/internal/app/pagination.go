package app

import (
	"time"

	"data_backend/internal/global"
	"data_backend/pkg/convert"

	"github.com/gin-gonic/gin"
)

type Pager struct {
	// 页码
	Page int `form:"page" json:"page" binding:"min=1"`
	// 每页数量
	PageSize int `form:"page_size" json:"page_size" binding:"min=1"`
	// 总行数
	TotalRows int64 `json:"total"`
}

func (p *Pager) Parse() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = global.APPSetting.DefaultPageSize
	}
	if p.PageSize > global.APPSetting.MaxPageSize {
		p.PageSize = global.APPSetting.MaxPageSize
	}
}

func (p Pager) GetPageOffset() int {
	result := 0
	if p.Page > 0 {
		result = (p.Page - 1) * p.PageSize
	}

	return result
}

// 日报表 按日查询时 会用默认值填充空缺的记录，当涉及到分页时会导致数据错乱、分页失效
// 通过 对 日期范围进行分页操作，取当前分页的时间范围(倒序)
// params: dateRange [2]time.Time 必须都是一天中的同一时刻（即差24小时的整数倍）
func (p *Pager) PaginationDateRange(dateRange [2]time.Time) [2]time.Time {
	p.TotalRows = int64(dateRange[1].Sub(dateRange[0]).Hours()/24) + 1
	if int64(p.PageSize) >= p.TotalRows {
		if p.Page == 1 {
			return dateRange
		} else { // 返回无效时间范围
			return [2]time.Time{}
		}
	} else { // 倒序
		offset := p.GetPageOffset()
		endDate := dateRange[1].AddDate(0, 0, -offset)
		if endDate.Before(dateRange[0]) {
			return [2]time.Time{}
		}

		startDate := endDate.AddDate(0, 0, -p.PageSize+1)
		if startDate.Before(dateRange[0]) {
			startDate = dateRange[0]
		}

		return [2]time.Time{startDate, endDate}
	}
}

func GetPage(ctx *gin.Context) int {
	page := convert.StrTo(ctx.Query(PAGE_KEY)).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(ctx *gin.Context) int {
	pageSize := convert.StrTo(ctx.Query(PAGE_SIZE_KEY)).MustInt()
	if pageSize <= 0 {
		return global.APPSetting.DefaultPageSize
	}
	if pageSize > global.APPSetting.MaxPageSize {
		return global.APPSetting.MaxPageSize
	}

	return pageSize
}

func GetPager(ctx *gin.Context) Pager {
	return Pager{
		Page:     GetPage(ctx),
		PageSize: GetPageSize(ctx),
	}
}

func GetStartEndRow(total, page, pageSize int) (startRow, endRow int) {
	startRow = (page-1)*pageSize + 1
	if startRow > total {
		return 0, 0
	}
	if total == 0 {
		return 0, 0
	}
	endRow = page * pageSize
	if endRow > total {
		endRow = total
	}
	return
}
