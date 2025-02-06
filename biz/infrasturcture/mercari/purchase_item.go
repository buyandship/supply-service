package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cenkalti/backoff/v5"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type PurchaseItemRequest struct {
	BuyerId            string `json:"buyer_id"`
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
		TrxId            int    `json:"trx_id"`
		PaidMethod       string `json:"paid_method"`
		Price            int    `json:"price"`
		PaidPrice        int    `json:"paid_price"`
		BuyerShippingFee int    `json:"buyer_shipping_fee"`
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

		reqBody, err := json.Marshal(req)
		if err != nil {
			hlog.Errorf("marshal json request error, %s", err.Error())
			return nil, bizErr.InternalError
		}
		data := bytes.NewBuffer(reqBody)
		httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/items/purchase", m.OpenApiDomain), data)
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
			errResp := &PurchaseItemErrorResponse{}
			if err := json.NewDecoder(httpRes.Body).Decode(errResp); err != nil {
				hlog.Errorf("decode http response error, err: %v", err)
				return nil, bizErr.InternalError
			}
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  errResp.FailureDetails.Reasons,
			}
		}

		resp := &PurchaseItemResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.Errorf("decode http response error, err: %v", err)
			return nil, bizErr.InternalError
		}
		return resp, nil
	}
	result, err := backoff.Retry(context.TODO(), purchaseItemFunc, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}
	return result, nil
}
