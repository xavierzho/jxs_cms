package form

import (
	"time"

	"data_backend/pkg"
)

type CachedRequest struct {
	DateTime string `form:"date_time" binding:"required"`
}

func (q *CachedRequest) Parse() (dateTime time.Time, err error) {
	dateTime, err = time.ParseInLocation(pkg.DATE_TIME_FORMAT, q.DateTime, pkg.Location)
	if err != nil {
		return dateTime, err
	}

	return dateTime, nil
}
