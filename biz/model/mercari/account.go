package mercari

import (
	"gorm.io/gorm"
)

// Account represents the account table in the database.
type Account struct {
	gorm.Model
	Email          string `gorm:"column:email"`          // `email` field
	BuyerID        int32  `gorm:"column:buyer_id;index"` // `buyer_id` field
	FamilyName     string `gorm:"column:family_name"`
	FirstName      string `gorm:"column:first_name"`
	FamilyNameKana string `gorm:"column:family_name_kana"`
	FirstNameKana  string `gorm:"column:first_name_kana"`
	Telephone      string `gorm:"column:telephone"`
	ZipCode1       string `gorm:"column:zipcode1"`
	ZipCode2       string `gorm:"column:zipcode2"`
	Prefecture     string `gorm:"column:prefecture"`
	City           string `gorm:"column:city"`
	Address1       string `gorm:"column:address1"`
	Address2       string `gorm:"column:address2"`
}

func (Account) TableName() string {
	return "account"
}
