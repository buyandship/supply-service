package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

// TransactionSearchRequest represents a transaction search request
type TransactionSearchRequest struct {
	YahooAccountID string `json:"yahoo_account_id"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	Status         string `json:"status,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
}

// SearchTransactions searches for transactions
func (c *Client) SearchTransactions(ctx context.Context, req TransactionSearchRequest) (*http.Response, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)

	if req.StartDate != "" {
		params.Set("start_date", req.StartDate)
	}
	if req.EndDate != "" {
		params.Set("end_date", req.EndDate)
	}
	if req.Status != "" {
		params.Set("status", req.Status)
	}
	if req.Limit > 0 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.Offset > 0 {
		params.Set("offset", strconv.Itoa(req.Offset))
	}

	return c.makeRequest(ctx, "GET", "/api/v1/transactions", params, nil, AuthTypeHMAC)
}

func (c *Client) MockGetTransaction(ctx context.Context, req *supply.YahooGetTransactionReq, yahooAccountID string) (*Transaction, error) {
	resp := Transaction{
		TransactionID:   req.TransactionID,
		YsRefID:         "YS-REF-001",
		AuctionID:       "x12345",
		CurrentPrice:    1000,
		TransactionType: "BID",
		Status:          "completed",
		ReqPrice:        1000,
		CreatedAt:       "2025-10-22T12:00:00Z",
		UpdatedAt:       "2025-10-22T12:00:01Z",
	}
	return &resp, nil
}
