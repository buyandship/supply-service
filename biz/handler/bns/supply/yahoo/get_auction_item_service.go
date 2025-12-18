package yahoo

import (
	"context"

	globalConfig "github.com/buyandship/bns-golib/config"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetAuctionItemService(ctx context.Context, req *supply.YahooGetAuctionItemReq) (*yahoo.AuctionItemDetail, error) {
	client := yahoo.GetClient()
	var yahooAccountID string
	switch globalConfig.GlobalAppConfig.Env {
	case "dev":
		yahooAccountID = config.DevYahoo02AccountID
	case "prod":
		yahooAccountID = config.ProdYahoo02AccountID
	}
	if err := mock.MockYahooGetAuctionItemError(req.AuctionID); err != nil {
		return nil, err
	}
	auctionItemResp, err := client.GetAuctionItemAuth(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID}, yahooAccountID)
	if err != nil {
		return nil, err
	}

	if err := mock.MockYahooGetAuctionItemDetail(auctionItemResp); err != nil {
		return nil, err
	}

	return &auctionItemResp.ResultSet.Result, nil
}
