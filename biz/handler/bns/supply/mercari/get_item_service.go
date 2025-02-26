package mercari

import (
	"context"
	"errors"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func GetItemService(ctx context.Context, req *supply.MercariGetItemReq) (*mercari.GetItemByIDResponse, error) {
	hlog.CtxInfof(ctx, "GetItemService is called, item_id: %s", req.GetItemID())

	if req.GetItemID() == "" {
		hlog.CtxErrorf(ctx, "empty item_id")
		return nil, bizErr.InvalidParameterError
	}

	// Mock
	if err := mock.MockMercariGetItemError(req.GetItemID()); err != nil {
		return nil, err
	}

	h := mercari.GetHandler()

	var prefecture string

	var buyerId int32 = 1
	if req.GetBuyerID() != 0 {
		buyerId = req.GetBuyerID()
	}
	if req.BuyerID != 0 {
		acc, err := db.GetHandler().GetAccount(ctx, buyerId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, bizErr.InvalidBuyerError
		}
		if err != nil {
			return nil, bizErr.InternalError
		}
		prefecture = acc.Prefecture
	}

	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:     req.GetItemID(),
		BuyerId:    req.GetBuyerID(),
		Prefecture: prefecture,
	})

	if err := mock.MockMercariItemResponse(resp); err != nil {
		return nil, err
	}

	return resp, err
}
