package mercari

import (
	"encoding/base64"
	"fmt"

	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
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
	Accounts          map[string]model.Account
}

func GetHandler() *Mercari {
	once.Do(func() {
		var url, authServiceDomain string
		if config.GlobalServerConfig.Env == "development" {
			authServiceDomain = "auth-sandbox.mercari.com"
			url = "https://api.mercari-sandbox.com"
		} else if config.GlobalServerConfig.Env == "production" {
			authServiceDomain = "auth.mercari.com"
			url = "https://api.mercari-sandbox.com"
		}

		accs, err := db.GetHandler().GetAccounts()
		if err != nil {
			hlog.Fatal(err)
		}

		Handler = &Mercari{
			AuthServiceDomain: authServiceDomain,
			OpenApiDomain:     url,
			Accounts:          accs,
		}
	})
	return Handler
}

func (m *Mercari) GenerateSecret(buyerID string) (string, error) {
	acc, ok := m.Accounts[buyerID]
	if !ok {
		hlog.Errorf("buyer not exists: %s", buyerID)
		return "", bizErr.InvalidBuyerError
	}
	basicSecret := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s:%s", acc.ClientID, acc.ClientSecret)))

	return basicSecret, nil
}
