package mercari

import (
	"context"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetItemService(ctx context.Context, req *supply.MercariGetItemReq) (*mercari.GetItemByIDResponse, error) {

	h := mercari.GetHandler()

	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:  req.GetItemID(),
		BuyerId: req.GetBuyerID(),
	})

	return resp, err
}
