package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/handler/bns/supply/utils"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetItemService(ctx context.Context, req *supply.MercariGetItemReq) (*mercari.GetItemByIDResponse, error) {
	if req.GetItemID() == "" {
		hlog.CtxErrorf(ctx, "empty item_id")
		return nil, bizErr.InvalidParameterError
	}

	// Mock
	if err := mock.MockMercariGetItemError(req.GetItemID()); err != nil {
		return nil, err
	}

	h := mercari.GetHandler()

	token, err := h.GetActiveToken(ctx)
	if err != nil {
		hlog.CtxInfof(ctx, "GetActiveToken error: %v", err)
		return nil, err
	}

	acc, err := utils.GetAccount(ctx, token.AccountID)
	if err != nil {
		hlog.CtxErrorf(ctx, "GetBuyer error: %v", err)
		return nil, err
	}

	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:     req.GetItemID(),
		Prefecture: acc.Prefecture,
	})

	if err := mock.MockMercariItemResponse(resp); err != nil {
		return nil, err
	}

	return resp, err
}
