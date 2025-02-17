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
	"time"
)

type GetUserByUserIDRequest struct {
	UserId string `json:"userID"`
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
		hlog.CtxInfof(ctx, "call /v1/users at %+v", time.Now())
		if ok := redis.GetHandler().Limit(ctx); ok {
			return nil, bizErr.RateLimitError
		}

		headers := map[string][]string{
			"Accept":        {"application/json"},
			"Authorization": {m.Token.AccessToken},
		}

		httpReq, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/users/%s", m.OpenApiDomain, req.UserId), nil)
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
			hlog.CtxErrorf(ctx, "get mercari user error: %s", respBody)
			return nil, backoff.Permanent(bizErr.BizError{
				Status:  httpRes.StatusCode,
				ErrCode: httpRes.StatusCode,
				ErrMsg:  string(respBody),
			})
		}

		resp := &GetUserByUserIDResponse{}
		if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
			hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
			return nil, backoff.Permanent(bizErr.InternalError)
		}
		hlog.CtxInfof(ctx, "get mercari user successfully")
		return resp, nil
	}
	result, err := backoff.Retry(ctx, getUserFunc, m.GetRetryOpts()...)
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
