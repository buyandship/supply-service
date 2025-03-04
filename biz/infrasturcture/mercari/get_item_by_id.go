package mercari

import (
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
	"net/url"
	"time"
)

type GetItemByIDRequest struct {
	ItemId     string `json:"itemID"`
	Prefecture string `json:"prefecture"`
}

type GetItemByIDResponse struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	ItemType    string `json:"item_type"`
	Description string `json:"description"`
	Updated     int    `json:"updated"`
	Created     int    `json:"created"`
	Seller      struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		NumSellItems int    `json:"num_sell_items"`
		Ratings      struct {
			Good   int `json:"good"`
			Normal int `json:"normal"`
			Bad    int `json:"bad"`
		} `json:"ratings"`
	} `json:"seller"`
	Photos       []string `json:"photos"`
	ItemCategory struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"item_category"`
	ItemCondition struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"item_condition"`
	ItemSize struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"item_size"`
	ItemBrand struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		SubName string `json:"sub_name"`
	} `json:"item_brand"`
	ShippingPayer struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Code string `json:"code"`
	} `json:"shipping_payer"`
	ShippingMethod struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"shipping_method"`
	ShippingFromArea struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"shipping_from_area"`
	ShippingDuration struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		MinDays int    `json:"min_days"`
		MaxDays int    `json:"max_days"`
	} `json:"shipping_duration"`
	ShippingClass struct {
		Id  int `json:"id"`
		Fee int `json:"fee"`
	} `json:"shipping_class"`
	CatalogDetails struct {
		ItemModel           string `json:"item_model"`
		ItemModelAttributes string `json:"item_model_attributes"`
		Color               string `json:"color"`
		Capacity            string `json:"capacity"`
		Carrier             string `json:"carrier"`
		Accessories         string `json:"accessories"`
		Imei                string `json:"imei"`
		Limitation          string `json:"limitation"`
	} `json:"catalog_details"`
	ItemDiscount struct {
		Expire         int `json:"expire"`
		CouponId       int `json:"coupon_id"`
		ReturnPercent  int `json:"return_percent"`
		ReturnAbsolute int `json:"return_absolute"`
	} `json:"item_discount"`
	Discounts struct {
		Breakdown struct {
			OfferToEveryone struct {
				Expire         int `json:"expire"`
				DiscountOrder  int `json:"discount_order"`
				ReturnPercent  int `json:"return_percent"`
				ReturnAbsolute int `json:"return_absolute"`
			} `json:"offer_to_everyone"`
			ItemCoupon struct {
				Expire         int `json:"expire"`
				CouponId       int `json:"coupon_id"`
				DiscountOrder  int `json:"discount_order"`
				ReturnPercent  int `json:"return_percent"`
				ReturnAbsolute int `json:"return_absolute"`
			} `json:"item_coupon"`
			InhouseItemDiscount struct {
				Expire         int `json:"expire"`
				DiscountOrder  int `json:"discount_order"`
				ReturnPercent  int `json:"return_percent"`
				ReturnAbsolute int `json:"return_absolute"`
			} `json:"inhouse_item_discount"`
		} `json:"breakdown"`
		TotalReturnPercent  int `json:"total_return_percent"`
		TotalReturnAbsolute int `json:"total_return_absolute"`
	} `json:"discounts"`
	NumComments              int    `json:"num_comments"`
	NumLikes                 int    `json:"num_likes"`
	Checksum                 string `json:"checksum"`
	AnshinItemAuthentication struct {
		IsAuthenticatable bool `json:"is_authenticatable"`
		Fee               int  `json:"fee"`
	} `json:"anshin_item_authentication"`
}

func (m *Mercari) GetItemByID(ctx context.Context, req *GetItemByIDRequest) (*GetItemByIDResponse, error) {
	getItemFunc := func() (*GetItemByIDResponse, error) {
		hlog.CtxInfof(ctx, "call /v1/items at %+v", time.Now().Local())
		if ok := redis.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}
		headers := map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {m.Token.AccessToken},
		}

		httpReq, err := http.NewRequest("GET",
			fmt.Sprintf("%s/v1/items/%s?prefecture=%s", m.OpenApiDomain, req.ItemId, url.QueryEscape(req.Prefecture)), nil)
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
			hlog.CtxErrorf(ctx, "get mercari item error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &GetItemByIDResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		hlog.CtxInfof(ctx, "get item by id successfully")
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
