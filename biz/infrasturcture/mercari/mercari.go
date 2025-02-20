package mercari

import (
	"encoding/base64"
	"fmt"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cenkalti/backoff/v5"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"sync"
)

var (
	once    sync.Once
	Handler *Mercari
)

type Mercari struct {
	AuthServiceDomain string
	OpenApiDomain     string
	ClientID          string
	ClientSecret      string
	CallbackUrl       string
	Token             *model.Token
}

func GetHandler() *Mercari {
	once.Do(func() {
		var url, authServiceDomain string
		if config.GlobalServerConfig.Env == "development" {
			authServiceDomain = "https://auth-sandbox.mercari.com"
			url = "https://api.mercari-sandbox.com"
		} else if config.GlobalServerConfig.Env == "production" {
			authServiceDomain = "https://auth.mercari.com"
			url = "https://api.jp-mercari.com"
		}

		t, err := db.GetHandler().GetToken()
		if err != nil {
			hlog.Fatal(err)
		}

		Handler = &Mercari{
			AuthServiceDomain: authServiceDomain,
			OpenApiDomain:     url,
			ClientID:          config.GlobalServerConfig.Mercari.ClientId,
			ClientSecret:      config.GlobalServerConfig.Mercari.ClientSecret,
			CallbackUrl:       config.GlobalServerConfig.Mercari.CallbackUrl,
			Token:             t,
		}
	})
	return Handler
}

func (m *Mercari) GenerateSecret() (string, error) {
	basicSecret := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", m.ClientID, m.ClientSecret)))
	return basicSecret, nil
}

func (m *Mercari) GetRetryOpts() []backoff.RetryOption {
	var opts []backoff.RetryOption
	opts = append(opts, backoff.WithBackOff(backoff.NewExponentialBackOff()))
	opts = append(opts, backoff.WithMaxTries(5))
	return opts
}

func (m *Mercari) TokenExpired() bool {
	if m.Token == nil {
		return true
	}
	// TODO: comment
	expiredTime := m.Token.CreatedAt.Add(time.Duration(m.Token.ExpiresIn-60) * time.Second)
	if time.Now().After(expiredTime) {
		return true
	}
	return false
}
