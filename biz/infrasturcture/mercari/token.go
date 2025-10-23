package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type SetTokenRequest struct {
	RedirectUrl string `json:"redirect_url"`
	Code        string `json:"code"`
	Scope       string `json:"scope"`
	State       string `json:"state"`
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (m *Mercari) SetToken(ctx context.Context, req *supply.MercariLoginCallBackReq) error {
	if ok := cache.GetRedisClient().Limit(ctx); ok {
		return bizErr.RateLimitError
	}

	secret, err := m.GenerateSecret()
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {fmt.Sprintf("Basic %s", secret)},
	}

	body := fmt.Sprintf("grant_type=%s&scope=%s&redirect_uri=%s&code=%s", "authorization_code",
		url.QueryEscape(req.Scope), m.CallbackUrl, req.Code)

	url := fmt.Sprintf("%s/jp/v1/token", m.AuthServiceDomain)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(body)))
	if err != nil {
		hlog.CtxErrorf(ctx, "http request error, err: %v", err)
		return bizErr.InternalError
	}
	httpReq.Header = headers

	httpRes, err := m.Client.Do(ctx, httpReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "http error, err: %v", err)
		return bizErr.InternalError
	}
	defer func() {
		if err := httpRes.Body.Close(); err != nil {
			hlog.CtxErrorf(ctx, "http close error: %s", err)
		}
	}()

	if httpRes.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpRes.Body)
		hlog.CtxErrorf(ctx, "get mercari token error: %s", respBody)
		return bizErr.BizError{
			Status:  httpRes.StatusCode,
			ErrCode: httpRes.StatusCode,
			ErrMsg:  string(respBody),
		}
	}

	resp := &GetTokenResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
		hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
		return bizErr.InternalError
	}

	hlog.CtxInfof(ctx, "get token success, resp: %+v", resp)

	// insert token
	ExtractAccountIDFunc := func(s string) int32 {
		parts := strings.SplitN(s, "-", 3)
		if len(parts) >= 3 {
			accountId, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return 0
			}
			return int32(accountId)
		}
		return 0
	}

	accountId := ExtractAccountIDFunc(req.State)
	if accountId == 0 {
		hlog.CtxErrorf(ctx, "invalid account id: %s", req.State)
		return bizErr.InvalidParameterError
	}

	if err := db.GetHandler().InsertTokenLog(ctx, &mercari.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		Scope:        resp.Scope,
		TokenType:    resp.TokenType,
		AccountID:    int32(accountId),
	}); err != nil {
		return bizErr.InternalError
	}

	if err := cache.GetRedisClient().Del(ctx, fmt.Sprintf(config.TokenRedisKeyPrefix, accountId)); err != nil {
		return err
	}

	return nil
}
