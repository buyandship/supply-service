package yahoo

// Category represents a Yahoo Auction category with nested child categories
type Category struct {
	CategoryID       int64       `json:"CategoryId,omitempty" example:"23336"`
	CategoryName     string      `json:"CategoryName,omitempty" example:"コンピュータ"`
	CategoryPath     string      `json:"CategoryPath,omitempty" example:"オークション > コンピュータ"`
	CategoryIDPath   string      `json:"CategoryIdPath,omitempty" example:"0,23336"`
	ParentCategoryID int64       `json:"ParentCategoryId,omitempty"`
	IsLeaf           bool        `json:"IsLeaf,omitempty" example:"false"`
	Depth            int         `json:"Depth,omitempty" example:"1"`
	Order            int         `json:"Order,omitempty" example:"0"`
	IsLink           bool        `json:"IsLink,omitempty" example:"false"`
	IsAdult          bool        `json:"IsAdult,omitempty" example:"false"`
	ChildCategoryNum int         `json:"ChildCategoryNum,omitempty" example:"15"`
	ChildCategory    []*Category `json:"ChildCategory,omitempty"`
	IsLeafToLink     *bool       `json:"IsLeafToLink,omitempty" example:"false"` // Only present in child categories
}

type Transaction struct {
	TransactionID   string  `json:"transaction_id,omitempty" example:"txn_abc123"`
	YsRefID         string  `json:"ys_ref_id,omitempty" example:"YS-REF-001"`
	AuctionID       string  `json:"auction_id,omitempty" example:"x123456789"`
	CurrentPrice    float64 `json:"current_price,omitempty" example:"1300"`
	TransactionType string  `json:"transaction_type,omitempty" example:"BID"`
	Status          string  `json:"status,omitempty" example:"completed"`
	ReqPrice        float64 `json:"req_price,omitempty" example:"1000"`
}

// Account represents a Yahoo account
type Account struct {
	YahooID  string `json:"yahoo_id" example:"chkyj_cp_evjr2p2v"`
	Email    string `json:"email" example:"bnstest.yahoo01@buyandship.com"`
	Password string `json:"password" example:"password"`
	Purpose  string `json:"purpose" example:"for bidding"`
}
