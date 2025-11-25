package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func LoginCallBackService(ctx context.Context, req *supply.MercariLoginCallBackReq) error {
	hlog.CtxInfof(ctx, "LoginCallBackService called, req: %+v", req)
	h := mercari.GetHandler()

	if req.GetCode() == "" {
		hlog.CtxErrorf(ctx, "empty code")
		return bizErr.InvalidParameterError
	}

	if req.GetScope() == "" {
		hlog.CtxErrorf(ctx, "empty scope")
		return bizErr.InvalidParameterError
	}

	if err := h.SetToken(ctx, req); err != nil {
		return err
	}

	return nil
}
