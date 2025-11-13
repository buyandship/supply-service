package yahoo

// Category represents a Yahoo Auction category with nested child categories
type Category struct {
	CategoryID       int64       `json:"CategoryId" example:"23336"`
	CategoryName     string      `json:"CategoryName" example:"コンピュータ"`
	CategoryPath     string      `json:"CategoryPath" example:"オークション > コンピュータ"`
	CategoryIDPath   string      `json:"CategoryIdPath" example:"0,23336"`
	ParentCategoryID int64       `json:"ParentCategoryId"`
	IsLeaf           bool        `json:"IsLeaf" example:"false"`
	Depth            int         `json:"Depth" example:"1"`
	Order            int         `json:"Order" example:"0"`
	IsLink           bool        `json:"IsLink" example:"false"`
	IsAdult          bool        `json:"IsAdult" example:"false"`
	ChildCategoryNum int         `json:"ChildCategoryNum" example:"15"`
	ChildCategory    []*Category `json:"ChildCategory"`
	IsLeafToLink     *bool       `json:"IsLeafToLink,omitempty" example:"false"` // Only present in child categories
}

type Transaction struct {
	TransactionID string `json:"TransactionID" example:"txn_abc123"`
}

// Account represents a Yahoo account
type Account struct {
}
