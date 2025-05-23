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
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

		if err := m.GetToken(ctx); err != nil {
			return nil, bizErr.InternalError
		}

		if ok := redis.GetHandler().Limit(ctx); ok {
			hlog.CtxErrorf(ctx, "hit rate limit")
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
		url := fmt.Sprintf("%s/v1/transactions/%s/post_review", m.OpenApiDomain, req.TrxId)
		httpReq, err := http.NewRequest("POST", url, data)
		if err != nil {
			hlog.CtxErrorf(ctx, "http request error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		httpReq.Header = headers

		httpRes, err := HttpDo(ctx, httpReq)
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
			hlog.CtxErrorf(ctx, "http unauthorized, refreshing token...")
			if err := m.RefreshToken(ctx); err != nil {
				hlog.CtxErrorf(ctx, "try to refresh token, but fails, err: %v", err)
				return nil, backoff.RetryAfter(1)
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

		resp := &PostTransactionReviewResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		if httpRes.StatusCode != http.StatusOK && httpRes.StatusCode != http.StatusCreated {
			errCode := httpRes.StatusCode
			if e, ok := FailureDetailsCodeMap[resp.FailureDetails.Code]; ok {
				errCode = e
			}
			hlog.CtxErrorf(ctx, "post mercari transaction review error: %s, trx_id: %s, request_id: %s, error_code: %s", resp.FailureDetails.Reasons, req.TrxId, resp.RequestId, resp.FailureDetails.Code)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: errCode,
				ErrMsg:  resp.FailureDetails.Reasons,
			})
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
