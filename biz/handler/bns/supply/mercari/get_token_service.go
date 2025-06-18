package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTokenService(ctx context.Context) (*supply.MercariGetTokenResp, error) {
	hlog.CtxInfof(ctx, "Getting token service")
	h := mercari.GetHandler()

	// TODO: get active account

	token, err := h.GetToken(ctx, 0)
	if err != nil {
		return nil, err
	}

	// TODO: get account from db
	return &supply.MercariGetTokenResp{
		Token:   token.AccessToken,
		Account: nil,
	}, nil
}
