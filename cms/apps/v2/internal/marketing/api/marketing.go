package api

import (
	"data_backend/apps/v2/internal/marketing/dao"
	"data_backend/apps/v2/internal/marketing/service"
	"data_backend/internal/app"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarketingAPI struct {
	svc *service.MarketingSvc
}

func NewMarketingAPI(engine *gorm.DB, log *logger.Logger) *MarketingAPI {
	return &MarketingAPI{
		svc: service.NewMarketingSvc(engine, log),
	}
}

type AttributionReq struct {
	UserID      int64  `json:"user_id" form:"user_id" binding:"required"`
	Channel     string `json:"channel" form:"channel" binding:"required"`
	OAID        string `json:"oaid" form:"oaid"`
	IMEI        string `json:"imei" form:"imei"`
	CallbackURL string `json:"callback_url" form:"callback_url"`
}

// RecordAttribution handles the request to save attribution data
func (api *MarketingAPI) RecordAttribution(c *gin.Context) {
	var req AttributionReq
	response := app.NewResponse(c)
	if err := c.ShouldBind(&req); err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}

	attr := &dao.UserAttribution{
		UserID:      req.UserID,
		Channel:     req.Channel,
		OAID:        req.OAID,
		IMEI:        req.IMEI,
		CallbackURL: req.CallbackURL,
	}

	if err := api.svc.RecordAttribution(attr); err != nil {
		response.ToErrorResponse(errcode.ServerError.WithDetails(err.Error()))
		return
	}

	response.ToResponseOK()
}
