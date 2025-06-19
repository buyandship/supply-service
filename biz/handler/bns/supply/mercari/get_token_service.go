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

	token, err := h.GetActiveToken(ctx)
	if err != nil {
		return nil, err
	}

	return &supply.MercariGetTokenResp{
		Token: token.AccessToken,
	}, nil
}
