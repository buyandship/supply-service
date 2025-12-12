package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/shopspring/decimal"
)

// SellingListRating represents seller rating information in selling list response
type SellingListRating struct {
	Point       int  `json:"Point,omitempty"`
	IsSuspended bool `json:"IsSuspended,omitempty"`
	IsDeleted   bool `json:"IsDeleted,omitempty"`
}

// SellingListSeller represents seller information in selling list response
type SellingListSeller struct {
	AucUserId            string            `json:"AucUserId,omitempty"`
	AucUserIdItemListURL string            `json:"AucUserIdItemListUrl,omitempty"`
	AucUserIdAboutURL    string            `json:"AucUserIdAboutUrl,omitempty"`
	AucUserIdRatingURL   string            `json:"AucUserIdRatingUrl,omitempty"`
	Rating               SellingListRating `json:"Rating,omitempty"`
}

// SellingListItem represents an auction item in the selling list response
type SellingListItem struct {
	AuctionID      string                `json:"AuctionID,omitempty"`
	Title          string                `json:"Title,omitempty"`
	ItemUrl        string                `json:"ItemUrl,omitempty"`
	AuctionItemUrl string                `json:"AuctionItemUrl,omitempty"`
	Image          AuctionImage          `json:"Image,omitempty"`
	CurrentPrice   float64               `json:"CurrentPrice,omitempty"`
	Bids           int                   `json:"Bids,omitempty"`
	EndTime        string                `json:"EndTime,omitempty"`
	BidOrbuy       float64               `json:"BidOrbuy,omitempty"`
	IsReserved     bool                  `json:"IsReserved,omitempty"`
	CharityOption  CharityOption         `json:"CharityOption,omitempty"`
	Option         AuctionItemListOption `json:"Option,omitempty"`
}

func (s *SellingListItem) GetBuyoutPriceString() string {
	return decimal.NewFromFloat(s.BidOrbuy).StringFixed(0)
}

func (s *SellingListItem) GetBidPriceString() string {
	return decimal.NewFromFloat(s.CurrentPrice).StringFixed(0)
}

// SellingListResult represents the result object in selling list response
type SellingListResult struct {
	Seller SellerInfo        `json:"Seller,omitempty"`
	Item   []SellingListItem `json:"Item,omitempty"`
}

// SellingListResponse represents the response from selling list API endpoint
type SellingListResponse struct {
	ResultSet struct {
		TotalResultsAvailable int               `json:"@totalResultsAvailable"`
		TotalResultsReturned  int               `json:"@totalResultsReturned"`
		FirstResultPosition   int               `json:"@firstResultPosition"`
		Result                SellingListResult `json:"Result"`
	} `json:"ResultSet"`
}

// GetSellingList retrieves seller's auction item list
func (c *Client) GetSellingList(ctx context.Context, req *supply.YahooGetSellingListReq) (*SellingListResponse, error) {
	params := url.Values{}
	params.Set("sellerAucUserId", req.SellerAucUserId)

	if req.YsRefID != nil {
		params.Set("ys_ref_id", *req.YsRefID)
	}
	if req.Start != nil {
		params.Set("start", strconv.Itoa(int(*req.Start)))
	}
	if req.Status != nil {
		params.Set("status", *req.Status)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/sellingList", params, nil, AuthTypeNone)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnprocessableEntity:
				return nil, bizErr.BizError{
					Status:  http.StatusBadRequest,
					ErrCode: http.StatusBadRequest,
					ErrMsg:  "validation error",
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
	var httpResp SellingListResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}
