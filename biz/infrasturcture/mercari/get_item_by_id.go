package mercari

import (
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"io"
	"net/http"
	"strconv"

	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type GetItemByIDRequest struct {
	ItemId  string `json:"itemID"`
	BuyerId string `json:"buyerID"`
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
		acc, ok := m.Accounts[req.BuyerId]
		if !ok {
			hlog.Errorf("buyer not exists, buyer_id: %s", req.BuyerId)
			return nil, bizErr.InvalidBuyerError
		}

		headers := map[string][]string{
			"Accept":        []string{"application/json"},
			"Authorization": []string{acc.AccessToken},
		}

		httpReq, err := http.NewRequest("GET",
			fmt.Sprintf("%s/v1/items/%s?prefecture=%s", m.OpenApiDomain, req.ItemId, acc.Prefecture), nil)
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
			hlog.Errorf("get mercari item error: %s", respBody)
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  "get mercari item fails",
			}
		}

		resp := &GetItemByIDResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.Errorf("decode http response error, err: %v", err)
			return nil, bizErr.InternalError
		}

		return resp, nil
	}
	result, err := backoff.Retry(context.TODO(), getItemFunc, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}
	return result, nil
}
