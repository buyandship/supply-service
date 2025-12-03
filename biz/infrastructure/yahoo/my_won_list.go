package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/buyandship/supply-service/biz/common/config"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

// MyWonListImage represents image information in my won list response
type MyWonListImage struct {
	URL    string `json:"Url,omitempty"`
	Width  int    `json:"@Width,omitempty"`
	Height int    `json:"@Height,omitempty"`
}

// MyWonListOption represents option information in my won list response
type MyWonListOption struct {
	StoreIconUrl        string `json:"StoreIconUrl,omitempty"`
	FreeShippingIconUrl string `json:"FreeShippingIconUrl,omitempty"`
	IsBold              bool   `json:"IsBold,omitempty"`
	IsBackGroundColor   bool   `json:"IsBackGroundColor,omitempty"`
}

// MyWonListSeller represents seller information in my won list response
type MyWonListSeller struct {
	AucUserId            string `json:"AucUserId,omitempty"`
	AucUserIdItemListURL string `json:"AucUserIdItemListUrl,omitempty"`
}

// MyWonListItem represents a won auction item in the my won list response
type MyWonListItem struct {
	AuctionID      string          `json:"AuctionID,omitempty"`
	Title          string          `json:"Title,omitempty"`
	WonPrice       float64         `json:"WonPrice,omitempty"`
	Bids           int             `json:"Bids,omitempty"`
	EndTime        string          `json:"EndTime,omitempty"`
	Seller         MyWonListSeller `json:"Seller,omitempty"`
	AuctionItemUrl string          `json:"AuctionItemUrl,omitempty"`
	Image          MyWonListImage  `json:"Image,omitempty"`
	Option         MyWonListOption `json:"Option,omitempty"`
}

// MyWonListResponse represents the response from my won list API endpoint
type MyWonListResponse struct {
	ResultSet struct {
		TotalResultsAvailable int             `json:"@totalResultsAvailable"`
		TotalResultsReturned  int             `json:"@totalResultsReturned"`
		FirstResultPosition   int             `json:"@firstResultPosition"`
		Result                []MyWonListItem `json:"Result"`
	} `json:"ResultSet"`
}

// GetMyWonList retrieves won auction list for authenticated user
func (c *Client) GetMyWonList(ctx context.Context, req *supply.YahooGetMyWonListReq) (*MyWonListResponse, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", config.MasterYahooAccountID)
	if req.YsRefID != nil {
		params.Set("ys_ref_id", *req.YsRefID)
	}
	if req.Start != nil {
		params.Set("start", strconv.Itoa(int(*req.Start)))
	}
	if req.ContactProgress != nil {
		params.Set("contact_progress", *req.ContactProgress)
	}
	if req.AuctionID != nil {
		params.Set("auction_id", *req.AuctionID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/myWonList", params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, bizErr.BizError{
					Status:  http.StatusUnauthorized,
					ErrCode: http.StatusUnauthorized,
					ErrMsg:  "OAuth token invalid or expired",
				}
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
	var httpResp MyWonListResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}
