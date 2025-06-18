package mercari

import (
	"gorm.io/gorm"
)

// Review represents the review table in the database.
type Review struct {
	gorm.Model
	TrxID     string `gorm:"size:255;index"`             // `trx_id` field
	Fame      string `gorm:"size:6;index"`               // 	`fame` field
	Review    string `gorm:"type:longtext;default:null"` // `review` field
	BuyerID   string `gorm:"size:255"`                   // `buyer_id` field
	AccountID int32  `gorm:"column:account_id;index"`
}

func (Review) TableName() string {
	return "review"
}
