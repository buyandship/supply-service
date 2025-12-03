package mercari

import (
	"time"

	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"gorm.io/gorm"
)

const (
	AccountStatusActive   = "active"
	AccountStatusInactive = "inactive"
	AccountStatusBanned   = "banned"
)

// Account represents the account table in the database.
type Account struct {
	gorm.Model
	Email          string     `gorm:"column:email"` // `email` field
	FamilyName     string     `gorm:"column:family_name"`
	FirstName      string     `gorm:"column:first_name"`
	FamilyNameKana string     `gorm:"column:family_name_kana"`
	FirstNameKana  string     `gorm:"column:first_name_kana"`
	Telephone      string     `gorm:"column:telephone"`
	ZipCode1       string     `gorm:"column:zipcode1"`
	ZipCode2       string     `gorm:"column:zipcode2"`
	Prefecture     string     `gorm:"column:prefecture"`
	City           string     `gorm:"column:city"`
	Address1       string     `gorm:"column:address1"`
	Address2       string     `gorm:"column:address2"`
	Priority       int        `gorm:"column:priority"`
	BannedAt       *time.Time `gorm:"column:banned_at"`
	ActiveAt       *time.Time `gorm:"column:active_at"`
}

func (Account) TableName() string {
	return "account"
}

func (a *Account) Thrift() *supply.Account {
	var bannedAt *string
	var activeAt *string
	status := AccountStatusInactive
	if a.BannedAt != nil {
		b := a.BannedAt.Format(time.RFC3339)
		bannedAt = &b
		status = AccountStatusBanned
	}
	if a.ActiveAt != nil {
		ts := a.ActiveAt.Format(time.RFC3339)
		activeAt = &ts
		status = AccountStatusActive
	}

	account := &supply.Account{
		ID:             int32(a.ID),
		Email:          a.Email,
		FamilyName:     a.FamilyName,
		FirstName:      a.FirstName,
		FamilyNameKana: a.FamilyNameKana,
		FirstNameKana:  a.FirstNameKana,
		Telephone:      a.Telephone,
		ZipCode1:       a.ZipCode1,
		ZipCode2:       a.ZipCode2,
		Prefecture:     a.Prefecture,
		City:           a.City,
		Address1:       a.Address1,
		Address2:       a.Address2,
		Status:         status,
		Priority:       int32(a.Priority),
		BannedAt:       bannedAt,
		ActiveAt:       activeAt,
	}
	return account
}

type SwitchAccountInfo struct {
	FromAccountID int32  `json:"from_account_id"`
	ToAccountID   int32  `json:"to_account_id"`
	Reason        string `json:"reason"`
}
