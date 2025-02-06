package mercari

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Account represents the account table in the database.
type Account struct {
	gorm.Model
	Email           string         `gorm:"size:255"`       // `email` field
	BuyerID         string         `gorm:"size:255;index"` // `buyer_id` field
	Prefecture      string         `gorm:"size:255"`
	DeliveryAddress datatypes.JSON `gorm:"type:json"` // `delivery_address` field
	AccessToken     string         `gorm:"size:255"`  // `access_token` field
	RefreshToken    string         `gorm:"size:255"`  // `refresh_token` field
	ClientID        string         `gorm:"size:255"`
	ClientSecret    string         `gorm:"size:255"`
}

func (Account) TableName() string {
	return "account"
}
