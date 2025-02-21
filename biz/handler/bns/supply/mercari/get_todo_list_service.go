package mercari

import (
	"context"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTodoListService(ctx context.Context, req *supply.MercariGetTodoListReq) (*mercari.GetTodoListResp, error) {
	hlog.CtxInfof(ctx, "GetTodoListService is called, req: %s", req)

	if req.GetLimit() > 60 {
		req.Limit = 60
	}

	h := mercari.GetHandler()
	resp, err := h.GetTodoList(ctx, &mercari.GetTodoListReq{
		Limit:     int(req.GetLimit()),
		PageToken: req.GetPageToken(),
	})

	return resp, err
}
