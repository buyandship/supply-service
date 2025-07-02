package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetAccountListService(ctx context.Context) (*supply.MercariGetAccountResp, error) {
	accounts, err := db.GetHandler().GetAccountList(ctx)
	if err != nil {
		return nil, err
	}

	resp := &supply.MercariGetAccountResp{
		Accounts: make([]*supply.Account, 0),
	}

	for _, account := range accounts {
		resp.Accounts = append(resp.Accounts, account.Thrift())
	}
	return resp, nil
}
