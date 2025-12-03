package mercari

import (
	"encoding/base64"
	"fmt"

	"sync"

	"github.com/buyandship/bns-golib/config"
	"github.com/buyandship/bns-golib/http"
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
		switch config.GlobalAppConfig.Env {
		case "dev":
			authServiceDomain = "https://auth-sandbox.mercari.com"
			url = "https://api.mercari-sandbox.com"
		case "prod":
			authServiceDomain = "https://auth.mercari.com"
			url = "https://api.jp-mercari.com"
		}
		client := http.NewClient()

		Handler = &Mercari{
			AuthServiceDomain: authServiceDomain,
			OpenApiDomain:     url,
			ClientID:          config.GlobalAppConfig.GetString("mercari.client_id"),
			ClientSecret:      config.GlobalAppConfig.GetString("mercari.client_secret"),
			CallbackUrl:       config.GlobalAppConfig.GetString("mercari.callback_url"),
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
