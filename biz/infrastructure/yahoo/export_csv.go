package yahoo

import (
	"context"
	"net/http"
	"net/url"
)

// ExportTransactionsCSV exports transactions as CSV
func (c *Client) ExportTransactionsCSV(ctx context.Context, req TransactionSearchRequest) (*http.Response, error) {
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

	return c.makeRequest(ctx, "GET", "/api/v1/transactions/export/csv", params, nil, AuthTypeHMAC)
}
