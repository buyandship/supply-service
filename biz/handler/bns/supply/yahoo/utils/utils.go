package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/mq"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/mercari"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func GetAccount(ctx context.Context, accId int32) (*mercari.Account, error) {

	acc := &mercari.Account{}
	if err := cache.GetRedisClient().Get(ctx, fmt.Sprintf(config.MercariAccountPrefix, accId), acc); err != nil {
		// degrade to load from
		acc, err := db.GetHandler().GetAccount(ctx, accId)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := cache.GetRedisClient().Set(context.Background(), fmt.Sprintf(config.MercariAccountPrefix, accId), acc, time.Hour); err != nil {
				hlog.Warnf("[goroutine] redis set buyer error: %+v", err)
			}
		}()
		return acc, err
	}
	return acc, nil
}

func GetAuctionItem(ctx context.Context, auctionID string) (*yahoo.AuctionItemResponse, error) {
	client := yahoo.GetClient()
	auctionItemResp, err := client.GetAuctionItemAuth(ctx, yahoo.AuctionItemRequest{AuctionID: auctionID})
	if err != nil {
		return nil, err
	}

	// TODO: update bid request status
	go func() {
		// get bid request
		bidRequest, err := db.GetHandler().GetBidRequestByAuctionID(ctx, auctionID)
		if err != nil {
			hlog.CtxErrorf(ctx, "get bid request failed: %+v", err)
			return
		}
		// check win bid

		// TBC: use TaxinBidPrice or Price?
		if int64(auctionItemResp.ResultSet.Result.Price) > bidRequest.MaxBid {
			// out bid
			if err := OutBid(ctx, bidRequest.OrderID); err != nil {
				hlog.CtxErrorf(ctx, "out bid failed: %+v", err)
				return
			}
		} else {
			if _, err := AddBid(ctx, &yahoo.PlaceBidRequest{
				YsRefID:         bidRequest.OrderID,
				TransactionType: bidRequest.RequestType,
				AuctionID:       auctionID,
				Price:           int(auctionItemResp.ResultSet.Result.Price),
			}, &auctionItemResp.ResultSet.Result); err != nil {
				hlog.CtxErrorf(ctx, "add bid failed: %+v", err)
				return
			}
		}
	}()

	return auctionItemResp, nil
}

func Buyout(ctx context.Context, req *yahoo.PlaceBidRequest, item *yahoo.AuctionItemDetail) (resp *yahoo.PlaceBidResult, err error) {
	// directly buyout
	placeBidResp, err := yahoo.GetClient().PlaceBid(ctx, req)
	if err != nil {
		if err := db.GetHandler().UpdateBuyoutRequest(ctx, &model.BidRequest{
			OrderID:      req.YsRefID,
			MaxBid:       int64(req.Price),
			Status:       "FAILED",
			ErrorMessage: err.Error(),
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}
		return nil, err
	}

	// save auction item when WIN_BID
	go func() {
		auctionItem := item.ToBidAuctionItem()
		auctionItem.BidRequestID = req.YsRefID
		auctionItem.ItemType = req.TransactionType
		if err := db.GetHandler().InsertBidAuctionItem(ctx, auctionItem); err != nil {
			hlog.CtxErrorf(ctx, "insert auction item failed: %+v", err)
		}
	}()

	if err := db.GetHandler().UpdateBuyoutRequest(ctx, &model.BidRequest{
		OrderID:       req.YsRefID,
		Status:        "WIN_BID",
		TransactionID: placeBidResp.ResultSet.Result.TransactionId,
		MaxBid:        int64(req.Price),
	}); err != nil {
		hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
	}

	return &placeBidResp.ResultSet.Result, nil
}

func Bid(ctx context.Context, req *yahoo.PlaceBidRequest, item *yahoo.AuctionItemDetail) (resp *yahoo.PlaceBidResult, err error) {
	// Place bid
	// 1. check if the bid request for this item is already exists
	cBid, err := db.GetHandler().GetBidRequestByAuctionID(ctx, req.AuctionID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get bid request failed: %+v", err)
		return nil, err
	}
	// if the bid reuqest for this item is already exists, compare the max bid price
	outBidOrderId := ""
	if cBid != nil {
		// compare the max bid price
		// if the max bid price is higher than the current max bid price, out bid the current max bid price
		if cBid.MaxBid < int64(req.Price) {
			// out bid this bid request
			outBidOrderId = cBid.OrderID
		} else {
			outBidOrderId = req.YsRefID
		}
		if err := OutBid(ctx, outBidOrderId); err != nil {
			hlog.CtxErrorf(ctx, "out bid failed: %+v", err)
			return nil, err
		}
	}
	// if not, out bid this bid request [LOST_BID]
	if outBidOrderId != req.YsRefID {
		return AddBid(ctx, req, item)
	} else {
		// TODO: return out bid error code
		return nil, errors.New("OUT BID")
	}

}

func AddBid(ctx context.Context, req *yahoo.PlaceBidRequest, item *yahoo.AuctionItemDetail) (resp *yahoo.PlaceBidResult, err error) {
	nextBidPrice := item.BidInfo.NextBid.Price
	// place bid
	placeBidResp, err := yahoo.GetClient().PlaceBid(ctx, &yahoo.PlaceBidRequest{
		YahooAccountID:  config.DevYahoo02AccountID,
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           nextBidPrice,
		Signature:       req.Signature,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid failed: %+v", err)
		return nil, err
	}

	// update bid request status to [BID_PROCESSING]
	if err := db.GetHandler().AddBidRequest(ctx, &model.YahooTransaction{
		BidRequestID:  req.YsRefID,
		Price:         int64(nextBidPrice),
		Status:        "BIDDING_IN_PROGRESS",
		TransactionID: placeBidResp.ResultSet.Result.TransactionId,
	}); err != nil {
		hlog.CtxErrorf(ctx, "add bid request failed: %+v", err)
		return nil, err
	}
	return &placeBidResp.ResultSet.Result, nil
}

func OutBid(ctx context.Context, bidId string) error {
	// update the bid request status to [LOST_BID]
	if err := db.GetHandler().OutBidRequest(ctx, bidId); err != nil {
		hlog.CtxErrorf(ctx, "out bid request failed: %+v", err)
		return err
	}
	// TODO: send message to mq
	go func() {
		//

	}()

	return nil
}

func WinBid(ctx context.Context, bidId string) error {
	return nil
}

func NotifyB4UBiddingResult(ctx context.Context, result *model.BiddingResult) error {
	jsonBody, err := json.Marshal(result)
	if err != nil {
		return err
	}

	headers := amqp.Table{}
	// TODO: generate batch number
	headers["x-batch-number"] = uuid.New().String()

	return mq.SendMessage(mq.Message{
		Exchange:   config.RetryExchange,
		RoutingKey: config.RetryRoutingKey,
		Publishing: amqp.Publishing{
			Body:    jsonBody,
			Headers: headers,
		},
	})
}
