package mercari

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Transaction represents the transaction table in the database.
type Transaction struct {
	gorm.Model
	TrxID         int64          `gorm:"column:trx_id"`                   // `trx_id` field
	RefID         string         `gorm:"size:255;unique"`                 // `ref_id` field
	ItemID        string         `gorm:"size:255"`                        // `item_id` field
	ItemType      string         `gorm:"size:255"`                        // `item_type` field
	ItemDetail    datatypes.JSON `gorm:"type:json"`                       // `item_detail` field
	BuyerID       int32          `gorm:"size:255"`                        // `buyer_id` field
	Price         int64          `gorm:"type:decimal(10,2);default:null"` // `price` field
	PaidPrice     int64          `gorm:"type:decimal(10,2);default:null"` // `paid_price` field
	RefPrice      int64          `gorm:"type:decimal(10,2);default:null"` // `ref_price` field
	FailureReason string         `gorm:"size:255"`                        // `failure_details` field
	Checksum      string         `gorm:"size:255"`                        // `checksum` field
	CouponID      string         `gorm:"size:255"`                        // `coupon_id` field
	Currency      string         `gorm:"size:255"`
}

func (Transaction) TableName() string {
	return "transaction"
}
