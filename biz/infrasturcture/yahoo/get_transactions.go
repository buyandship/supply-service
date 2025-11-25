package yahoo

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

type Transaction struct {
	TransactionID    string  `json:"transaction_id,omitempty"`
	RequestGroupID   string  `json:"request_group_id,omitempty"`
	RetryCount       int     `json:"retry_count,omitempty"`
	YsRefID          string  `json:"ys_ref_id"`
	YahooAccountID   string  `json:"yahoo_account_id,omitempty"`
	AuctionID        string  `json:"auction_id,omitempty"`
	CurrentPrice     float64 `json:"current_price,omitempty"`
	TransactionType  string  `json:"transaction_type,omitempty"`
	Status           string  `json:"status,omitempty"`
	APIEndpoint      string  `json:"api_endpoint,omitempty"`
	HTTPStatus       int     `json:"http_status,omitempty"`
	ProcessingTimeMS int     `json:"processing_time_ms,omitempty"`
	ReqPrice         float64 `json:"req_price,omitempty"`
	CreatedAt        string  `json:"created_at,omitempty"`
	UpdatedAt        string  `json:"updated_at,omitempty"`
	// RequestData      supply.YahooPlaceBidReq `json:"request_data"`
	// ResponseData     PlaceBidResponse        `json:"response_data,omitempty"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Count        int           `json:"count"`
	NextCursor   string        `json:"next_cursor"`
}

// GetTransaction gets specific transaction details
func (c *Client) GetTransaction(ctx context.Context, req *supply.YahooGetTransactionReq, yahooAccountID string) (*Transaction, error) {
	path := fmt.Sprintf("/api/v1/transactions/%s", req.TransactionID)
	params := url.Values{}
	params.Set("yahoo_account_id", yahooAccountID)

	resp, err := c.makeRequest(ctx, "GET", path, params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusNotFound:
				return nil, bizErr.BizError{
					Status:  http.StatusNotFound,
					ErrCode: http.StatusNotFound,
					ErrMsg:  "Transaction not found",
				}
			case http.StatusUnprocessableEntity:
				return nil, bizErr.BizError{
					Status:  http.StatusBadRequest,
					ErrCode: http.StatusBadRequest,
					ErrMsg:  "invalid transaction id",
				}
			default:
				return nil, bizErr.BizError{
					Status:  http.StatusInternalServerError,
					ErrCode: http.StatusInternalServerError,
					ErrMsg:  "internal server error",
				}
			}
		}
		return nil, bizErr.BizError{
			Status:  http.StatusInternalServerError,
			ErrCode: http.StatusInternalServerError,
			ErrMsg:  err.Error(),
		}
	}

	var tx Transaction
	if err := c.parseResponse(resp, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}

func (c *Client) GetTransactions(ctx context.Context, req *supply.YahooGetTransactionsReq) (*GetTransactionsResponse, error) {
	path := "/api/v1/transactions"
	params := url.Values{}
	params.Set("yahoo_account_id", config.Yahoo02AccountID)
	if req.TransactionID != "" {
		params.Set("transaction_id", req.TransactionID)
	}
	if req.YsRefID != "" {
		params.Set("ys_ref_id", req.YsRefID)
	}
	if req.AuctionID != "" {
		params.Set("auction_id", req.AuctionID)
	}

	resp, err := c.makeRequest(ctx, "GET", path, params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	var getTransactionsResponse GetTransactionsResponse
	if err := c.parseResponse(resp, &getTransactionsResponse); err != nil {
		return nil, err
	}

	return &getTransactionsResponse, nil
}
