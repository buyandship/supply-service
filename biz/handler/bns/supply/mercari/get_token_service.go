package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetTokenService(ctx context.Context) (*supply.MercariGetTokenResp, error) {
	h := mercari.GetHandler()

	token, err := h.GetActiveToken(ctx)
	if err != nil {
		return nil, err
	}

	return &supply.MercariGetTokenResp{
		Token: token.AccessToken,
	}, nil
}
