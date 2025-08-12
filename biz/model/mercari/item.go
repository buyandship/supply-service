package mercari

// Item represents a Mercari item with all its details
type Item struct {
	ID                       string                   `json:"id"`
	Status                   string                   `json:"status"`
	Name                     string                   `json:"name"`
	Price                    int                      `json:"price"`
	ItemType                 string                   `json:"item_type"`
	Description              string                   `json:"description"`
	Updated                  int64                    `json:"updated"`
	Created                  int64                    `json:"created"`
	Seller                   Seller                   `json:"seller"`
	Thumbnail                string                   `json:"thumbnail"`
	Photos                   []string                 `json:"photos"`
	ItemCondition            ItemCondition            `json:"item_condition"`
	ShippingPayer            ShippingPayer            `json:"shipping_payer"`
	ShippingDuration         ShippingDuration         `json:"shipping_duration"`
	ItemCategory             ItemCategory             `json:"item_category"`
	ItemBrand                ItemBrand                `json:"item_brand"`
	AnshinItemAuthentication AnshinItemAuthentication `json:"anshin_item_authentication"`
	ItemSizes                []ItemSize               `json:"item_sizes"`
}

// Seller represents the seller information
type Seller struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
}

// ItemCondition represents the condition of the item
type ItemCondition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ShippingPayer represents who pays for shipping
type ShippingPayer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ShippingDuration represents shipping time information
type ShippingDuration struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	MinDays int    `json:"min_days"`
	MaxDays int    `json:"max_days"`
}

// ItemCategory represents the category of the item
type ItemCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ItemBrand represents the brand of the item
type ItemBrand struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SubName string `json:"sub_name"`
}

// AnshinItemAuthentication represents authentication information
type AnshinItemAuthentication struct {
	IsAuthenticatable bool `json:"is_authenticatable"`
	Fee               int  `json:"fee"`
}

// ItemSize represents the size of the item
type ItemSize struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
