package mercari

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
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
		if ok := redis.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}
		headers := map[string][]string{
			"Accept":              {"application/json"},
			"Authorization":       {m.Token.AccessToken},
			"x-crossborder-fresh": {"true"},
		}

		httpReq, err := http.NewRequest("GET",
			fmt.Sprintf("%s/v1/master/item_categories", m.OpenApiDomain), nil)
		if err != nil {
			hlog.CtxErrorf(ctx, "http request error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		httpReq.Header = headers

		client := &http.Client{}
		httpRes, err := client.Do(httpReq)
		defer func() {
			if err := httpRes.Body.Close(); err != nil {
				hlog.CtxErrorf(ctx, "http close error: %s", err)
			}
		}()
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		if httpRes.StatusCode == http.StatusUnauthorized {
			hlog.CtxErrorf(ctx, "http unauthorized, refreshing token...")
			if err := m.RefreshToken(ctx); err != nil {
				hlog.CtxErrorf(ctx, "try to refresh token, but fails, err: %v", err)
			}
			return nil, bizErr.UnauthorisedError
		}
		// retry code: 409, 429, 5xx
		if httpRes.StatusCode == http.StatusTooManyRequests {
			hlog.CtxErrorf(ctx, "http too many requests, retrying...")
			return nil, backoff.RetryAfter(1)
		}
		if httpRes.StatusCode == http.StatusConflict {
			hlog.CtxErrorf(ctx, "http conflict, retrying...")
			return nil, bizErr.ConflictError
		}
		if httpRes.StatusCode >= 500 && httpRes.StatusCode < 600 {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxErrorf(ctx, "http error, error_code: [%d], error_msg: [%s], retrying...", httpRes.StatusCode, respBody)
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			}
		}

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxErrorf(ctx, "get mercari item error: %s", respBody)
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
			berr := pErr.Unwrap()
			return nil, berr
		}
		return nil, err
	}
	return result, nil
}
