package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

// CategoryLeafOption represents option information in category leaf response
// Extends AuctionItemListOption with additional icon fields
type CategoryLeafOption struct {
	NewIcon           string `json:"NewIcon,omitempty"`
	StoreIcon         string `json:"StoreIcon,omitempty"`
	FeaturedIcon      string `json:"FeaturedIcon,omitempty"`
	FreeshippingIcon  string `json:"FreeshippingIcon,omitempty"`
	NewItemIcon       string `json:"NewItemIcon,omitempty"`
	WrappingIcon      string `json:"WrappingIcon,omitempty"`
	BuynowIcon        string `json:"BuynowIcon,omitempty"`
	EasyPaymentIcon   string `json:"EasyPaymentIcon,omitempty"`
	PointIcon         string `json:"PointIcon,omitempty"`
	CharityOptionIcon string `json:"CharityOptionIcon,omitempty"`
	IsBold            bool   `json:"IsBold,omitempty"`
	IsBackGroundColor bool   `json:"IsBackGroundColor,omitempty"`
	IsOffer           bool   `json:"IsOffer,omitempty"`
	IsCharity         bool   `json:"IsCharity,omitempty"`
}

// CategoryLeafItem represents an auction item in the category leaf response
type CategoryLeafItem struct {
	AuctionID        string                `json:"AuctionID,omitempty"`
	Title            string                `json:"Title,omitempty"`
	Seller           AuctionItemListSeller `json:"Seller,omitempty"`
	ItemUrl          string                `json:"ItemUrl,omitempty"`
	AuctionItemUrl   string                `json:"AuctionItemUrl,omitempty"`
	Image            AuctionImage          `json:"Image,omitempty"`
	OriginalImageNum int                   `json:"OriginalImageNum,omitempty"`
	CurrentPrice     float64               `json:"CurrentPrice,omitempty"`
	Bids             int                   `json:"Bids,omitempty"`
	StartTime        string                `json:"StartTime,omitempty"`
	EndTime          string                `json:"EndTime,omitempty"`
	BidOrBuy         float64               `json:"BidOrBuy,omitempty"`
	IsReserved       bool                  `json:"IsReserved,omitempty"`
	CharityOption    CharityOption         `json:"CharityOption,omitempty"`
	Affiliate        Affiliate             `json:"Affiliate,omitempty"`
	Option           CategoryLeafOption    `json:"Option,omitempty"`
	IsAdult          bool                  `json:"IsAdult,omitempty"`
}

// CategoryLeafResult represents the result object in category leaf response
type CategoryLeafResult struct {
	CategoryPath   string             `json:"CategoryPath,omitempty"`
	CategoryIdPath string             `json:"CategoryIdPath,omitempty"`
	Item           []CategoryLeafItem `json:"Item,omitempty"`
}

// CategoryLeafResponse represents the response from category leaf API endpoint
type CategoryLeafResponse struct {
	ResultSet struct {
		TotalResultsAvailable int                `json:"@totalResultsAvailable"`
		TotalResultsReturned  int                `json:"@totalResultsReturned"`
		FirstResultPosition   int                `json:"@firstResultPosition"`
		Result                CategoryLeafResult `json:"Result"`
	} `json:"ResultSet"`
}

func (c *Client) GetCategoryLeaf(ctx context.Context, req *supply.YahooGetCategoryLeafReq) (*CategoryLeafResponse, error) {
	params := url.Values{}
	params.Set("category", strconv.Itoa(int(req.Category)))
	if req.ExceptCategory != nil {
		params.Set("except_category", *req.ExceptCategory)
	}
	if req.Featured != nil {
		params.Set("featured", strconv.Itoa(int(*req.Featured)))
	}
	if req.Page != nil {
		params.Set("page", strconv.Itoa(int(*req.Page)))
	}
	if req.Sort != nil {
		params.Set("sort", *req.Sort)
	}
	if req.Order != nil {
		params.Set("order", *req.Order)
	}
	if req.Store != nil {
		params.Set("store", strconv.Itoa(int(*req.Store)))
	}
	if req.Aucminprice != nil {
		params.Set("aucminprice", strconv.Itoa(int(*req.Aucminprice)))
	}
	if req.Aucmaxprice != nil {
		params.Set("aucmaxprice", strconv.Itoa(int(*req.Aucmaxprice)))
	}
	if req.AucminBidorbuyPrice != nil {
		params.Set("aucmin_bidorbuy_price", strconv.Itoa(int(*req.AucminBidorbuyPrice)))
	}
	if req.AucmaxBidorbuyPrice != nil {
		params.Set("aucmax_bidorbuy_price", strconv.Itoa(int(*req.AucmaxBidorbuyPrice)))
	}
	if req.Easypayment != nil {
		params.Set("easypayment", strconv.Itoa(int(*req.Easypayment)))
	}
	if req.New != nil {
		params.Set("new", strconv.Itoa(int(*req.New)))
	}
	if req.Freeshipping != nil {
		params.Set("freeshipping", strconv.Itoa(int(*req.Freeshipping)))
	}
	if req.Wrappingicon != nil {
		params.Set("wrappingicon", strconv.Itoa(int(*req.Wrappingicon)))
	}
	if req.Buynow != nil {
		params.Set("buynow", strconv.Itoa(int(*req.Buynow)))
	}
	if req.Thumbnail != nil {
		params.Set("thumbnail", strconv.Itoa(int(*req.Thumbnail)))
	}
	if req.Attn != nil {
		params.Set("attn", strconv.Itoa(int(*req.Attn)))
	}
	if req.Point != nil {
		params.Set("point", strconv.Itoa(int(*req.Point)))
	}
	if req.ItemStatus != nil {
		params.Set("item_status", *req.ItemStatus)
	}
	if req.Adf != nil {
		params.Set("adf", strconv.Itoa(int(*req.Adf)))
	}
	if req.MinCharity != nil {
		params.Set("min_charity", strconv.Itoa(int(*req.MinCharity)))
	}
	if req.MaxCharity != nil {
		params.Set("max_charity", strconv.Itoa(int(*req.MaxCharity)))
	}
	if req.MinAffiliate != nil {
		params.Set("min_affiliate", strconv.Itoa(int(*req.MinAffiliate)))
	}
	if req.MaxAffiliate != nil {
		params.Set("max_affiliate", strconv.Itoa(int(*req.MaxAffiliate)))
	}
	if req.Timebuf != nil {
		params.Set("timebuf", strconv.Itoa(int(*req.Timebuf)))
	}
	if req.Ranking != nil {
		params.Set("ranking", strconv.Itoa(int(*req.Ranking)))
	}
	if req.SellerAucUserID != nil {
		params.Set("seller_auc_user_id", *req.SellerAucUserID)
	}
	if req.BlackSellerAucUserID != nil {
		params.Set("black_seller_auc_user_id", *req.BlackSellerAucUserID)
	}
	if req.Sort2 != nil {
		params.Set("sort2", *req.Sort2)
	}
	if req.Order2 != nil {
		params.Set("order2", *req.Order2)
	}
	if req.LocCd != nil {
		params.Set("loc_cd", *req.LocCd)
	}
	if req.Fixed != nil {
		params.Set("fixed", strconv.Itoa(int(*req.Fixed)))
	}
	if req.MaxStart != nil {
		params.Set("max_start", strconv.FormatInt(*req.MaxStart, 10))
	}
	if req.MinStart != nil {
		params.Set("min_start", strconv.FormatInt(*req.MinStart, 10))
	}
	if req.ExceptShoppingitem != nil {
		params.Set("except_shoppingitem", strconv.Itoa(int(*req.ExceptShoppingitem)))
	}
	if req.Callback != nil {
		params.Set("callback", *req.Callback)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/categoryLeaf", params, nil, AuthTypeNone)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnprocessableEntity:
				return nil, bizErr.BizError{
					Status:  http.StatusBadRequest,
					ErrCode: http.StatusBadRequest,
					ErrMsg:  "validation error",
				}
			case http.StatusBadRequest:
				return nil, bizErr.BizError{
					Status:  http.StatusNotFound,
					ErrCode: http.StatusNotFound,
					ErrMsg:  "category not found",
				}
			case http.StatusInternalServerError:
				return nil, bizErr.BizError{
					Status:  http.StatusInternalServerError,
					ErrCode: 10001,
					ErrMsg:  "internal server error",
				}
			}
		}
		return nil, bizErr.BizError{
			Status:  http.StatusInternalServerError,
			ErrCode: http.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}
	var httpResp CategoryLeafResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}
