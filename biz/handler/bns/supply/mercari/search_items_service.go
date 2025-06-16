package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func SearchItemsService(ctx context.Context, req *supply.MercariSearchItemsReq) (*mercari.SearchItemsResponse, error) {
	h := mercari.GetHandler()
	hlog.CtxInfof(ctx, "search items req: %+v", req)
	return h.SearchItems(ctx, req)
}
