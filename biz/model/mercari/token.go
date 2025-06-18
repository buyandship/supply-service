package mercari

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	AccessToken  string `gorm:"size:255"` // `access_token` field
	RefreshToken string `gorm:"size:255"` // `refresh_token` field
	ExpiresIn    int32  `gorm:"column:expires_in"`
	Scope        string `gorm:"size:255"`
	TokenType    string `gorm:"size:255"`
	AccountID    int32  `gorm:"column:account_id;index"`
}

func (Token) TableName() string {
	return "token"
}

func (m *Token) Expired() bool {
	if m == nil {
		return true
	}
	expiredTime := m.CreatedAt.Add(time.Duration(m.ExpiresIn-60) * time.Second)
	if time.Now().After(expiredTime) {
		return true
	}
	return false
}
