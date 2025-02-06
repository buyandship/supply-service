package mercari

import (
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cenkalti/backoff/v5"
	"io"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type GetUserByUserIDRequest struct {
	UserId  string `json:"userID"`
	BuyerId string `json:"buyerID"`
}

type GetUserByUserIDResponse struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	PhotoUrl          string `json:"photo_url"`
	PhotoThumbnailUrl string `json:"photo_thumbnail_url"`
	Ratings           struct {
		Good int `json:"good"`
		Bad  int `json:"bad"`
	} `json:"ratings"`
	NumRatings          int    `json:"num_ratings"`
	StarRatingScore     int    `json:"star_rating_score"`
	Created             int    `json:"created"`
	Proper              bool   `json:"proper"`
	Introduction        string `json:"introduction"`
	NumSellItems        int    `json:"num_sell_items"`
	HasIdentityVerified bool   `json:"has_identity_verified"`
	UserBadges          []struct {
		ID          int    `json:"ID"`
		Name        string `json:"Name"`
		Description string `json:"Description"`
		IconURL     string `json:"IconURL"`
	} `json:"user_badges"`
}

func (m *Mercari) GetUser(ctx context.Context, req *GetUserByUserIDRequest) (*GetUserByUserIDResponse, error) {
	getUserFunc := func() (*GetUserByUserIDResponse, error) {
		acc, ok := m.Accounts[req.BuyerId]
		if !ok {
			hlog.Errorf("buyer not exists, buyer_id: %s", req.BuyerId)
			return nil, bizErr.InvalidBuyerError
		}

		headers := map[string][]string{
			"Accept":        []string{"application/json"},
			"Authorization": []string{acc.AccessToken},
		}

		httpReq, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/users/%s", m.OpenApiDomain, req.UserId), nil)
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
			hlog.Errorf("get mercari user error: %s", respBody)
			return nil, bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  "get mercari user fails",
			}
		}

		resp := &GetUserByUserIDResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.Errorf("decode http response error, err: %v", err)
			return nil, bizErr.InternalError
		}
		return resp, nil
	}
	result, err := backoff.Retry(context.TODO(), getUserFunc, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	if err != nil {
		return nil, err
	}
	return result, nil
}
