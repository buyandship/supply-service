package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var FailureDetailsCodeMap = map[string]int{
	"F0017": 17,
	"F1004": 1004,
}

type PurchaseItemRequest struct {
	AccountId          int32  `json:"account_id"`
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
	CountryCode        string `json:"country_code"`
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

type GenericErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details []struct {
		Type      string `json:"@type"`
		RequestId string `json:"request_id"`
	} `json:"details"`
}

func (m *Mercari) PurchaseItem(ctx context.Context, refId string, req *PurchaseItemRequest) error {
	purchaseItemFunc := func() (*PurchaseItemResponse, error) {
		token, err := m.GetActiveToken(ctx)
		if err != nil {
			return nil, err
		}

		if ok := cache.GetHandler().Limit(ctx); ok {
			hlog.CtxWarnf(ctx, "hit rate limit")
			return nil, backoff.RetryAfter(1)
		}

		headers := map[string][]string{
			"Content-Type":  {"application/json"},
			"Accept":        {"application/json"},
			"Authorization": {token.AccessToken},
		}
		reqBody, err := json.Marshal(req)
		if err != nil {
			hlog.CtxErrorf(ctx, "marshal json request error, %s", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		data := bytes.NewBuffer(reqBody)
		url := fmt.Sprintf("%s/v1/items/purchase", m.OpenApiDomain)
		httpReq, err := http.NewRequest("POST", url, data)
		if err != nil {
			hlog.CtxErrorf(ctx, "http request error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		httpReq.Header = headers

		httpRes := &http.Response{}
		if config.GlobalServerConfig.Env == "development" && token.AccountID == 4 {
			// mock acl ban error
			// delete after testing
			httpRes.StatusCode = http.StatusForbidden
			httpRes.Body = io.NopCloser(bytes.NewBufferString(`{"request_id": "08f40f07-f67c-42c4-a47c-b794a7d7aa76", "transaction_status": "failure", "failure_details": {"code": "F0017", "reasons": "The buyer is currently forbidden to purchase items due to ACL Ban"}}`))
		} else {
			httpRes, err = HttpDo(ctx, httpReq)
			if err != nil {
				hlog.CtxErrorf(ctx, "http error, err: %v", err)
				return nil, backoff.Permanent(bizErr.InternalError)
			}
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

		// purchasing fails
		if httpRes.StatusCode != http.StatusOK {
			errResp := &PurchaseItemErrorResponse{}
			if err := json.NewDecoder(httpRes.Body).Decode(errResp); err != nil {
				hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
				return nil, backoff.Permanent(bizErr.InternalError)
			}

			hlog.CtxInfof(ctx, "purchase item error, error_code: [%d], error_msg: [%s]", httpRes.StatusCode, errResp.FailureDetails.Reasons)

			// if we purchase item fails, query by item
			getTxResp, err := m.GetTransactionByItemID(ctx, req.ItemId, req.AccountId)
			if err != nil {
				if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
					RefID:         refId,
					FailureReason: fmt.Sprintf("%s|%s", errResp.RequestId, errResp.FailureDetails.Reasons),
					AccountID:     req.AccountId,
				}); err != nil {
					hlog.CtxErrorf(ctx, "UpdateTransaction fail, [%s]", err.Error())
					return nil, backoff.Permanent(bizErr.InternalError)
				}
				var errMsg string
				if errResp.FailureDetails.Reasons == "" {
					// if the failure_details.reasons is empty, try to get the error from the http response
					errResp := &GenericErrorResponse{}
					if err := json.NewDecoder(httpRes.Body).Decode(errResp); err != nil {
						hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
					} else {
						if len(errResp.Details) > 0 {
							errMsg = fmt.Sprintf("%s|Generic Error", errResp.Details[0].RequestId)
						} else {
							errMsg = fmt.Sprintf("%d", httpRes.StatusCode)
						}
					}
				} else {
					errMsg = fmt.Sprintf("%s|%s", errResp.RequestId, errResp.FailureDetails.Reasons)
				}

				errCode := httpRes.StatusCode
				if e, ok := FailureDetailsCodeMap[errResp.FailureDetails.Code]; ok {
					errCode = e
				}

				if errCode == 17 {
					// failover
					if err := m.Failover(ctx, token.AccountID); err != nil {
						hlog.CtxWarnf(ctx, "The account:[%d] is banned,try to failover but fails, error: %s", token.AccountID, err.Error())
						time.Sleep(1 * time.Second)
						return nil, bizErr.ACLBanError
					}
					return nil, bizErr.ACLBanError
				}

				return nil, backoff.Permanent(bizErr.BizError{
					Status:  httpRes.StatusCode,
					ErrCode: errCode,
					ErrMsg:  errMsg,
				})
			}

			// if the query transaction by item_id return success, update the transaction
			if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
				RefID:            refId,
				PaidPrice:        getTxResp.PaidPrice,
				Price:            getTxResp.Price,
				TrxID:            getTxResp.Id,
				BuyerShippingFee: getTxResp.ShippingInfo.BuyerShippingFee,
				AccountID:        token.AccountID,
			}); err != nil {
				hlog.CtxErrorf(ctx, "UpdateTransaction fail, [%s]", err.Error())
				return nil, backoff.Permanent(bizErr.InternalError)
			}
			return nil, nil
		}

		// purchase success.
		resp := &PurchaseItemResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		// purchasing successfully
		trxId := strconv.FormatInt(resp.TransactionDetails.TrxId, 10)
		if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
			RefID:            refId,
			TrxID:            trxId,
			PaidPrice:        resp.TransactionDetails.PaidPrice,
			Price:            resp.TransactionDetails.Price,
			BuyerShippingFee: resp.TransactionDetails.BuyerShippingFee,
			AccountID:        token.AccountID,
		}); err != nil {
			hlog.CtxErrorf(ctx, "UpdateTransaction fail, [%s]", err.Error())
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		return nil, nil
	}

	_, err := backoff.Retry(ctx, purchaseItemFunc, m.GetRetryOpts()...)
	if err != nil {
		pErr := &backoff.PermanentError{}
		if errors.As(err, &pErr) {
			berr := pErr.Unwrap()
			return berr
		}
		return err
	}

	return nil
}
