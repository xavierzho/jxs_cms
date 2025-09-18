package dao

type Item struct {
	ItemID         int64  `gorm:"column:item_id; type:bigint" json:"item_id"`
	ItemName       string `gorm:"column:item_name; type:varchar(64)" json:"item_name"`
	LevelName      string `gorm:"column:level_name; type:longtext" json:"level_name"`
	CoverThumb     string `gorm:"column:cover_thumb; type:varchar(255)" json:"cover_thumb"`
	ShowPrice      int64  `gorm:"column:show_price; type:bigint" json:"show_price"`
	InnerPrice     int64  `gorm:"column:inner_price; type:bigint" json:"inner_price"`
	RecyclingPrice int64  `gorm:"column:recycling_price; type:bigint" json:"recycling_price"`
}
