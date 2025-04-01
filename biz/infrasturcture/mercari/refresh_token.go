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
	"github.com/buyandship/supply-svr/biz/common/trace"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
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
	if err := m.LoadTokenFromCache(ctx); err != nil {
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
		js, err := json.Marshal(m.Token)
		if err != nil {
			hlog.CtxInfof(ctx, "marshal json failed, err:%v", err)
		} else {
			if err := redis.GetHandler().Set(ctx, redis.TokenRedisKey, string(js), 5*time.Minute); err != nil {
				hlog.CtxErrorf(ctx, "redis set failed, err:%v", err)
				return err
			}
		}
	}
	if m.TokenExpired() {
		if err := m.RefreshToken(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mercari) LoadTokenFromCache(ctx context.Context) error {
	s, err := redis.GetHandler().Get(ctx, redis.TokenRedisKey)
	if err != nil {
		return err
	}
	if t, ok := s.(string); ok {
		hlog.CtxInfof(ctx, "get token from cache %s", t)
		if err := json.Unmarshal([]byte(t), m.Token); err != nil {
			return bizErr.InternalError
		}
		return nil
	}
	return fmt.Errorf("get token from cache failed")
}

func (m *Mercari) RefreshToken(ctx context.Context) error {
	if err := m.refreshToken(ctx); err != nil {
		return err
	}
	if err := redis.GetHandler().Del(ctx, redis.TokenRedisKey); err != nil {
		return err
	}
	return nil
}

func (m *Mercari) refreshToken(ctx context.Context) error {
	// Try to acquire lock
	locked, err := redis.GetHandler().TryLock(ctx, "mercari_refresh_token")
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to acquire lock: %v", err)
		return bizErr.InternalError
	}
	if !locked {
		hlog.CtxErrorf(ctx, "failed to acquire lock: another refresh is in progress")
		return bizErr.ConflictError
	}
	defer func() {
		if err := redis.GetHandler().Unlock(ctx, "mercari_refresh_token"); err != nil {
			hlog.CtxErrorf(ctx, "failed to release lock: %v", err)
		}
	}()

	if ok := redis.GetHandler().Limit(ctx); ok {
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
	c := &http.Client{}

	ctx, span := trace.StartHTTPOperation(ctx, "POST", url)
	defer trace.EndSpan(span, nil)

	hlog.CtxInfof(ctx, "refresh token request, refresh_token=%s, access_token=%s", m.Token.RefreshToken, m.Token.AccessToken)

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
	trace.RecordHTTPResponse(span, httpRes)
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
