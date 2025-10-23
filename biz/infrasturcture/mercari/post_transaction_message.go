package mercari

import (
	"bytes"
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

type PostTransactionMessageRequest struct {
	TransactionId string `json:"transactionID"`
	Message       string `json:"message"`
	AccountID     int32  `json:"account_id"`
}

type PostTransactionMessageResponse struct {
	Id        string `json:"id"`
	Body      string `json:"body"`
	UserId    string `json:"user_id"`
	Created   int    `json:"created"`
	AccountID int32  `json:"account_id"`
}

func (m *Mercari) PostTransactionMessage(ctx context.Context, req *PostTransactionMessageRequest) (*PostTransactionMessageResponse, error) {
	postTransactionMessageFunc := func() (*PostTransactionMessageResponse, error) {
		token, err := m.GetToken(ctx, req.AccountID)
		if err != nil {
			hlog.CtxErrorf(ctx, "get token failed: %v", err)
			return nil, err
		}

		if ok := cache.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {token.AccessToken},
		}
		jsonReq := map[string]string{
			"message": req.Message,
		}
		reqBody, err := json.Marshal(jsonReq)
		if err != nil {
			hlog.CtxErrorf(ctx, "marshal json request error, %s", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		data := bytes.NewBuffer(reqBody)
		url := fmt.Sprintf("%s/v2/transactions/%s/messages", m.OpenApiDomain, req.TransactionId)
		httpReq, err := http.NewRequest("POST", url, data)
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
		if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusCreated {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxInfof(ctx, "post mercari transaction message error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}
		resp := &PostTransactionMessageResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxInfof(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		return resp, nil
	}

	result, err := backoff.Retry(ctx, postTransactionMessageFunc, m.GetRetryOpts()...)
	if err != nil {
		pErr := &backoff.PermanentError{}
		if errors.As(err, &pErr) {
			hlog.CtxInfof(ctx, "post mercari transaction message error: %v", err)
			berr := pErr.Unwrap()
			return nil, berr
		}
		hlog.CtxInfof(ctx, "post mercari transaction message error: %v", err)
		return nil, err
	}

	result.AccountID = req.AccountID
	return result, nil
}
