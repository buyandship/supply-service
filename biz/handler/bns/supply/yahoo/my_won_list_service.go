package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetMyWonListService(ctx context.Context, req *supply.YahooGetMyWonListReq) (*yahoo.MyWonListResponse, error) {
	client := yahoo.GetClient()
	resp, err := client.GetMyWonList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
