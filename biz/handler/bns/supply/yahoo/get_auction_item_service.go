package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetAuctionItemService(ctx context.Context, req *supply.YahooGetAuctionItemReq) (*supply.YahooGetAuctionItemResp, error) {
	client := yahoo.GetClient()
	auctionItemResp, err := client.MockGetAuctionItem(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID})
	if err != nil {
		return nil, err
	}

	return &supply.YahooGetAuctionItemResp{
		AuctionID:    auctionItemResp.ResultSet.Result.AuctionID,
		Title:        auctionItemResp.ResultSet.Result.Title,
		Description:  auctionItemResp.ResultSet.Result.Description,
		CurrentPrice: int64(auctionItemResp.ResultSet.Result.CurrentPrice),
		StartPrice:   int64(auctionItemResp.ResultSet.Result.StartPrice),
		Bids:         int32(auctionItemResp.ResultSet.Result.Bids),
		ItemStatus:   auctionItemResp.ResultSet.Result.ItemStatus,
		EndTime:      auctionItemResp.ResultSet.Result.EndTime,
		StartTime:    auctionItemResp.ResultSet.Result.StartTime,
		Seller: &supply.Seller{
			ID:     auctionItemResp.ResultSet.Result.Seller.ID,
			Rating: auctionItemResp.ResultSet.Result.Seller.Rating,
		},
		Image: auctionItemResp.ResultSet.Result.Image,
	}, nil
}
