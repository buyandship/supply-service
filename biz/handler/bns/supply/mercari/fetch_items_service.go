package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func FetchItemsService(ctx context.Context, req *supply.MercariFetchItemsReq) (*mercari.FetchItemsResponse, error) {

	h := mercari.GetHandler()

	return h.FetchItems(ctx, req)
}
