package service

import (
	"context"
	"time"

	"data_backend/apps/v2/internal/marketing/dao"
	"data_backend/apps/v2/internal/marketing/xiaomi"
	"data_backend/pkg/logger"

	"gorm.io/gorm"
)

type MarketingSvc struct {
	dao          *dao.AttributionDao
	xiaomiClient *xiaomi.Client
	logger       *logger.Logger
}

func NewMarketingSvc(engine *gorm.DB, log *logger.Logger) *MarketingSvc {
	log = log.WithContext(context.WithValue(context.Background(), logger.ModuleKey, log.ModuleKey().Add(".MarketingSvc")))
	return &MarketingSvc{
		dao:          dao.NewAttributionDao(engine, log),
		xiaomiClient: xiaomi.NewClient(log),
		logger:       log,
	}
}

// RecordAttribution saves the attribution info when app starts or user clicks ad
func (s *MarketingSvc) RecordAttribution(attr *dao.UserAttribution) error {
	// Check if already exists to avoid duplicates if necessary, or just log every click
	// For simplicity, we'll just create a new record.
	return s.dao.Create(attr)
}

// ReportEvent checks if the user is attributed to Xiaomi and reports the event
func (s *MarketingSvc) ReportEvent(userID int64, eventType xiaomi.EventType, amount int64) error {
	// 1. Find attribution info for user
	attr, err := s.dao.GetByUserID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Not an attributed user, ignore
			return nil
		}
		return err
	}

	// 2. Check if it's Xiaomi channel
	if attr.Channel != "xiaomi" {
		return nil
	}

	// 3. Send report
	req := &xiaomi.ReportRequest{
		OAID:        attr.OAID,
		IMEI:        attr.IMEI,
		CallbackURL: attr.CallbackURL,
		EventType:   eventType,
		Amount:      amount,
		Timestamp:   time.Now().UnixMilli(),
	}

	return s.xiaomiClient.ReportEvent(req)
}
