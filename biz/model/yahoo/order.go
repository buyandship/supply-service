package yahoo

import "gorm.io/gorm"

type BidRequest struct {
	gorm.Model
	RequestType  string `gorm:"column:request_type"`
	OrderID      string `gorm:"column:order_id"`
	AuctionID    string `gorm:"column:auction_id"`
	MaxBid       int64  `gorm:"column:max_bid"`
	Quantity     int32  `gorm:"column:quantity"`
	Partial      bool   `gorm:"column:partial"`
	Status       string `gorm:"column:status"`
	ErrorMessage string `gorm:"column:error_message"`
}

func (o *BidRequest) TableName() string {
	return "yahoo.bid_request"
}

type YahooTransaction struct {
	gorm.Model
	BidRequestID string `gorm:"column:bid_request_id"`
	Price        int64  `gorm:"column:price"`
	Status       string `gorm:"column:status"`
	ErrorMessage string `gorm:"column:error_message"`
}

func (o *YahooTransaction) TableName() string {
	return "yahoo.transaction"
}
