package mercari

import (
	"bytes"
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
	"time"
)

type PostTransactionReviewRequest struct {
	TrxId   string `json:"trx_id"`
	Fame    string `json:"fame"`
	Message string `json:"message"`
}

type PostTransactionReviewResponse struct {
	FailureDetails struct {
		Code    string `json:"code"`
		Reasons string `json:"reasons"`
	} `json:"failure_details,omitempty"`
	RequestId    string `json:"request_id"`
	ReviewStatus string `json:"review_status"`
}

func (m *Mercari) PostTransactionReview(ctx context.Context, req *PostTransactionReviewRequest) (*PostTransactionReviewResponse, error) {
	postTransactionReviewFunc := func() (*PostTransactionReviewResponse, error) {
		hlog.CtxInfof(ctx, "call /v1/transactions/{transactionID}/post_review at %+v", time.Now())
		if ok := redis.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {m.Token.AccessToken},
		}
		jsonReq := map[string]string{
			"fame":    req.Fame,
			"message": req.Message,
		}
		reqBody, err := json.Marshal(jsonReq)
		if err != nil {
			hlog.CtxErrorf(ctx, "marshal json request error, %s", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		data := bytes.NewBuffer(reqBody)
		httpReq, err := http.NewRequest("POST",
			fmt.Sprintf("%s/v1/transactions/%s/post_review", m.OpenApiDomain, req.TrxId), data)
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
			if err := m.GetToken(ctx); err != nil {
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
			hlog.CtxErrorf(ctx, "http error, error_code: [%d], error_msg: [%s], retrying at [%+v]...",
				httpRes.StatusCode, respBody, time.Now().Local())
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			}
		}
		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxErrorf(ctx, "post mercari transaction review error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}
		resp := &PostTransactionReviewResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		hlog.CtxInfof(ctx, "post mercari transaction review response: %+v", resp)
		return resp, nil
	}

	result, err := backoff.Retry(ctx, postTransactionReviewFunc, m.GetRetryOpts()...)
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
