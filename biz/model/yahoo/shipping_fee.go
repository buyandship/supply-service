package yahoo

import "gorm.io/gorm"

type ShippingFee struct {
	gorm.Model
	ServiceCode int     `json:"service_code"`
	From        string  `json:"from"`
	Size        string  `json:"size"`
	Fee         float64 `json:"fee"`
}

func (s ShippingFee) TableName() string {
	return "yahoo.shipping_fee"
}
