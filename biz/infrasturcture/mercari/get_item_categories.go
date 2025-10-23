package mercari

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/buyandship/bns-golib/cache"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type GetItemCategoriesResp struct {
	MasterCategories []struct {
		Id       string `json:"id"`
		Name     string `json:"name"`
		Level    string `json:"level,omitempty"`
		ParentId string `json:"parent_id,omitempty"`
	} `json:"master_categories"`
}

func (m *Mercari) GetItemCategories(ctx context.Context) (*GetItemCategoriesResp, error) {
	getItemFunc := func() (*GetItemCategoriesResp, error) {

		token, err := m.GetActiveToken(ctx)
		if err != nil {
			return nil, err
		}

		if ok := cache.GetRedisClient().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}
		headers := map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {token.AccessToken},
		}

		url := fmt.Sprintf("%s/v1/master/item_categories", m.OpenApiDomain)
		httpReq, err := http.NewRequest("GET",
			url, nil)
		if err != nil {
			hlog.CtxErrorf(ctx, "http request error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		httpReq.Header = headers

		httpRes, err := m.Client.Do(ctx, httpReq)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		defer func() {
			if err := httpRes.Body.Close(); err != nil {
				hlog.CtxErrorf(ctx, "http close error: %s", err)
			}
		}()

		if httpRes.StatusCode == http.StatusUnauthorized {
			hlog.CtxInfof(ctx, "http unauthorized, refreshing token...")
			if err := m.RefreshToken(ctx, token); err != nil {
				hlog.CtxErrorf(ctx, "try to refresh token, but fails, err: %v", err)
				return nil, backoff.RetryAfter(1)
			}
			return nil, bizErr.UnauthorisedError
		}
		// retry code: 409, 429, 5xx
		if httpRes.StatusCode == http.StatusTooManyRequests {
			hlog.CtxWarnf(ctx, "http too many requests, retrying...")
			return nil, backoff.RetryAfter(1)
		}
		if httpRes.StatusCode == http.StatusConflict {
			hlog.CtxWarnf(ctx, "http conflict, retrying...")
			return nil, bizErr.ConflictError
		}
		if httpRes.StatusCode >= 500 && httpRes.StatusCode < 600 {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxInfof(ctx, "http error, error_code: [%d], error_msg: [%s], retrying at [%+v]...",
				httpRes.StatusCode, respBody, time.Now().Local())
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			}
		}

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxInfof(ctx, "get mercari item categories error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &GetItemCategoriesResp{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		return resp, nil
	}
	result, err := backoff.Retry(ctx, getItemFunc, m.GetRetryOpts()...)
	if err != nil {
		pErr := &backoff.PermanentError{}
		if errors.As(err, &pErr) {
			hlog.CtxInfof(ctx, "get mercari item categories error: %v", err)
			berr := pErr.Unwrap()
			return nil, berr
		}
		hlog.CtxInfof(ctx, "get mercari item categories error: %v", err)
		return nil, err
	}
	return result, nil
}
