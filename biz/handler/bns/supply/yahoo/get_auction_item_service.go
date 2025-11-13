package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetAuctionItemService(ctx context.Context, req *supply.YahooGetAuctionItemReq) (*yahoo.AuctionItemDetail, error) {
	client := yahoo.GetClient()
	auctionItemResp, err := client.MockGetAuctionItem(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID})
	if err != nil {
		return nil, err
	}

	return &auctionItemResp.ResultSet.Result, nil
}
