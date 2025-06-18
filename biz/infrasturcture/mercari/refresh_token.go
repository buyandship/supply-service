package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int32  `json:"expires_in"` // in second, e.g. 3600
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (m *Mercari) GetToken(ctx context.Context) error {
	// load from redis cache
	if err := cache.GetHandler().Get(ctx, cache.TokenRedisKey, m.Token); err != nil {
		hlog.CtxInfof(ctx, "load from cache failed, err:%v", err)
		// Degrade to load from mysql
		t, err := db.GetHandler().GetToken(ctx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return bizErr.UnloginError
			}
			return err
		}
		m.Token = t
		go func() {
			if err := cache.GetHandler().Set(ctx, cache.TokenRedisKey, m.Token, 5*time.Minute); err != nil {
				hlog.Warnf("redis set failed, err:%v", err)
			}
		}()
	}
	if m.TokenExpired() {
		if err := m.RefreshToken(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mercari) RefreshToken(ctx context.Context) error {
	if err := m.refreshToken(ctx); err != nil {
		return err
	}
	if err := cache.GetHandler().Del(ctx, cache.TokenRedisKey); err != nil {
		return err
	}
	return nil
}

func (m *Mercari) refreshToken(ctx context.Context) error {
	// Try to acquire lock
	locked, err := cache.GetHandler().TryLock(ctx, "mercari_refresh_token")
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to acquire lock: %v", err)
		return bizErr.InternalError
	}
	if !locked {
		hlog.CtxErrorf(ctx, "failed to acquire lock: another refresh is in progress")
		return bizErr.ConflictError
	}
	defer func() {
		if err := cache.GetHandler().Unlock(ctx, "mercari_refresh_token"); err != nil {
			hlog.CtxErrorf(ctx, "failed to release lock: %v", err)
		}
	}()

	if ok := cache.GetHandler().Limit(ctx); ok {
		hlog.CtxErrorf(ctx, "rate limit error")
		return bizErr.RateLimitError
	}

	secret, err := m.GenerateSecret()
	if err != nil {
		return bizErr.InternalError
	}

	headers := map[string][]string{
		"Content-Type":  {"application/x-www-form-urlencoded"},
		"Authorization": {fmt.Sprintf("Basic %s", secret)},
	}

	body := fmt.Sprintf("grant_type=%s&scope=%s&refresh_token=%s", "refresh_token",
		url.QueryEscape(m.Token.Scope), m.Token.RefreshToken)

	data := bytes.NewBuffer([]byte(body))
	url := fmt.Sprintf("%s/jp/v1/token", m.AuthServiceDomain)
	httpReq, err := http.NewRequest("POST", url, data)
	if err != nil {
		hlog.CtxErrorf(ctx, "http request error, err: %v", err)
		return bizErr.InternalError
	}
	httpReq.Header = headers

	httpRes, err := HttpDo(ctx, httpReq)
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
		hlog.CtxErrorf(ctx, "refresh token error: %s", respBody)
		return bizErr.UnauthorisedError
	}

	resp := &RefreshTokenResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
		hlog.CtxErrorf(ctx, "decode http response error, err: %v", err)
		return bizErr.InternalError
	}

	hlog.CtxInfof(ctx, "refresh token success, resp: %+v", resp)

	if err := db.GetHandler().InsertTokenLog(ctx, &mercari.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		Scope:        resp.Scope,
		TokenType:    resp.TokenType,
	}); err != nil {
		return bizErr.InternalError
	}

	m.Token = &mercari.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		Scope:        resp.Scope,
		TokenType:    resp.TokenType,
	}

	return nil
}
