package mercari

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func SearchItemsService(ctx context.Context, req *supply.MercariSearchItemsReq) (*mercari.SearchItemsResponse, error) {
	h := mercari.GetHandler()
	return h.SearchItems(ctx, req)
}
