package mercari

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/assert"
)

func TestGetItem_Concurrent(t *testing.T) {
	// Load test config

	// Create HTTP client
	c, err := client.NewClient()
	assert.NoError(t, err)

	ctx := context.Background()
	numGoroutines := 2
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
			req.SetRequestURI("http://mkp-t2.buynship.com/v1/supplysrv/internal/mercari/item")
			req.Header.Set("Content-Type", "application/json")
			req.Header.SetMethod(consts.MethodGet)
			req.Header.Set("timestamp", "1742969664845")
			req.Header.Set("hmac", "bc94e449adcf48949a384c1385817958eaa96260dc0010efefd604226a5a9183")
			req.Header.Set("X-Request-ID", fmt.Sprintf("TEST_CASE_%d", i))
			req.SetBody([]byte(`{"item_id": "m79752279771"}`))
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
				fmt.Println("status code", res.StatusCode())
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

	// Verify results
	t.Logf("Test results: success=%d, conflict=%d, rate_limit=%d", successCount, conflictCount, rateLimitCount)

	// Check that we got at least one successful token
	assert.Greater(t, successCount, 0, "Should have at least one successful token refresh")
}
