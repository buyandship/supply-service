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
)

type PurchaseItemRequest struct {
	BuyerId            int32  `json:"buyer_id"`
	ItemId             string `json:"item_id"`
	FamilyName         string `json:"family_name"`
	FirstName          string `json:"first_name"`
	FamilyNameKana     string `json:"family_name_kana"`
	FirstNameKana      string `json:"first_name_kana"`
	Telephone          string `json:"telephone"`
	ZipCode1           string `json:"zip_code1"`
	ZipCode2           string `json:"zip_code2"`
	Prefecture         string `json:"prefecture"`
	City               string `json:"city"`
	Address1           string `json:"address1"`
	Address2           string `json:"address2"`
	DeliveryIdentifier string `json:"delivery_identifier"`
	Checksum           string `json:"checksum"`
	CouponId           int    `json:"coupon_id"`
	ItemAuthentication bool   `json:"item_authentication"`
}

type PurchaseItemResponse struct {
	RequestId          string `json:"request_id"`
	TransactionStatus  string `json:"transaction_status"`
	TransactionDetails struct {
		TrxId            int64  `json:"trx_id"`
		PaidMethod       string `json:"paid_method"`
		Price            int64  `json:"price"`
		PaidPrice        int64  `json:"paid_price"`
		BuyerShippingFee string `json:"buyer_shipping_fee"`
		ItemId           string `json:"item_id"`
		Checksum         string `json:"checksum"`
		UserAddress      struct {
			ZipCode1       string `json:"zip_code1"`
			ZipCode2       string `json:"zip_code2"`
			FamilyName     string `json:"family_name"`
			FirstName      string `json:"first_name"`
			FamilyNameKana string `json:"family_name_kana"`
			FirstNameKana  string `json:"first_name_kana"`
			Prefecture     string `json:"prefecture"`
			City           string `json:"city"`
			Address1       string `json:"address1"`
			Address2       string `json:"address2"`
			Telephone      string `json:"telephone"`
		} `json:"user_address"`
		ShippingMethodId int `json:"shipping_method_id"`
	} `json:"transaction_details"`

	CouponId int `json:"coupon_id"`
}

type PurchaseItemErrorResponse struct {
	RequestId         string `json:"request_id"`
	TransactionStatus string `json:"transaction_status"`
	FailureDetails    struct {
		Code    string `json:"code"`
		Reasons string `json:"reasons"`
	} `json:"failure_details"`
}

func (m *Mercari) PurchaseItem(ctx context.Context, req *PurchaseItemRequest) (*PurchaseItemResponse, error) {
	purchaseItemFunc := func() (*PurchaseItemResponse, error) {
		if ok := redis.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {m.Token.AccessToken},
		}
		reqBody, err := json.Marshal(req)
		if err != nil {
			hlog.CtxErrorf(ctx, "marshal json request error, %s", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		data := bytes.NewBuffer(reqBody)
		httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/items/purchase", m.OpenApiDomain), data)
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
			hlog.CtxErrorf(ctx, "http error, error_code: [%d], err_msg:[%s], retrying...", httpRes.StatusCode, respBody)
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			}
		}

		if httpRes.StatusCode != http.StatusOK {
			hlog.CtxErrorf(ctx, "http error, error_code: [%d]", httpRes.StatusCode)
			errResp := &PurchaseItemErrorResponse{}
			if err := json.NewDecoder(httpRes.Body).Decode(errResp); err != nil {
				hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
				return nil, backoff.Permanent(bizErr.InternalError)
			}
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  errResp.FailureDetails.Reasons,
			})
		}

		resp := &PurchaseItemResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		return resp, nil
	}
	result, err := backoff.Retry(ctx, purchaseItemFunc, m.GetRetryOpts()...)
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
