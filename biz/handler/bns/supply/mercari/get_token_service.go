package mercari

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
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
