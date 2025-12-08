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

	"github.com/buyandship/bns-golib/config"
	bnsHttp "github.com/buyandship/bns-golib/http"
	"github.com/buyandship/bns-golib/retry"
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
		switch config.GlobalAppConfig.Env {
		case "dev":
			baseURL = "http://staging.yahoo-bridge.internal" // TODO: change to actual url
			// baseURL = "https://internal-stagin20251027043053843000000001-645109195.ap-northeast-1.elb.amazonaws.com"
		case "prod":
			baseURL = "http://production.yahoo-bridge.internal" // TODO: change to actual url
		}
		apiKey := config.GlobalAppConfig.GetString("yahoo.api_key")
		secretKey := config.GlobalAppConfig.GetString("yahoo.secret_key")
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
	AuthTypeHMAC AuthType = "hmac"
	AuthTypeNone AuthType = "none"
)

type Error struct {
	Message string `json:"Message"`
	Code    int    `json:"Code"`
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
	req.Header.Set("Content-Type", "application/json")

	// Make request with retry mechanism
	var resp *http.Response
	operation := func() (*http.Response, error) {
		var err error
		resp, err = c.httpClient.Do(ctx, req)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(err)
		}

		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
			// TODO: handle retryable error
		default:
			respBody, _ := io.ReadAll(resp.Body)
			hlog.CtxInfof(ctx, "status code: [%d], response body: [%s]",
				resp.StatusCode, string(respBody))
			return resp, backoff.Permanent(fmt.Errorf("%s", string(respBody)))
		}
	}

	resp, err = backoff.Retry(ctx, operation, retry.GetDefaultRetryOpts()...)
	if err != nil {
		return resp, fmt.Errorf("failed to send request after retries: %w", err)
	}

	return resp, nil
}

// API Methods
func (c *Client) parseResponse(resp *http.Response, v interface{}) error {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			hlog.Errorf("http close error: %s", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	hlog.Debugf("http response body: %s", string(body))

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
