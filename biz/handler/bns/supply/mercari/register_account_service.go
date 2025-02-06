package mercari

import (
	"context"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
)

func RegisterAccountService(ctx context.Context, req *supply.MercariRegisterAccountReq) (*supply.MercariRegisterAccountResp, error) {
	// 1. get access token and refresh token.
	h := mercari.GetHandler()
	resp, err := h.GetToken(ctx, &mercari.GetTokenRequest{
		BuyerID:     req.BuyerID,
		RedirectUrl: req.RedirectUrl,
	})
	if err != nil {
		return nil, err
	}

	if err := db.GetHandler().UpsertAccount(ctx, &model.Account{
		Email:        req.GetEmail(),
		BuyerID:      req.GetBuyerID(),
		Prefecture:   req.GetPrefecture(),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ClientID:     req.GetClientID(),
		ClientSecret: req.GetClientSecret(),
	}); err != nil {
		return nil, err
	}

	return &supply.MercariRegisterAccountResp{}, nil
}
