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
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type ToDo struct {
	ItemId   string `json:"item_id"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	PhotoUrl string `json:"photo_url"`
	Created  int    `json:"created"`
}

type GetTodoListReq struct {
	Limit     int    `json:"limit"`
	PageToken string `json:"page_token"`
}

type GetTodoListResp struct {
	Data          []ToDo `json:"data"`
	NextPageToken string `json:"next_page_token"`
}

func (m *Mercari) GetTodoList(ctx context.Context, req *GetTodoListReq) (*GetTodoListResp, error) {
	getItemFunc := func() (*GetTodoListResp, error) {
		hlog.CtxInfof(ctx, "call /v1/todolist at %+v", time.Now())

		token := &model.Token{}
		// TODO: get active token

		if ok := cache.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}
		headers := map[string][]string{
			"Authorization": {token.AccessToken},
		}

		reqBody, err := json.Marshal(req)
		if err != nil {
			hlog.CtxErrorf(ctx, "marshal json request error, %s", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		data := bytes.NewBuffer(reqBody)
		url := fmt.Sprintf("%s/v1/todolist", m.OpenApiDomain)
		httpReq, err := http.NewRequest("GET", url, data)
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
			if err := m.RefreshToken(ctx, token); err != nil {
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

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxErrorf(ctx, "get mercari todo list error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &GetTodoListResp{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		hlog.CtxInfof(ctx, "get item categories successfully")
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
