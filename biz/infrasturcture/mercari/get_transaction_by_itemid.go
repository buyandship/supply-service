package mercari

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type GetTransactionByItemIDResponse struct {
	Id           string `json:"id"`
	Status       string `json:"status"`
	ItemId       string `json:"item_id"`
	SellerId     string `json:"seller_id"`
	Price        int64  `json:"price"`
	PaidPrice    int64  `json:"paid_price"`
	UpdatedTime  int    `json:"updated_time"`
	ShippingInfo struct {
		ShippingMethodName string `json:"shipping_method_name"`
		TrackingNumber     string `json:"tracking_number"`
		Status             string `json:"status"`
		BuyerShippingFee   string `json:"buyer_shipping_fee"`
	} `json:"shipping_info"`
	AccountID int32 `json:"account_id"`
}

func (m *Mercari) GetTransactionByItemID(ctx context.Context, itemId string, accountId int32) (*GetTransactionByItemIDResponse, error) {
	getItemFunc := func() (*GetTransactionByItemIDResponse, error) {
		token, err := m.GetToken(ctx, accountId)
		if err != nil {
			return nil, err
		}

		if ok := cache.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}
		headers := map[string][]string{
			"Authorization": {token.AccessToken},
		}

		url := fmt.Sprintf("%s/v2/transactions/%s", m.OpenApiDomain, itemId)
		httpReq, err := http.NewRequest("GET", url, nil)
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
				hlog.CtxWarnf(ctx, "try to refresh token, but fails, err: %v", err)
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
			hlog.CtxWarnf(ctx, "http error, error_code: [%d], error_msg: [%s], retrying at [%+v]...",
				httpRes.StatusCode, respBody, time.Now().Local())
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			}
		}

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxInfof(ctx, "get mercari transaction by itemid error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &GetTransactionByItemIDResponse{}
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
			hlog.CtxInfof(ctx, "get mercari transaction by itemid error: %v", err)
			berr := pErr.Unwrap()
			return nil, berr
		}
		hlog.CtxInfof(ctx, "get mercari transaction by itemid error: %v", err)
		return nil, err
	}
	result.AccountID = accountId
	return result, nil
}
