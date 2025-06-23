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
	model "github.com/buyandship/supply-svr/biz/model/mercari"
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

func (m *Mercari) GetActiveToken(ctx context.Context) (*mercari.Token, error) {
	accountId := 0
	if err := cache.GetHandler().Get(ctx, cache.ActiveAccountId, &accountId); err != nil {
		accs, err := db.GetHandler().GetAccountList(ctx)
		if err != nil {
			return nil, err
		}
		if len(accs) == 0 {
			return nil, bizErr.InternalError
		}
		for _, acc := range accs {
			if acc.ActiveAt != nil {
				accountId = int(acc.ID)
				break
			}
		}
		if accountId == 0 {
			return nil, bizErr.InternalError
		}
		go func() {
			if err := cache.GetHandler().Set(context.Background(), cache.ActiveAccountId, accountId, time.Hour); err != nil {
				hlog.Warnf("redis set failed, err:%v", err)
			}
		}()
	}

	return m.GetToken(ctx, int32(accountId))
}

func (m *Mercari) GetToken(ctx context.Context, accountId int32) (*mercari.Token, error) {
	// load from redis cache
	token := &mercari.Token{}
	if err := cache.GetHandler().Get(ctx, fmt.Sprintf(cache.TokenRedisKeyPrefix, accountId), token); err != nil {
		hlog.CtxInfof(ctx, "the token is not in cache, load from mysql")
		// Degrade to load from mysql
		t, err := db.GetHandler().GetToken(ctx, accountId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizErr.UnloginError
			}
			return nil, err
		}
		go func() {
			if err := cache.GetHandler().Set(context.Background(), fmt.Sprintf(cache.TokenRedisKeyPrefix, accountId), t, 5*time.Minute); err != nil {
				hlog.Warnf("[goroutine] redis set failed, err:%v", err)
			}
		}()
		return t, nil
	}
	if token.Expired() {
		if err := m.RefreshToken(ctx, token); err != nil {
			return nil, err
		}
	}
	return token, nil
}

func (m *Mercari) RefreshToken(ctx context.Context, token *mercari.Token) error {
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
		url.QueryEscape(token.Scope), token.RefreshToken)

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

	token = &mercari.Token{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
		Scope:        resp.Scope,
		TokenType:    resp.TokenType,
		AccountID:    token.AccountID,
	}

	if err := db.GetHandler().InsertTokenLog(ctx, token); err != nil {
		hlog.CtxErrorf(ctx, "insert token log failed, err: %v", err)
		return bizErr.InternalError
	}

	if err := cache.GetHandler().Del(ctx, fmt.Sprintf(cache.TokenRedisKeyPrefix, token.AccountID)); err != nil {
		hlog.CtxErrorf(ctx, "delete token from cache failed, err: %v", err)
		return err
	}
	return nil
}

func (m *Mercari) Failover(ctx context.Context, accountId int32) error {
	hlog.CtxInfof(ctx, "failover account: %d", accountId)
	// set banned_at
	now := time.Now()
	if err := db.GetHandler().BanAccount(ctx, accountId); err != nil {
		return err
	}

	// get all accounts
	accs, err := db.GetHandler().GetAccountList(ctx)
	if err != nil {
		return err
	}

	activeAccountId := 0
	for _, acc := range accs {
		if acc.BannedAt == nil && acc.Priority > 0 {
			// set active_at
			if err := db.GetHandler().UpdateAccount(ctx, &model.Account{
				Model: gorm.Model{
					ID: uint(acc.ID),
				},
				ActiveAt: &now,
			}); err != nil {
				return err
			}
			hlog.CtxInfof(ctx, "set active account: %d", acc.ID)
			activeAccountId = int(acc.ID)
			break
		}
	}

	if activeAccountId == 0 {
		// alert
		hlog.CtxErrorf(ctx, "no active account found")
		return bizErr.InternalError
	}

	if err := cache.GetHandler().Set(ctx, cache.ActiveAccountId, activeAccountId, time.Hour); err != nil {
		return err
	}

	// get token
	token, err := m.GetToken(ctx, int32(activeAccountId))
	if err != nil {
		return err
	}

	// refresh token
	if err := m.RefreshToken(ctx, token); err != nil {
		return err
	}

	return nil
}
