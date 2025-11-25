package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetAuctionItemService(ctx context.Context, req *supply.YahooGetAuctionItemReq) (*yahoo.AuctionItemDetail, error) {
	client := yahoo.GetClient()
	auctionItemResp, err := client.GetAuctionItemAuth(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID}, config.MasterYahooAccountID)
	if err != nil {
		return nil, err
	}

	return &auctionItemResp.ResultSet.Result, nil
}
