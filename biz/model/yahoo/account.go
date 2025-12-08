package yahoo

// Category represents a Yahoo Auction category with nested child categories
type Category struct {
	CategoryID       int64       `json:"CategoryId,omitempty" example:"23336"`
	CategoryName     string      `json:"CategoryName,omitempty" example:"コンピュータ"`
	CategoryPath     string      `json:"CategoryPath,omitempty" example:"オークション > コンピュータ"`
	CategoryIDPath   string      `json:"CategoryIdPath,omitempty" example:"0,23336"`
	NumOfAuctions    *int64      `json:"NumOfAuctions,omitempty" example:"100"`
	ParentCategoryID *int64      `json:"ParentCategoryId,omitempty"`
	IsLeaf           *bool       `json:"IsLeaf,omitempty" example:"false"`
	Depth            *int        `json:"Depth,omitempty" example:"1"`
	Order            *int        `json:"Order,omitempty" example:"0"`
	IsLink           *bool       `json:"IsLink,omitempty" example:"false"`
	IsAdult          *bool       `json:"IsAdult,omitempty" example:"false"`
	ChildCategoryNum *int        `json:"ChildCategoryNum,omitempty" example:"15"`
	IsLeafToLink     *bool       `json:"IsLeafToLink,omitempty" example:"false"` // Only present in child categories
	ChildCategory    []*Category `json:"ChildCategory,omitempty"`
}

// Account represents a Yahoo account
type Account struct {
	YahooID  string `json:"yahoo_id" example:"chkyj_cp_evjr2p2v"`
	Email    string `json:"email" example:"bnstest.yahoo01@buyandship.com"`
	Password string `json:"password" example:"password"`
	Purpose  string `json:"purpose" example:"for bidding"`
}
