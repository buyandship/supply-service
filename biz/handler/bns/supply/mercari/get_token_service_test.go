package mercari

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
)

func TestGetToken_Concurrent(t *testing.T) {
	// Load test config

	// Create HTTP client
	c, err := client.NewClient()
	assert.NoError(t, err)

	ctx := context.Background()
	numGoroutines := 10
	var wg sync.WaitGroup
	results := make(chan *struct {
		token string
		err   error
	}, numGoroutines)

	// Launch multiple goroutines to call GetToken endpoint concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resp supply.MercariGetTokenResp
			req := &protocol.Request{}
			req.SetRequestURI("http://mkp-t2.buynship.com/v1/supplysrv/internal/mercari/token")
			req.Header.SetMethod(consts.MethodGet)
			req.Header.Set("timestamp", "1742969664845")
			req.Header.Set("hmac", "bc94e449adcf48949a384c1385817958eaa96260dc0010efefd604226a5a9183")
			req.Header.Set("X-Request-ID", fmt.Sprintf("TEST_CASE_%d", i))
			res := &protocol.Response{}
			err := c.Do(ctx, req, res)
			if err != nil {
				fmt.Println("err", err)
				results <- &struct {
					token string
					err   error
				}{"", err}
				return
			}
			if res.StatusCode() != 200 {
				results <- &struct {
					token string
					err   error
				}{"", bizErr.InternalError}
				return
			}

			// Decode response body
			if err := json.Unmarshal(res.Body(), &resp); err != nil {
				fmt.Println("decode err", err)
				results <- &struct {
					token string
					err   error
				}{"", err}
				return
			}

			results <- &struct {
				token string
				err   error
			}{resp.Token, nil}
		}()
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(results)

	// Collect results
	var successCount, conflictCount, rateLimitCount int
	var tokens []string
	for result := range results {
		if result.err != nil {
			if result.err == bizErr.ConflictError {
				conflictCount++
			} else if result.err == bizErr.RateLimitError {
				rateLimitCount++
			}
			continue
		}
		successCount++
		tokens = append(tokens, result.token)
	}

	// Verify results
	t.Logf("Test results: success=%d, conflict=%d, rate_limit=%d", successCount, conflictCount, rateLimitCount)

	// Check that we got at least one successful token
	assert.Greater(t, successCount, 0, "Should have at least one successful token refresh")

	// Check that all successful tokens are the same (indicating lock worked)
	if len(tokens) > 1 {
		firstToken := tokens[0]
		for _, token := range tokens[1:] {
			assert.Equal(t, firstToken, token, "All successful tokens should be the same")
		}
	}

	// Check that we got some rate limit errors (indicating rate limiting worked)
	assert.Greater(t, rateLimitCount, 0, "Should have some rate limit errors")
}

func TestGetToken_ConcurrentWithTimeout(t *testing.T) {

	// Create HTTP client
	c, err := client.NewClient()
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	numGoroutines := 5
	var wg sync.WaitGroup
	results := make(chan *struct {
		token string
		err   error
	}, numGoroutines)

	// Launch goroutines with small delays between them
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(delay time.Duration) {
			defer wg.Done()
			time.Sleep(delay)
			var resp supply.MercariGetTokenResp
			req := &protocol.Request{}
			req.SetRequestURI("http://mkp-t2.buynship.com/v1/supplysrv/internal/mercari/token")
			req.Header.SetMethod(consts.MethodGet)
			req.Header.Set("timestamp", "1742969664845")
			req.Header.Set("hmac", "bc94e449adcf48949a384c1385817958eaa96260dc0010efefd604226a5a9183")
			req.Header.Set("X-Request-ID", "TEST_CASE")
			res := &protocol.Response{}
			err := c.Do(ctx, req, res)
			if err != nil {
				results <- &struct {
					token string
					err   error
				}{"", err}
				return
			}
			if res.StatusCode() != 200 {
				results <- &struct {
					token string
					err   error
				}{"", bizErr.InternalError}
				return
			}

			// Decode response body
			if err := json.Unmarshal(res.Body(), &resp); err != nil {
				results <- &struct {
					token string
					err   error
				}{"", err}
				return
			}

			results <- &struct {
				token string
				err   error
			}{resp.Token, nil}
		}(time.Duration(i) * 100 * time.Millisecond)
	}

	// Wait for all goroutines to complete or context to timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		t.Log("Test timed out as expected")
	case <-done:
		t.Log("All goroutines completed before timeout")
	}

	close(results)

	// Collect results
	var successCount, conflictCount, rateLimitCount int
	for result := range results {
		if result.err != nil {
			if result.err == bizErr.ConflictError {
				conflictCount++
			} else if result.err == bizErr.RateLimitError {
				rateLimitCount++
			}
			continue
		}
		successCount++
	}

	// Log results
	t.Logf("Test results: success=%d, conflict=%d, rate_limit=%d", successCount, conflictCount, rateLimitCount)

	// Verify we got some successful responses
	assert.Greater(t, successCount, 0, "Should have at least one successful token refresh")
}

func TestGetToken_Sequential(t *testing.T) {
	// Create HTTP client
	c, err := client.NewClient()
	assert.NoError(t, err)

	ctx := context.Background()

	// Make first request
	var resp1 supply.MercariGetTokenResp
	req1 := &protocol.Request{}
	req1.SetRequestURI("http://mkp-t2.buynship.com/v1/supplysrv/internal/mercari/token")
	req1.Header.SetMethod(consts.MethodGet)
	req1.Header.Set("timestamp", "1742969664845")
	req1.Header.Set("hmac", "bc94e449adcf48949a384c1385817958eaa96260dc0010efefd604226a5a9183")
	req1.Header.Set("X-Request-ID", "TEST_CASE")
	res1 := &protocol.Response{}
	err1 := c.Do(ctx, req1, res1)
	assert.NoError(t, err1)
	assert.Equal(t, 200, res1.StatusCode())

	// Decode first response
	if err := json.Unmarshal(res1.Body(), &resp1); err != nil {
		t.Fatalf("Failed to decode first response: %v", err)
	}
	assert.NotEmpty(t, resp1.Token)

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Make second request
	var resp2 supply.MercariGetTokenResp
	req2 := &protocol.Request{}
	req2.SetRequestURI("http://mkp-t2.buynship.com/v1/supplysrv/internal/mercari/token")
	req2.Header.SetMethod(consts.MethodGet)
	req2.Header.Set("timestamp", "1742969664845")
	req2.Header.Set("hmac", "bc94e449adcf48949a384c1385817958eaa96260dc0010efefd604226a5a9183")
	req2.Header.Set("X-Request-ID", "TEST_CASE")
	res2 := &protocol.Response{}
	err2 := c.Do(ctx, req2, res2)
	assert.NoError(t, err2)
	assert.Equal(t, 200, res2.StatusCode())

	// Decode second response
	if err := json.Unmarshal(res2.Body(), &resp2); err != nil {
		t.Fatalf("Failed to decode second response: %v", err)
	}
	assert.NotEmpty(t, resp2.Token)

	// Verify tokens are the same (cached)
	assert.Equal(t, resp1.Token, resp2.Token, "Tokens should be the same due to caching")
}
