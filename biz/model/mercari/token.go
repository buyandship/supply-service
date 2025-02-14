package mercari

import "gorm.io/gorm"

type Token struct {
	gorm.Model
	AccessToken  string `gorm:"size:255"` // `access_token` field
	RefreshToken string `gorm:"size:255"` // `refresh_token` field
	ExpiresIn    int32  `gorm:"column:expires_in"`
	Scope        string `gorm:"size:255"`
	TokenType    string `gorm:"size:255"`
}

func (Token) TableName() string {
	return "token"
}
