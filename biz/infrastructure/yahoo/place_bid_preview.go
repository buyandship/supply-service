package yahoo

import (
	"context"
	"net/url"
	"strconv"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type PlaceBidPreviewResponse struct {
	ResultSet struct {
		Result                PlaceBidResult `json:"Result"`
		TotalResultsAvailable int            `json:"totalResultsAvailable"`
		TotalResultsReturned  int            `json:"totalResultsReturned"`
		FirstResultPosition   int            `json:"firstResultPosition"`
	} `json:"ResultSet,omitempty"`
	Error Error `json:"Error,omitempty"`
}

// PlaceBidPreviewRequest represents a bid preview request
type PlaceBidPreviewRequest struct {
	YahooAccountID  string `json:"yahoo_account_id"`
	YsRefID         string `json:"ys_ref_id"`
	TransactionType string `json:"transaction_type"`
	AuctionID       string `json:"auction_id"`
	Price           int    `json:"price"`
	Quantity        int    `json:"quantity,omitempty"`
	Partial         bool   `json:"partial,omitempty"`
}

// PlaceBidPreview gets bid preview with signature
func (c *Client) PlaceBidPreview(ctx context.Context, req *PlaceBidPreviewRequest) (*PlaceBidPreviewResponse, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)
	params.Set("ys_ref_id", req.YsRefID)
	params.Set("transaction_type", req.TransactionType)
	params.Set("auction_id", req.AuctionID)
	params.Set("price", strconv.Itoa(req.Price))

	if req.Quantity > 0 {
		params.Set("quantity", strconv.Itoa(req.Quantity))
	}
	if req.Partial {
		params.Set("partial", "true")
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/placeBidPreview", params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	var placeBidPreviewResponse PlaceBidPreviewResponse
	if err := c.parseResponse(resp, &placeBidPreviewResponse); err != nil {
		return nil, err
	}

	if placeBidPreviewResponse.Error.Code != 0 {
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "You were not able to place your bid. This auction has already ended.",
		}
	}

	return &placeBidPreviewResponse, nil
}

func (c *Client) MockPlaceBidPreview(ctx context.Context, req *PlaceBidPreviewRequest) (*PlaceBidPreviewResponse, error) {
	placeBidPreviewResponse := PlaceBidPreviewResponse{
		ResultSet: struct {
			Result                PlaceBidResult `json:"Result"`
			TotalResultsAvailable int            `json:"totalResultsAvailable"`
			TotalResultsReturned  int            `json:"totalResultsReturned"`
			FirstResultPosition   int            `json:"firstResultPosition"`
		}{
			Result: PlaceBidResult{
				Signature: "abc123def456...",
			},
		},
	}
	return &placeBidPreviewResponse, nil
}
