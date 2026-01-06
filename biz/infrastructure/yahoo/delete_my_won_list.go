package yahoo

import (
	"context"
	"net/http"
	"net/url"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

type DeleteMyWonListResponse struct {
	ResultSet struct {
		TotalResultsAvailable int `json:"@totalResultsAvailable"`
		TotalResultsReturned  int `json:"@totalResultsReturned"`
		FirstResultPosition   int `json:"@firstResultPosition"`
		DeleteResultSet       struct {
			Result struct {
				AuctionID string `json:"auctionID,omitempty"`
				Result    bool   `json:"result,omitempty"`
			} `json:"Result,omitempty"`
			MyWonListURL string `json:"MyWonListUrl,omitempty"`
		} `json:"DeleteResultSet,omitempty"`
	} `json:"ResultSet,omitempty"`
}

func (c *Client) DeleteMyWonList(ctx context.Context, req *supply.YahooDeleteMyWonListReq) (*DeleteMyWonListResponse, error) {
	params := url.Values{}

	// TODO: get yahoo account
	// params.Set("yahoo_account_id", config.DevYahoo02AccountID)
	params.Set("ys_ref_id", req.YsRefID)
	params.Set("auction_id", req.AuctionID)

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/deleteMyWonList", params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnauthorized:
				return nil, bizErr.BizError{
					Status:  http.StatusUnauthorized,
					ErrCode: http.StatusUnauthorized,
					ErrMsg:  "unauthorized",
				}
			case http.StatusInternalServerError:
				return nil, bizErr.BizError{
					Status:  http.StatusInternalServerError,
					ErrCode: http.StatusInternalServerError,
					ErrMsg:  "internal server error",
				}
			case http.StatusUnprocessableEntity:
				return nil, bizErr.BizError{
					Status:  http.StatusUnprocessableEntity,
					ErrCode: http.StatusUnprocessableEntity,
					ErrMsg:  "field required",
				}
			}
		}
	}

	var httpResp DeleteMyWonListResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}
