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

	"github.com/buyandship/bns-golib/cache"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/logger/zap"
)

type PostTransactionReviewRequest struct {
	TrxId     string `json:"trx_id"`
	Fame      string `json:"fame"`
	Message   string `json:"message"`
	AccountID int32  `json:"account_id"`
}

type PostTransactionReviewResponse struct {
	FailureDetails struct {
		Code    string `json:"code"`
		Reasons string `json:"reasons"`
	} `json:"failure_details,omitempty"`
	RequestId    string `json:"request_id"`
	ReviewStatus string `json:"review_status"`
	AccountID    int32  `json:"account_id"`
}

func (m *Mercari) PostTransactionReview(ctx context.Context, req *PostTransactionReviewRequest) (*PostTransactionReviewResponse, error) {
	postTransactionReviewFunc := func() (*PostTransactionReviewResponse, error) {
		token, err := m.GetToken(ctx, req.AccountID)
		if err != nil {
			hlog.CtxErrorf(ctx, "get token failed: %v", err)
			return nil, err
		}

		if ok := cache.GetRedisClient().Limit(ctx); ok {
			hlog.CtxWarnf(ctx, "hit rate limit")
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {token.AccessToken},
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

		httpRes, err := m.Client.Do(ctx, httpReq)
		if err != nil {
			hlog.CtxErrorf(ctx, "http error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		defer func() {
			if err := httpRes.Body.Close(); err != nil {
				hlog.CtxWarnf(ctx, "http close error: %s", err)
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

		if httpRes.StatusCode == http.StatusNotFound {
			resp := &GenericErrorResponse{}
			if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
				hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
				return nil, backoff.Permanent(bizErr.InternalError)
			}
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  fmt.Sprintf("[error_message: %s, details: %+v, request_id: %s]", resp.Message, resp.Details, ctx.Value(zap.ExtraKey("X-Request-ID"))),
			})
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
			hlog.CtxInfof(ctx, "post mercari transaction review error: %s, trx_id: %s, request_id: %s, error_code: %s", resp.FailureDetails.Reasons, req.TrxId, resp.RequestId, resp.FailureDetails.Code)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: errCode,
				ErrMsg:  fmt.Sprintf("[error_message: %s, code: %+v, request_id: %s]", resp.FailureDetails.Reasons, resp.FailureDetails.Code, ctx.Value(zap.ExtraKey("X-Request-ID"))),
			})
		}

		return resp, nil
	}

	result, err := backoff.Retry(ctx, postTransactionReviewFunc, m.GetRetryOpts()...)
	if err != nil {
		pErr := &backoff.PermanentError{}
		if errors.As(err, &pErr) {
			hlog.CtxInfof(ctx, "post mercari transaction review error: %v", err)
			berr := pErr.Unwrap()
			return nil, berr
		}
		hlog.CtxInfof(ctx, "post mercari transaction review error: %v", err)
		return nil, err
	}

	result.AccountID = req.AccountID
	return result, nil
}
