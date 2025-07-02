package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTransactionByItemIdService(ctx context.Context, req *supply.MercariGetTransactionByItemIdReq) (*mercari.GetTransactionByItemIDResponse, error) {
	if req.GetItemID() == "" {
		hlog.CtxErrorf(ctx, "empty item_id")
		return nil, bizErr.InvalidParameterError
	}
	h := mercari.GetHandler()

	trx, err := db.GetHandler().GetTransaction(ctx, &model.Transaction{
		ItemID: req.GetItemID(),
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "get transaction failed: %v", err)
		return nil, bizErr.NotFoundError
	}

	resp, err := h.GetTransactionByItemID(ctx, req.GetItemID(), trx.AccountID)
	if err != nil {
		return nil, err
	}

	if err := mock.MockMercariGetTransactionByItemId(resp); err != nil {
		return nil, err
	}

	return resp, err
}
