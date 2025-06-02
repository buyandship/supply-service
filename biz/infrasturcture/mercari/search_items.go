package mercari

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type SearchItemsRequest struct {
	Keyword                 *string `json:"keyword,omitempty"`
	ExcludeKeyword          *string `json:"exclude_keyword,omitempty"`
	CategoryId              *int    `json:"category_id,omitempty"`
	SellerId                *int    `json:"seller_id,omitempty"`
	ShopId                  *string `json:"shop_id,omitempty"`
	SizeId                  *int    `json:"size_id,omitempty"`
	ColorId                 *int    `json:"color_id,omitempty"`
	PriceMin                *int    `json:"price_min,omitempty"`
	PriceMax                *int    `json:"price_max,omitempty"`
	CreatedBeforeDate       *int    `json:"created_before_date,omitempty"`
	CreatedAfterDate        *int    `json:"created_after_date,omitempty"`
	ItemConditionId         *int    `json:"item_condition_id,omitempty"`
	ShippingPayerId         *int    `json:"shipping_payer_id,omitempty"`
	Status                  *string `json:"status,omitempty"`
	Marketplace             *int    `json:"marketplace,omitempty" validate:"oneof=1 2 3"`
	Sort                    *string `json:"sort,omitempty" validate:"oneof=score created_time price num_likes"`
	Order                   *string `json:"order,omitempty" validate:"oneof=asc desc"`
	Page                    *int    `json:"page,omitempty"`
	Limit                   *int    `json:"limit,omitempty"`
	ItemAuthentication      *bool   `json:"item_authentication,omitempty"`
	TimeSale                *bool   `json:"time_sale,omitempty"`
	WithOfferPricePromotion *bool   `json:"with_offer_price_promotion,omitempty"`
}

type SearchItemsResponse struct {
	Data []Item `json:"data"`
	Meta struct {
		HasNext  bool `json:"has_next"`
		NumFound int  `json:"num_found"`
	} `json:"meta"`
}

type Item struct {
	ID                       string                    `json:"id"`
	Status                   string                    `json:"status"`
	Name                     string                    `json:"name"`
	Price                    int                       `json:"price"`
	ItemType                 string                    `json:"item_type"`
	Description              string                    `json:"description"`
	Updated                  int64                     `json:"updated"`
	Created                  int64                     `json:"created"`
	Seller                   Seller                    `json:"seller"`
	Thumbnail                string                    `json:"thumbnail"`
	Photos                   []string                  `json:"photos"`
	ItemCondition            ItemCondition             `json:"item_condition"`
	ShippingPayer            ShippingPayer             `json:"shipping_payer"`
	ShippingDuration         ShippingDuration          `json:"shipping_duration"`
	ItemCategory             ItemCategory              `json:"item_category"`
	ItemBrand                *ItemBrand                `json:"item_brand,omitempty"`
	ItemSizes                []ItemSize                `json:"item_sizes,omitempty"`
	ItemDiscount             *ItemDiscount             `json:"item_discount,omitempty"`
	AnshinItemAuthentication *AnshinItemAuthentication `json:"anshin_item_authentication,omitempty"`
}

type Seller struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Rating float64 `json:"rating"`
	ShopID string  `json:"shop_id,omitempty"`
}

type ItemCondition struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShippingPayer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type ShippingDuration struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	MinDays int    `json:"min_days"`
	MaxDays int    `json:"max_days"`
}

type ItemCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemBrand struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	SubName string `json:"sub_name"`
}

type ItemSize struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ItemDiscount struct {
	ReturnPercent int `json:"return_percent"`
}

type AnshinItemAuthentication struct {
	IsAuthenticatable bool `json:"is_authenticatable"`
	Fee               int  `json:"fee"`
}

func (m *Mercari) SearchItems(ctx context.Context, req *supply.MercariSearchItemsReq) (*SearchItemsResponse, error) {
	SearchItemsFunc := func() (*SearchItemsResponse, error) {
		hlog.CtxInfof(ctx, "call /v1/items at %+v", time.Now().Local())

		if err := m.GetToken(ctx); err != nil {
			return nil, bizErr.InternalError
		}

		if ok := redis.GetHandler().Limit(ctx); ok {
			hlog.CtxErrorf(ctx, "hit rate limit")
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {m.Token.AccessToken},
		}

		baseUrl, err := url.Parse(fmt.Sprintf("%s/v3/items/search", m.OpenApiDomain))
		if err != nil {
			hlog.CtxErrorf(ctx, "url parse error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		queryParams := url.Values{}
		if req.Keyword != nil {
			queryParams.Add("keyword", *req.Keyword)
		}
		if req.ExcludeKeyword != nil {
			queryParams.Add("exclude_keyword", *req.ExcludeKeyword)
		}
		if req.BrandID != nil {
			queryParams.Add("brand_id", strconv.Itoa(int(*req.BrandID)))
		}
		if req.CategoryID != nil {
			queryParams.Add("category_id", strconv.Itoa(int(*req.CategoryID)))
		}
		if req.SellerID != nil {
			queryParams.Add("seller_id", strconv.Itoa(int(*req.SellerID)))
		}
		if req.ShopID != nil {
			queryParams.Add("shop_id", *req.ShopID)
		}
		if req.SizeID != nil {
			queryParams.Add("size_id", strconv.Itoa(int(*req.SizeID)))
		}
		if req.ColorID != nil {
			queryParams.Add("color_id", strconv.Itoa(int(*req.ColorID)))
		}
		if req.PriceMin != nil {
			queryParams.Add("price_min", strconv.Itoa(int(*req.PriceMin)))
		}
		if req.PriceMax != nil {
			queryParams.Add("price_max", strconv.Itoa(int(*req.PriceMax)))
		}
		if req.CreatedBeforeDate != nil {
			queryParams.Add("created_before_date", strconv.Itoa(int(*req.CreatedBeforeDate)))
		}
		if req.CreatedAfterDate != nil {
			queryParams.Add("created_after_date", strconv.Itoa(int(*req.CreatedAfterDate)))
		}
		if req.ItemConditionID != nil {
			queryParams.Add("item_condition_id", strconv.Itoa(int(*req.ItemConditionID)))
		}
		if req.ShippingPayerID != nil {
			queryParams.Add("shipping_payer_id", strconv.Itoa(int(*req.ShippingPayerID)))
		}
		if req.Status != nil {
			queryParams.Add("status", *req.Status)
		}
		if req.Marketplace != nil {
			queryParams.Add("marketplace", strconv.Itoa(int(*req.Marketplace)))
		}
		if req.Sort != nil {
			queryParams.Add("sort", *req.Sort)
		}
		if req.Order != nil {
			queryParams.Add("order", *req.Order)
		}
		if req.Page != nil {
			queryParams.Add("page", strconv.Itoa(int(*req.Page)))
		}
		if req.Limit != nil {
			queryParams.Add("limit", strconv.Itoa(int(*req.Limit)))
		}
		if req.ItemAuthentication != nil {
			queryParams.Add("item_authentication", strconv.FormatBool(*req.ItemAuthentication))
		}
		if req.TimeSale != nil {
			queryParams.Add("time_sale", strconv.FormatBool(*req.TimeSale))
		}
		if req.WithOfferPricePromotion != nil {
			queryParams.Add("with_offer_price_promotion", strconv.FormatBool(*req.WithOfferPricePromotion))
		}
		baseUrl.RawQuery = queryParams.Encode()

		httpReq, err := http.NewRequest("GET", baseUrl.String(), nil)
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

		if httpRes.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(httpRes.Body)
			hlog.CtxInfof(ctx, "get mercari item error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &SearchItemsResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}

		return resp, nil
	}
	result, err := backoff.Retry(ctx, SearchItemsFunc, m.GetRetryOpts()...)
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
