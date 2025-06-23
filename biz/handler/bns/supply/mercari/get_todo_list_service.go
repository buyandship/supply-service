package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetTodoListService(ctx context.Context, req *supply.MercariGetTodoListReq) (*mercari.GetTodoListResp, error) {
	if req.GetLimit() > 60 {
		req.Limit = 60
	}

	h := mercari.GetHandler()
	return h.GetTodoList(ctx, &mercari.GetTodoListReq{
		Limit:     int(req.GetLimit()),
		PageToken: req.GetPageToken(),
	})
}
