package form

import "github.com/shopspring/decimal"

type Item struct {
	ItemID         string          `json:"item_id"`
	ItemName       string          `json:"item_name"`
	LevelName      string          `json:"level_name"`
	CoverThumb     string          `json:"cover_thumb"`
	ShowPrice      decimal.Decimal `json:"show_price"`
	InnerPrice     decimal.Decimal `json:"inner_price"`
	RecyclingPrice decimal.Decimal `json:"recycling_price"`
}
