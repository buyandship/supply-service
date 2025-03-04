package mercari

import (
	"context"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetSellerService(ctx context.Context, req *supply.MercariGetSellerReq) (*mercari.GetUserByUserIDResponse, error) {
	hlog.CtxInfof(ctx, "GetSellerService is called, %s", req.GetSellerID())

	if err := mock.MockMercariSellerError(req.GetSellerID()); err != nil {
		return nil, err
	}

	if req.GetSellerID() == "" {
		hlog.CtxInfof(ctx, "empty seller_id")
		return nil, bizErr.InvalidParameterError
	}

	h := mercari.GetHandler()

	resp, err := h.GetUser(ctx, &mercari.GetUserByUserIDRequest{
		UserId: req.GetSellerID(),
	})

	return resp, err
}
