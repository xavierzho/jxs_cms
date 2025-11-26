package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/report/bet/dao"
	"data_backend/apps/v2/internal/report/bet/form"
	iDao "data_backend/internal/dao"
	iErrcode "data_backend/internal/errcode"
	"data_backend/pkg"
	"data_backend/pkg/errcode"
	"data_backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BetSvc struct {
	logger *logger.Logger
	dao    *dao.BetDao
}

func NewBetSvc(ctx *gin.Context, engine, center *gorm.DB, log *logger.Logger) *BetSvc {
	log = log.WithContext(context.WithValue(ctx, logger.ModuleKey, log.ModuleKey().Add(".BetSvc")))
	return &BetSvc{
		logger: log,
		dao:    dao.NewBetDao(engine, center, log),
	}
}

func (svc *BetSvc) Generate(params *form.GenerateRequest) (e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return errcode.InvalidParams.WithDetails(err.Error())
	}

	for cDate := dateRange[0]; !cDate.After(dateRange[1]); cDate = cDate.AddDate(0, 0, 1) {
		if err = svc.generate(cDate); err != nil {
			return iErrcode.SQLExecFail.WithDetails(err.Error())
		}
	}

	return nil
}

func (svc *BetSvc) generate(cDate time.Time) (err error) {
	dataList, err := svc._generate(cDate)
	if err != nil {
		return err
	}

	if err = svc.dao.Save(dataList...); err != nil {
		return err
	}

	return nil
}

func (svc *BetSvc) _generate(cDate time.Time) (dataList []*dao.Bet, err error) {
	dataList = make([]*dao.Bet, 0, len(form.BET_TYPE_LIST))
	dataMap := make(map[string]*dao.Bet, len(form.BET_TYPE_LIST))
	for _, betType := range form.BET_TYPE_LIST {
		dataMap[betType] = &dao.Bet{
			DailyTypeModel: iDao.DailyTypeModel{
				Date:     cDate.Format(pkg.DATE_FORMAT),
				DataType: betType,
			},
		}
		dataList = append(dataList, dataMap[betType])
	}

	dataBet, dataAmount, dataBox, dataPay, err := svc.dao.Generate(cDate)
	if err != nil {
		return nil, err
	}

	for _, item := range dataBet {
		betType := convertBetType(item.DataType)
		dataMap[betType].UserCnt = item.UserCnt
		dataMap[betType].BetNums = item.BetNums
	}

	for _, item := range dataAmount {
		betType := convertBetType(item.DataType)
		dataMap[betType].Amount = item.Amount
	}

	for _, item := range dataBox {
		betType := convertBetType(item.DataType)
		dataMap[betType].BoxCntRemaining = item.BoxCntRemaining
		dataMap[betType].BoxCntNew = item.BoxCntNew
		dataMap[betType].BoxCntClose = item.BoxCntClose
		// 若跑历史数据则剩余箱子置为0
		if dataMap[betType].Date < time.Now().Format(pkg.DATE_FORMAT) {
			dataMap[betType].BoxCntRemaining = 0
		}
	}

	for _, item := range dataPay {
		betType := convertBetType(item.DataType)
		dataMap[betType].AmountWeChat = item.AmountWeChat
		dataMap[betType].AmountAli = item.AmountAli
	}

	return
}

func convertBetType(dataType string) string {
	switch dataType {
	case "1", "101":
		return form.BET_TYPE_FIRST_PRIZE
	case "2", "102":
		return form.BET_TYPE_GASHAPON
	case "3", "103":
		return form.BET_TYPE_CHAO
	case "4", "104":
		return form.BET_TYPE_HOLE
	case "5", "105":
		return form.BET_TYPE_CHAO_SHE
	case "6", "106":
		return form.BET_TYPE_SHARE_BILL
	default:
		return ""
	}
}

func (svc *BetSvc) All(params *form.AllRequest) (result []form.Bet, e *errcode.Error) {
	dateRange, err := params.Parse()
	if err != nil {
		return nil, errcode.InvalidParams.WithDetails(err.Error())
	}

	data, err := svc.dao.All(dateRange, params.DataType, nil)
	if err != nil {
		return nil, errcode.QueryFail.WithDetails(err.Error())
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, pkg.Location)
	if !dateRange[0].After(today) && !dateRange[1].Before(today) {
		todayData, err := svc._generate(today)
		if err != nil {
			return nil, errcode.QueryFail.WithDetails(err.Error())
		}

		for _, item := range todayData {
			if item.DataType == params.DataType {
				data = append(data, item)
				break
			}
		}

	}

	result, err = form.Format(dateRange, data)
	if err != nil {
		svc.logger.Errorf("All, Format: %v", err)
		return nil, errcode.TransformFail.WithDetails(err.Error())
	}

	return result, nil
}
