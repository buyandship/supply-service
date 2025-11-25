package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
)

// SearchAuctionsRequest represents the request parameters for searching Yahoo Auctions
type SearchAuctionsRequest struct {
	// Query (required): Search keywords
	Query string `json:"query"`

	// YahooAccountID (optional): Yahoo account ID for audit logging
	YahooAccountID string `json:"yahoo_account_id,omitempty"`

	// YsRefID (optional): Y!S request unique identifier for tracking
	YsRefID string `json:"ys_ref_id,omitempty"`

	// Type: Search type (all/any)
	Type string `json:"type,omitempty"`

	// Category: Category ID for filtering results
	Category int `json:"category,omitempty"`

	// ExceptCategory: Comma-separated category IDs to exclude
	ExceptCategory string `json:"except_category,omitempty"`

	// Page: Page number for pagination (default: 1)
	Page int `json:"page,omitempty"`

	// Sort: Sort field (end, img, bids, cbids, bidorbuy, affiliate)
	Sort string `json:"sort,omitempty"`

	// Order: Sort direction (a=ascending, d=descending)
	Order string `json:"order,omitempty"`

	// Store: Store filter (0=all, 1=store only, 2=individual sellers only)
	Store int `json:"store,omitempty"`

	// AucMinPrice: Minimum current price
	AucMinPrice int `json:"aucminprice,omitempty"`

	// AucMaxPrice: Maximum current price
	AucMaxPrice int `json:"aucmaxprice,omitempty"`

	// AucMinBidorbuyPrice: Minimum buy-now price
	AucMinBidorbuyPrice int `json:"aucmin_bidorbuy_price,omitempty"`

	// AucMaxBidorbuyPrice: Maximum buy-now price
	AucMaxBidorbuyPrice int `json:"aucmax_bidorbuy_price,omitempty"`

	// LocCd: Location code (1-48)
	LocCd int `json:"loc_cd,omitempty"`

	// EasyPayment: Yahoo Easy Payment available (1=yes, 0=all)
	EasyPayment int `json:"easypayment,omitempty"`

	// New: New icon items only (1=yes, 0=all)
	New int `json:"new,omitempty"`

	// FreeShipping: Free shipping items only (1=yes, 0=all)
	FreeShipping int `json:"freeshipping,omitempty"`

	// WrappingIcon: Gift wrapping available (1=yes, 0=all)
	WrappingIcon int `json:"wrappingicon,omitempty"`

	// BuyNow: Buy-now price set (1=yes, 0=all)
	BuyNow int `json:"buynow,omitempty"`

	// Thumbnail: Items with images (1=yes, 0=all)
	Thumbnail int `json:"thumbnail,omitempty"`

	// Attn: Featured auctions only (1=yes, 0=all)
	Attn int `json:"attn,omitempty"`

	// Point: Yahoo Points items (1=yes, 0=all)
	Point int `json:"point,omitempty"`

	// ItemStatus: Item condition (0=all, 1=new, 2=used)
	ItemStatus int `json:"item_status,omitempty"`

	// Adf: Include adult category items (1=yes)
	Adf int `json:"adf,omitempty"`

	// SellerAucUserID: Filter by seller user ID
	SellerAucUserID string `json:"seller_auc_user_id,omitempty"`

	// F: Search target field (0x2/0x4/0x8)
	F string `json:"f,omitempty"`

	// Ngram: Search method (0=MA, 1=NGram)
	Ngram int `json:"ngram,omitempty"`

	// Fixed: Auction type (0=all, 1=fixed price only, 2=auction only)
	Fixed int `json:"fixed,omitempty"`

	// MinCharity: Minimum charity donation rate (%)
	MinCharity int `json:"min_charity,omitempty"`

	// MaxCharity: Maximum charity donation rate (%)
	MaxCharity int `json:"max_charity,omitempty"`

	// MinAffiliate: Minimum affiliate rate (%)
	MinAffiliate int `json:"min_affiliate,omitempty"`

	// MaxAffiliate: Maximum affiliate rate (%)
	MaxAffiliate int `json:"max_affiliate,omitempty"`

	// Timebuf: Filter by remaining time (seconds)
	Timebuf int `json:"timebuf,omitempty"`

	// Ranking: Ranking type (current/popular)
	Ranking string `json:"ranking,omitempty"`

	// BlackSellerAucUserID: Comma-separated seller IDs to exclude
	BlackSellerAucUserID string `json:"black_seller_auc_user_id,omitempty"`

	// Featured: Featured auctions (on/off)
	Featured string `json:"featured,omitempty"`

	// Sort2: Secondary sort field
	Sort2 string `json:"sort2,omitempty"`

	// Order2: Secondary sort order
	Order2 string `json:"order2,omitempty"`

	// MinStart: Minimum start time (Unix timestamp)
	MinStart int64 `json:"min_start,omitempty"`

	// MaxStart: Maximum start time (Unix timestamp)
	MaxStart int64 `json:"max_start,omitempty"`

	// ExceptShoppingItem: Exclude Yahoo! Shopping items
	ExceptShoppingItem bool `json:"except_shoppingitem,omitempty"`
}

// AuctionItemListSeller represents seller information in auction item list
type AuctionItemListSeller struct {
	AucUserId            string `json:"AucUserId,omitempty"`
	AucUserIdItemListURL string `json:"AucUserIdItemListUrl,omitempty"`
	AucUserIdRatingURL   string `json:"AucUserIdRatingUrl,omitempty"`
}

// AuctionItemListImage represents image information in auction item list
type AuctionItemListImage struct {
	URL    string `json:"Url,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// AuctionItemListOption represents option information in auction item list
type AuctionItemListOption struct {
	NewIcon           string `json:"NewIcon,omitempty"`
	StoreIcon         string `json:"StoreIcon,omitempty"`
	FeaturedIcon      string `json:"FeaturedIcon,omitempty"`
	FreeshippingIcon  string `json:"FreeshippingIcon,omitempty"`
	BuynowIcon        string `json:"BuynowIcon,omitempty"`
	EasyPaymentIcon   string `json:"EasyPaymentIcon,omitempty"`
	IsBold            bool   `json:"IsBold,omitempty"`
	IsBackGroundColor bool   `json:"IsBackGroundColor,omitempty"`
	IsOffer           bool   `json:"IsOffer,omitempty"`
	IsCharity         bool   `json:"IsCharity,omitempty"`
}

// Affiliate represents affiliate information
type Affiliate struct {
	Rate int `json:"Rate"`
}

// AuctionItemListDetail represents an auction item in the search/list results
type AuctionItemListDetail struct {
	AuctionID        string                `json:"AuctionID,omitempty"`
	Title            string                `json:"Title,omitempty"`
	CategoryId       int                   `json:"CategoryId,omitempty"`
	IsCrossListing   bool                  `json:"IsCrossListing,omitempty"`
	IsFleaMarket     bool                  `json:"IsFleaMarket,omitempty"`
	Seller           AuctionItemListSeller `json:"Seller,omitempty"`
	ItemUrl          string                `json:"ItemUrl,omitempty"`
	AuctionItemUrl   string                `json:"AuctionItemUrl,omitempty"`
	Image            AuctionItemListImage  `json:"Image,omitempty"`
	OriginalImageNum int32                 `json:"OriginalImageNum,omitempty"`
	CurrentPrice     float64               `json:"CurrentPrice,omitempty"`
	Bids             int                   `json:"Bids,omitempty"`
	EndTime          string                `json:"EndTime,omitempty"`
	StartTime        string                `json:"StartTime,omitempty"`
	BidOrBuy         float64               `json:"BidOrBuy,omitempty"`
	IsReserved       bool                  `json:"IsReserved,omitempty"`
	CharityOption    CharityOption         `json:"CharityOption,omitempty"`
	Affiliate        Affiliate             `json:"Affiliate,omitempty"`
	Option           AuctionItemListOption `json:"Option,omitempty"`
	IsAdult          bool                  `json:"IsAdult,omitempty"`
}

type SearchAuctionResult struct {
	Item []AuctionItemListDetail `json:"Item"`
}

type SearchAuctionsResponse struct {
	ResultSet struct {
		TotalResultsAvailable int                 `json:"@totalResultsAvailable"`
		TotalResultsReturned  int                 `json:"@totalResultsReturned"`
		FirstResultPosition   int                 `json:"@firstResultPosition"`
		Result                SearchAuctionResult `json:"Result"`
	} `json:"ResultSet"`
}

func (c *Client) SearchAuctions(ctx context.Context, req *SearchAuctionsRequest) (*SearchAuctionsResponse, error) {
	params := url.Values{}
	params.Set("query", req.Query)
	params.Set("yahoo_account_id", config.MasterYahooAccountID)
	if req.YsRefID != "" {
		params.Set("ys_ref_id", req.YsRefID)
	}
	if req.Type != "" {
		params.Set("type", req.Type)
	}
	if req.Category != 0 {
		params.Set("category", strconv.Itoa(req.Category))
	}
	if req.ExceptCategory != "" {
		params.Set("except_category", req.ExceptCategory)
	}
	if req.Page != 0 {
		params.Set("page", strconv.Itoa(req.Page))
	}
	if req.Sort != "" {
		params.Set("sort", req.Sort)
	}
	if req.Order != "" {
		params.Set("order", req.Order)
	}
	if req.Store != 0 {
		params.Set("store", strconv.Itoa(req.Store))
	}
	if req.AucMinPrice != 0 {
		params.Set("aucminprice", strconv.Itoa(req.AucMinPrice))
	}
	if req.AucMaxPrice != 0 {
		params.Set("aucmaxprice", strconv.Itoa(req.AucMaxPrice))
	}
	if req.AucMinBidorbuyPrice != 0 {
		params.Set("aucmin_bidorbuy_price", strconv.Itoa(req.AucMinBidorbuyPrice))
	}
	if req.AucMaxBidorbuyPrice != 0 {
		params.Set("aucmax_bidorbuy_price", strconv.Itoa(req.AucMaxBidorbuyPrice))
	}
	if req.LocCd != 0 {
		params.Set("loc_cd", strconv.Itoa(req.LocCd))
	}
	if req.EasyPayment != 0 {
		params.Set("easypayment", strconv.Itoa(req.EasyPayment))
	}
	if req.New != 0 {
		params.Set("new", strconv.Itoa(req.New))
	}
	if req.FreeShipping != 0 {
		params.Set("freeshipping", strconv.Itoa(req.FreeShipping))
	}
	if req.WrappingIcon != 0 {
		params.Set("wrappingicon", strconv.Itoa(req.WrappingIcon))
	}
	if req.BuyNow != 0 {
		params.Set("buynow", strconv.Itoa(req.BuyNow))
	}
	if req.Thumbnail != 0 {
		params.Set("thumbnail", strconv.Itoa(req.Thumbnail))
	}
	if req.Attn != 0 {
		params.Set("attn", strconv.Itoa(req.Attn))
	}
	if req.Point != 0 {
		params.Set("point", strconv.Itoa(req.Point))
	}
	if req.ItemStatus != 0 {
		params.Set("item_status", strconv.Itoa(req.ItemStatus))
	}
	if req.Adf != 0 {
		params.Set("adf", strconv.Itoa(req.Adf))
	}
	if req.SellerAucUserID != "" {
		params.Set("seller_auc_user_id", req.SellerAucUserID)
	}
	if req.F != "" {
		params.Set("f", req.F)
	}
	if req.Ngram != 0 {
		params.Set("ngram", strconv.Itoa(req.Ngram))
	}
	if req.Fixed != 0 {
		params.Set("fixed", strconv.Itoa(req.Fixed))
	}
	if req.MinCharity != 0 {
		params.Set("min_charity", strconv.Itoa(req.MinCharity))
	}
	if req.MaxCharity != 0 {
		params.Set("max_charity", strconv.Itoa(req.MaxCharity))
	}
	if req.MinAffiliate != 0 {
		params.Set("min_affiliate", strconv.Itoa(req.MinAffiliate))
	}
	if req.MaxAffiliate != 0 {
		params.Set("max_affiliate", strconv.Itoa(req.MaxAffiliate))
	}
	if req.Timebuf != 0 {
		params.Set("timebuf", strconv.Itoa(req.Timebuf))
	}
	if req.Ranking != "" {
		params.Set("ranking", req.Ranking)
	}
	if req.BlackSellerAucUserID != "" {
		params.Set("black_seller_auc_user_id", req.BlackSellerAucUserID)
	}
	if req.Featured != "" {
		params.Set("featured", req.Featured)
	}
	if req.Sort2 != "" {
		params.Set("sort2", req.Sort2)
	}
	if req.Order2 != "" {
		params.Set("order2", req.Order2)
	}
	if req.MinStart != 0 {
		params.Set("min_start", strconv.FormatInt(req.MinStart, 10))
	}
	if req.MaxStart != 0 {
		params.Set("max_start", strconv.FormatInt(req.MaxStart, 10))
	}
	if req.ExceptShoppingItem {
		params.Set("except_shoppingitem", "true")
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/search", params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusInternalServerError:
				{
					return nil, bizErr.BizError{
						Status:  http.StatusBadRequest,
						ErrCode: http.StatusBadRequest,
						ErrMsg:  "yahoo returns internal server error",
					}
				}
			}
		}
		return nil, err
	}

	var searchAuctionsResponse SearchAuctionsResponse
	if err := c.parseResponse(resp, &searchAuctionsResponse); err != nil {
		return nil, err
	}

	return &searchAuctionsResponse, nil
}
