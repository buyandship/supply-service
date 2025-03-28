package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTransactionByItemIdService(ctx context.Context, req *supply.MercariGetTransactionByItemIdReq) (*mercari.GetTransactionByItemIDResponse, error) {
	hlog.CtxInfof(ctx, "GetTransactionByItemIdService is called, item_id: %s", req.GetItemID())

	if req.GetItemID() == "" {
		hlog.CtxErrorf(ctx, "empty item_id")
		return nil, bizErr.InvalidParameterError
	}
	h := mercari.GetHandler()

	resp, err := h.GetTransactionByItemID(ctx, req.GetItemID())
	if err != nil {
		return nil, err
	}

	if err := mock.MockMercariGetTransactionByItemId(resp); err != nil {
		return nil, err
	}

	return resp, err
}
