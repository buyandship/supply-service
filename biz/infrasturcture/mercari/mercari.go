package mercari

import (
	"encoding/base64"
	"fmt"

	"github.com/cenkalti/backoff/v5"

	"sync"

	"github.com/buyandship/bns-golib/http"
	"github.com/buyandship/supply-svr/biz/common/config"
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
	Client            *http.Client
}

func GetHandler() *Mercari {
	once.Do(func() {
		var url, authServiceDomain string
		switch config.GlobalServerConfig.Env {
		case "development":
			authServiceDomain = "https://auth-sandbox.mercari.com"
			url = "https://api.mercari-sandbox.com"
		case "production":
			authServiceDomain = "https://auth.mercari.com"
			url = "https://api.jp-mercari.com"
		}
		client := http.NewClient()

		Handler = &Mercari{
			AuthServiceDomain: authServiceDomain,
			OpenApiDomain:     url,
			ClientID:          config.GlobalServerConfig.Mercari.ClientId,
			ClientSecret:      config.GlobalServerConfig.Mercari.ClientSecret,
			CallbackUrl:       config.GlobalServerConfig.Mercari.CallbackUrl,
			Client:            client,
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
