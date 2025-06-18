package mercari

import (
	"gorm.io/gorm"
)

// Message represents the message table in the database.
type Message struct {
	gorm.Model
	TrxID     string `gorm:"size:255;index"`             // `trx_id` field
	Message   string `gorm:"type:longtext;default:null"` // `message` field
	BuyerID   string `gorm:"size:255"`                   // `buyer_id` field
	AccountID int32  `gorm:"column:account_id;index"`
}

func (Message) TableName() string {
	return "message"
}
