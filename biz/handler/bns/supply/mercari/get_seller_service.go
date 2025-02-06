package mercari

import (
	"context"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetSellerService(ctx context.Context, req *supply.MercariGetSellerReq) (*mercari.GetUserByUserIDResponse, error) {
	h := mercari.GetHandler()

	resp, err := h.GetUser(ctx, &mercari.GetUserByUserIDRequest{
		UserId:  req.GetSellerID(),
		BuyerId: req.GetBuyerID(),
	})

	return resp, err
}
