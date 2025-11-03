package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetAuctionItemAuthService(ctx context.Context, req *supply.YahooGetAuctionItemAuthReq) (*supply.YahooGetAuctionItemAuthResp, error) {
	client := yahoo.GetClient()
	// TODO: get account_id
	auctionItemAuthResp, err := client.MockGetAuctionItemAuth(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID}, "")
	if err != nil {
		return nil, err
	}

	return &supply.YahooGetAuctionItemAuthResp{
		AuctionItem: &supply.YahooGetAuctionItemResp{
			AuctionID:    auctionItemAuthResp.ResultSet.Result.AuctionID,
			Title:        auctionItemAuthResp.ResultSet.Result.Title,
			Description:  auctionItemAuthResp.ResultSet.Result.Description,
			CurrentPrice: int64(auctionItemAuthResp.ResultSet.Result.CurrentPrice),
			StartPrice:   int64(auctionItemAuthResp.ResultSet.Result.StartPrice),
			Bids:         int32(auctionItemAuthResp.ResultSet.Result.Bids),
			ItemStatus:   auctionItemAuthResp.ResultSet.Result.ItemStatus,
			EndTime:      auctionItemAuthResp.ResultSet.Result.EndTime,
			StartTime:    auctionItemAuthResp.ResultSet.Result.StartTime,
			Seller: &supply.Seller{
				ID:     auctionItemAuthResp.ResultSet.Result.Seller.ID,
				Rating: auctionItemAuthResp.ResultSet.Result.Seller.Rating,
			},
			Image: auctionItemAuthResp.ResultSet.Result.Image,
		},
		IsWatching: auctionItemAuthResp.ResultSet.Result.IsWatching,
		BidStatus: &supply.BidStatus{
			HasBid:       auctionItemAuthResp.ResultSet.Result.MyBidStatus.HasBid,
			MyHighestBid: int64(auctionItemAuthResp.ResultSet.Result.MyBidStatus.MyHighestBid),
			IsWinning:    auctionItemAuthResp.ResultSet.Result.MyBidStatus.IsWinning,
		},
	}, nil
}
