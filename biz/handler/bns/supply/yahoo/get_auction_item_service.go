package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/handler/bns/supply/utils"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetAuctionItemService(ctx context.Context, req *supply.YahooGetAuctionItemReq) (*yahoo.AuctionItemDetail, error) {

	if err := mock.MockYahooGetAuctionItemError(req.AuctionID); err != nil {
		return nil, err
	}

	auctionItemResp, err := utils.GetAuctionItem(ctx, req.AuctionID)
	if err != nil {
		return nil, err
	}

	if err := mock.MockYahooGetAuctionItemDetail(auctionItemResp); err != nil {
		return nil, err
	}

	return &auctionItemResp.ResultSet.Result, nil
}
