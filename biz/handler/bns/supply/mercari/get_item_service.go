package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/handler/bns/supply/mercari/utils"
	"github.com/buyandship/supply-service/biz/infrastructure/mercari"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
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
