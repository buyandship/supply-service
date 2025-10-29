package yahoo

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"sync"
	"time"

	"net/http"

	bnsHttp "github.com/buyandship/bns-golib/http"
	"github.com/buyandship/bns-golib/retry"
	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var (
	once    sync.Once
	Handler *Client
)

// Client represents the Yahoo Auction Bridge API client
type Client struct {
	baseURL    string
	apiKey     string
	secretKey  string
	httpClient *bnsHttp.Client
}

// NewClient creates a new Yahoo Auction Bridge client
func GetClient() *Client {
	once.Do(func() {
		client := bnsHttp.NewClient(
			bnsHttp.WithTimeout(10 * time.Second), // TODO: change to actual timeout
		)
		var baseURL string
		switch config.GlobalServerConfig.Env {
		case "development":
			baseURL = "https://internal-stagin20251027043053843000000001-645109195.ap-northeast-1.elb.amazonaws.com" // TODO: change to actual url
		case "production":
			baseURL = "https://mock-api.yahoo-auction.jp" // TODO: change to actual url
		}
		apiKey := config.GlobalServerConfig.Yahoo.ApiKey
		secretKey := config.GlobalServerConfig.Yahoo.SecretKey
		Handler = &Client{
			baseURL:    baseURL,
			apiKey:     apiKey,
			secretKey:  secretKey,
			httpClient: client,
		}
	})
	return Handler
}

// Authentication types
type AuthType string

const (
	AuthTypeHMAC  AuthType = "hmac"
	AuthTypeOAuth AuthType = "oauth"
	AuthTypeNone  AuthType = "none"
)

// Request/Response Models

// PlaceBidRequest represents a bid request
type PlaceBidRequest struct {
	YahooAccountID  string `json:"yahoo_account_id"`
	YsRefID         string `json:"ys_ref_id"`
	TransactionType string `json:"transaction_type"` // BID or BUYOUT
	AuctionID       string `json:"auction_id"`
	Price           int    `json:"price"`
	Signature       string `json:"signature"`
	Quantity        int    `json:"quantity,omitempty"`
	Partial         bool   `json:"partial,omitempty"`
}

type PlaceBidResponse struct {
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

type PlaceBidPreviewResponse struct {
}

// AuctionItemRequest represents a request for auction item information
type AuctionItemRequest struct {
	AuctionID string `json:"auctionID"`
	AppID     string `json:"appid,omitempty"`
}

// TransactionSearchRequest represents a transaction search request
type TransactionSearchRequest struct {
	YahooAccountID string `json:"yahoo_account_id"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	Status         string `json:"status,omitempty"`
	Limit          int    `json:"limit,omitempty"`
	Offset         int    `json:"offset,omitempty"`
}

// OAuthAuthorizeRequest represents OAuth authorization request
type OAuthAuthorizeRequest struct {
	YahooAccountID string `json:"yahoo_account_id"`
}

// TokenRefreshRequest represents token refresh request
type TokenRefreshRequest struct {
	YahooAccountID string `json:"yahoo_account_id"`
}

// Response models
type AuctionItemResponse struct {
	ResultSet struct {
		Result struct {
			AuctionID    string `json:"AuctionID"`
			Title        string `json:"Title"`
			Description  string `json:"Description"`
			CurrentPrice int    `json:"CurrentPrice"`
			StartPrice   int    `json:"StartPrice"`
			Bids         int    `json:"Bids"`
			ItemStatus   string `json:"ItemStatus"`
			EndTime      string `json:"EndTime"`
			StartTime    string `json:"StartTime"`
			Seller       Seller `json:"Seller"`
			Image        string `json:"Image"`
		} `json:"Result"`
	} `json:"ResultSet"`
}

type Seller struct {
	ID     string  `json:"Id"`
	Rating float64 `json:"Rating"`
}

type ErrorResponse struct {
	Detail []ErrorDetail `json:"detail"`
}

type ErrorDetail struct {
	Type  string      `json:"type"`
	Loc   []string    `json:"loc"`
	Msg   string      `json:"msg"`
	Input interface{} `json:"input,omitempty"`
}

// Helper method to generate HMAC signature
func (c *Client) generateHMACSignature(timestamp, method, path, body string) string {
	message := timestamp + method + path + body
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

// Helper method to make authenticated requests
func (c *Client) makeRequest(ctx context.Context, method, path string, params url.Values, body interface{}, authType AuthType) (*http.Response, error) {
	// Build URL
	fullURL := c.baseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	// Prepare request body
	var bodyReader io.Reader
	var bodyStr string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyStr = string(bodyBytes)
		bodyReader = bytes.NewBufferString(bodyStr)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication headers
	if authType == AuthTypeHMAC {
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		signature := c.generateHMACSignature(timestamp, method, path, bodyStr)

		req.Header.Set("X-API-Key", c.apiKey)
		req.Header.Set("X-Timestamp", timestamp)
		req.Header.Set("X-Signature", signature)
	}

	// Set content type for POST requests
	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Make request with retry mechanism
	var resp *http.Response
	operation := func() (*http.Response, error) {
		var err error
		resp, err = c.httpClient.Do(ctx, req)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				hlog.CtxErrorf(ctx, "http close error: %s", err)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			hlog.CtxWarnf(ctx, "http error, error_code: [%d], error_msg: [%s]",
				resp.StatusCode, string(respBody))
		}

		// TODO: retrable error

		if resp.StatusCode >= 500 {
			// non-retryable error
			return nil, backoff.Permanent(fmt.Errorf("server error: %d", resp.StatusCode))
		}

		return resp, nil
	}

	resp, err = backoff.Retry(ctx, operation, retry.GetDefaultRetryOpts()...)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to send request after retries: %v", err)
		return nil, fmt.Errorf("failed to send request after retries: %w", err)
	}

	return resp, nil
}

// API Methods
func (c *Client) parseResponse(resp *http.Response, v interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}

// Authorize initiates Yahoo OAuth 2.0 authorization flow
func (c *Client) Authorize(ctx context.Context, req OAuthAuthorizeRequest) (*http.Response, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)

	return c.makeRequest(ctx, "GET", "/auth/authorize", params, nil, AuthTypeNone)
}

// OAuthCallback handles OAuth callback
func (c *Client) OAuthCallback(ctx context.Context, code, state string) (*http.Response, error) {
	params := url.Values{}
	params.Set("code", code)
	params.Set("state", state)

	return c.makeRequest(ctx, "GET", "/auth/callback", params, nil, AuthTypeNone)
}

// GetTokenStatus gets OAuth token status
func (c *Client) GetTokenStatus(ctx context.Context, yahooAccountID string) (*http.Response, error) {
	path := fmt.Sprintf("/auth/token/status/%s", yahooAccountID)
	return c.makeRequest(ctx, "GET", path, nil, nil, AuthTypeNone)
}

// RefreshToken refreshes OAuth token
func (c *Client) RefreshToken(ctx context.Context, req TokenRefreshRequest) (*http.Response, error) {
	path := fmt.Sprintf("/auth/token/refresh/%s", req.YahooAccountID)
	return c.makeRequest(ctx, "POST", path, nil, nil, AuthTypeNone)
}

// RevokeToken revokes OAuth token
func (c *Client) RevokeToken(ctx context.Context, yahooAccountID string) (*http.Response, error) {
	path := fmt.Sprintf("/auth/token/%s", yahooAccountID)
	return c.makeRequest(ctx, "DELETE", path, nil, nil, AuthTypeNone)
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

	return &placeBidPreviewResponse, nil
}

// PlaceBid executes a bid on Yahoo Auction
func (c *Client) PlaceBid(ctx context.Context, req *PlaceBidRequest) (*supply.YahooPlaceBidResp, error) {
	params := url.Values{}
	params.Set("yahoo_account_id", req.YahooAccountID)
	params.Set("ys_ref_id", req.YsRefID)
	params.Set("transaction_type", req.TransactionType)
	params.Set("auction_id", req.AuctionID)
	params.Set("price", strconv.Itoa(req.Price))
	params.Set("signature", req.Signature)

	if req.Quantity > 0 {
		params.Set("quantity", strconv.Itoa(req.Quantity))
	}
	if req.Partial {
		params.Set("partial", "true")
	}

	resp, err := c.makeRequest(ctx, "POST", "/api/v1/placeBid", params, nil, AuthTypeHMAC)
	if err != nil {
		return nil, err
	}

	placeBidResponse := supply.YahooPlaceBidResp{}
	if err := c.parseResponse(resp, &placeBidResponse); err != nil {
		return nil, err
	}

	return &placeBidResponse, nil
}

func (c *Client) MockPlaceBid(ctx context.Context, req *PlaceBidRequest) (*supply.YahooPlaceBidResp, error) {
	return nil, nil
}

// GetAuctionItem gets auction item information (public API)
func (c *Client) GetAuctionItem(ctx context.Context, req AuctionItemRequest) (*http.Response, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	return c.makeRequest(ctx, "GET", "/api/v1/auctionItem", params, nil, AuthTypeNone)
}

// GetAuctionItemAuth gets authenticated auction item information
func (c *Client) GetAuctionItemAuth(ctx context.Context, req AuctionItemRequest, yahooAccountID string) (*http.Response, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	params.Set("yahoo_account_id", yahooAccountID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	return c.makeRequest(ctx, "GET", "/api/v1/auctionItemAuth", params, nil, AuthTypeHMAC)
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

// GetTransaction gets specific transaction details
func (c *Client) GetTransaction(ctx context.Context, transactionID, yahooAccountID string) (*http.Response, error) {
	path := fmt.Sprintf("/api/v1/transactions/%s", transactionID)
	params := url.Values{}
	params.Set("yahoo_account_id", yahooAccountID)

	return c.makeRequest(ctx, "GET", path, params, nil, AuthTypeHMAC)
}

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

// Health check
func (c *Client) HealthCheck() (*http.Response, error) {
	return c.makeRequest(context.Background(), "GET", "/health", nil, nil, AuthTypeNone)
}

// Get API info
func (c *Client) GetAPIInfo() (*http.Response, error) {
	return c.makeRequest(context.Background(), "GET", "/", nil, nil, AuthTypeNone)
}

// Helper method to parse error response
func ParseErrorResponse(resp *http.Response) (*ErrorResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var errorResp ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return nil, fmt.Errorf("failed to parse error response: %w", err)
	}

	return &errorResp, nil
}

// Helper method to parse auction item response
func ParseAuctionItemResponse(resp *http.Response) (*AuctionItemResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var auctionResp AuctionItemResponse
	if err := json.Unmarshal(body, &auctionResp); err != nil {
		return nil, fmt.Errorf("failed to parse auction item response: %w", err)
	}

	return &auctionResp, nil
}
