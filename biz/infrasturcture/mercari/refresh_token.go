package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
	"net/url"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"` // in second, e.g. 3600
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type T struct {
	Info struct {
		AccessToken string `json:"access_token"`
	} `json:"info"`
}

func (m *Mercari) GetToken(ctx context.Context) error {
	// TODO: use redis
	t, err := db.GetHandler().GetToken()
	if err != nil {
		return err
	}
	m.Token = t

	if m.TokenExpired() {
		if err := m.refreshToken(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mercari) refreshToken(ctx context.Context) error {
	if ok := redis.GetHandler().Limit(ctx); ok {
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

	body := fmt.Sprintf("grant_type=%s&scope=%s&refresh_token=%s", "refresh_token",
		url.QueryEscape(m.Token.Scope), m.Token.RefreshToken)

	data := bytes.NewBuffer([]byte(body))
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/jp/v1/token", m.AuthServiceDomain), data)
	if err != nil {
		hlog.CtxErrorf(ctx, "http request error, err: %v", err)
		return bizErr.InternalError
	}
	httpReq.Header = headers
	c := &http.Client{}
	httpRes, err := c.Do(httpReq)
	defer func() {
		if err := httpRes.Body.Close(); err != nil {
			hlog.CtxErrorf(ctx, "http close error: %s", err)
		}
	}()
	if err != nil {
		hlog.CtxErrorf(ctx, "http error, err: %v", err)
		return bizErr.InternalError
	}
	if httpRes.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpRes.Body)
		hlog.CtxErrorf(ctx, "refresh token error: %s", respBody)
		return bizErr.UnauthorisedError
	}

	resp := &RefreshTokenResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
		hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
		return bizErr.InternalError
	}

	if err := db.GetHandler().InsertTokenLog(context.Background(), &mercari.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		Scope:        resp.Scope,
		TokenType:    resp.TokenType,
	}); err != nil {
		return bizErr.InternalError
	}

	m.Token, err = db.GetHandler().GetToken()
	if err != nil {
		return bizErr.InternalError
	}

	return nil
}
