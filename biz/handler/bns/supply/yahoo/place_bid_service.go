package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func PlaceBidService(ctx context.Context, req *supply.YahooPlaceBidReq) (resp *supply.YahooPlaceBidResp, err error) {
	client := yahoo.GetClient()

	bidReq := yahoo.PlaceBidRequest{}
	placeBidResp, err := client.PlaceBid(ctx, &bidReq)
	if err != nil {
		return nil, err
	}

	return placeBidResp, nil

}
