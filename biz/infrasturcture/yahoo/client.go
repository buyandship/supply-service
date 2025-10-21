package yahoo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Configuration for Yahoo! Auction Bridge Service
const (
	APIKey    = "buyship-service-001"
	SecretKey = "your-secret-key-here"
	BaseURL   = "https://yahoo-auction-bridge.example.com"
)

// Client represents the Yahoo Auction Bridge client
type Client struct {
	apiKey     string
	secretKey  string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Yahoo Auction Bridge client
func NewClient() *Client {
	// load from config
	return &Client{
		apiKey:    APIKey,
		secretKey: SecretKey,
		baseURL:   BaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeAuthenticatedRequest makes an authenticated request to Yahoo! Auction Bridge Service
func (c *Client) makeAuthenticatedRequest(method, path string, body interface{}) (*http.Response, error) {
	// 1. Generate timestamp
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 2. Convert body to JSON string (NO SPACES)
	var bodyStr string
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyStr = string(bodyBytes)
	}

	// 3. Create request
	url := c.baseURL + path
	var bodyReader io.Reader
	if bodyStr != "" {
		bodyReader = bytes.NewBufferString(bodyStr)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	req.Header.Set("X-Timestamp", timestamp)

	// 4. Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	return resp, nil
}

// GetItem retrieves product information
func (c *Client) GetItem(aID string) (*http.Response, error) {
	path := fmt.Sprintf("/api/v1/auction/item?aID=%s", aID)
	return c.makeAuthenticatedRequest("GET", path, nil)
}

// PurchaseItem performs a buy-out purchase
func (c *Client) PurchaseItem(aID string, quantity int) (*http.Response, error) {
	path := "/api/v1/auction/purchase"
	body := map[string]interface{}{
		"aID":      aID,
		"Quantity": quantity,
	}
	return c.makeAuthenticatedRequest("POST", path, body)
}

// SearchItems searches for products
func (c *Client) SearchItems(query, category string) (*http.Response, error) {
	path := fmt.Sprintf("/api/v1/auction/search?query=%s&category=%s", query, category)
	return c.makeAuthenticatedRequest("GET", path, nil)
}
