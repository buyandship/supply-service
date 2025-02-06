package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
	"strconv"
)

type PostTransactionMessageRequest struct {
	TransactionId string `json:"transactionID"`
	Message       string `json:"message"`
	BuyerId       string `json:"buyer_id"`
}

type PostTransactionMessageResponse struct {
	Id      string `json:"id"`
	Body    string `json:"body"`
	UserId  string `json:"user_id"`
	Created int    `json:"created"`
}

func (m *Mercari) PostTransactionMessage(ctx context.Context, req *PostTransactionMessageRequest) (*PostTransactionMessageResponse, error) {
	postTransactionMessageFunc := func() (*PostTransactionMessageResponse, error) {
		acc, ok := m.Accounts[req.BuyerId]
		if !ok {
			hlog.Errorf("buyer not exists, buyer_id: %s", req.BuyerId)
			return nil, bizErr.InvalidBuyerError
		}

		headers := map[string][]string{
			"Content-Type":  []string{"application/json"},
			"Accept":        []string{"application/json"},
			"Authorization": []string{acc.AccessToken},
		}

		jsonReq := map[string]string{
			"message": req.Message,
		}
		reqBody, err := json.Marshal(jsonReq)
		if err != nil {
			hlog.Errorf("marshal json request error, %s", err.Error())
			return nil, bizErr.InternalError
		}
		data := bytes.NewBuffer(reqBody)
		httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/v2/transactions/%s/messages", m.OpenApiDomain, req.TransactionId), data)
		if err != nil {
			hlog.Errorf("http request error, err: %v", err)
			return nil, bizErr.InternalError
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
			hlog.Errorf("http error, err: %v", err)
			return nil, bizErr.InternalError
		}

		if httpRes.StatusCode == http.StatusUnauthorized {
			if err := m.RefreshToken(req.BuyerId); err != nil {
				hlog.Errorf("try to refresh token, but fails, err: %v", err)
			}
			seconds, err := strconv.ParseInt(httpRes.Header.Get("Retry-After"), 10, 64)
			if err == nil {
				return nil, backoff.RetryAfter(int(seconds))
			}
		}
		// retry code: 409, 429, 5xx
		if httpRes.StatusCode == http.StatusTooManyRequests {
			seconds, err := strconv.ParseInt(httpRes.Header.Get("Retry-After"), 10, 64)
			if err == nil {
				return nil, backoff.RetryAfter(int(seconds))
			}
		}
		if httpRes.StatusCode == http.StatusConflict {
			seconds, err := strconv.ParseInt(httpRes.Header.Get("Retry-After"), 10, 64)
			if err == nil {
				return nil, backoff.RetryAfter(int(seconds))
			}
		}
		if httpRes.StatusCode >= 500 && httpRes.StatusCode < 600 {
			seconds, err := strconv.ParseInt(httpRes.Header.Get("Retry-After"), 10, 64)
			if err == nil {
				return nil, backoff.RetryAfter(int(seconds))
			}
		}

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.Errorf("post mercari transaction message error: %s", respBody)
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  "post mercari transaction message fails",
			}
		}

		resp := &PostTransactionMessageResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.Errorf("decode http response error, err: %v", err)
			return nil, bizErr.InternalError
		}
		return resp, nil
	}

	result, err := backoff.Retry(context.TODO(), postTransactionMessageFunc, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}
	return result, nil
}
