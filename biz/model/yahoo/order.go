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

	TransactionID string `gorm:"-"`
}

func (o *BidRequest) TableName() string {
	return "yahoo.bid_request"
}

type BidAuctionItem struct {
	gorm.Model
	BidRequestID     string `gorm:"column:bid_request_id"`
	AuctionID        string `gorm:"column:auction_id"`
	Status           string `gorm:"column:status"`
	Name             string `gorm:"column:name"`
	CurrentPrice     string `gorm:"column:current_price"`
	BuyoutPrice      string `gorm:"column:buyout_price"`
	ItemType         string `gorm:"column:item_type"`
	Description      string `gorm:"column:description"`
	StartTime        string `gorm:"column:start_time"`
	EndTime          string `gorm:"column:end_time"`
	ItemCategoryID   string `gorm:"column:item_category_id"`
	ItemCondition    string `gorm:"column:item_condition"`
	ItemBrand        string `gorm:"column:item_brand"`
	ItemWatchListNum string `gorm:"column:item_watch_list_num"`
	SellerID         string `gorm:"column:seller_id"`
	SellerRating     string `gorm:"column:seller_rating"`
}

func (o *BidAuctionItem) TableName() string {
	return "yahoo.bid_auction_item"
}

type YahooTransaction struct {
	gorm.Model
	BidRequestID  string `gorm:"column:bid_request_id"`
	TransactionID string `gorm:"column:transaction_id"`
	Price         int64  `gorm:"column:price"`
	Status        string `gorm:"column:status"`
	ErrorMessage  string `gorm:"column:error_message"`
}

func (o *YahooTransaction) TableName() string {
	return "yahoo.transaction"
}
