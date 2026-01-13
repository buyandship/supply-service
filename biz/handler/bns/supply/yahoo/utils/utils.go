package utils

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/buyandship/bns-golib/cache"
	GlobalConfig "github.com/buyandship/bns-golib/config"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/common/consts"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/mq"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/mercari"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
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

func generateHmac(bidId string) string {
	hmac := hmac.New(sha256.New, []byte(GlobalConfig.GlobalAppConfig.GetString("b4u_secret")))
	hmac.Write([]byte(bidId))
	return hex.EncodeToString(hmac.Sum(nil))
}

func GetAuctionItem(ctx context.Context, auctionID string) (*yahoo.AuctionItemResponse, error) {
	auctionItemResp := &yahoo.AuctionItemResponse{}

	client := yahoo.GetClient()
	if GlobalConfig.GlobalAppConfig.Env == "dev" && auctionID == "bravo_test_item" {
		auctionItemResp = mock.TestAuction()
	} else {
		resp, err := client.GetAuctionItemAuth(ctx, yahoo.AuctionItemRequest{AuctionID: auctionID})
		if err != nil {
			return nil, err
		}
		auctionItemResp = resp
	}

	// TODO: update bid request status
	go func() {
		// get bid request
		bidRequest, err := db.GetHandler().GetBidRequestByAuctionID(ctx, auctionID, "")
		if err != nil {
			hlog.CtxErrorf(ctx, "get bid request failed: %+v", err)
			return
		}
		if bidRequest == nil {
			return
		}
		// TODO: check win bid
		if auctionItemResp != nil && auctionItemResp.ResultSet.Result.WinnersInfo != nil && auctionItemResp.ResultSet.Result.WinnersInfo.Winner != nil {
			if auctionItemResp.ResultSet.Result.WinnersInfo.Winner[0].AucUserId == "AnzTKsBM5HUpBc3CCQc3dHpETkds1" { // TODO: change to list
				// Win bid
				if err := WinBid(ctx, bidRequest.OrderID, int64(auctionItemResp.ResultSet.Result.WinnersInfo.Winner[0].WonPrice)); err != nil {
					hlog.CtxErrorf(ctx, "win bid failed: %+v", err)
					return
				}
			} else {
				// OUTBID
				if err := OutBid(ctx, bidRequest.OrderID); err != nil {
					hlog.CtxErrorf(ctx, "out bid failed: %+v", err)
					return
				}
			}
			return
		}

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
		Status:        model.StatusWinBid,
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
	cBid, err := db.GetHandler().GetBidRequestByAuctionID(ctx, req.AuctionID, req.YsRefID)
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
		hlog.CtxDebugf(ctx, "out bid order id: %s, cBid.MaxBid: %d, req.Price: %d", outBidOrderId, cBid.MaxBid, req.Price)
		if err := OutBid(ctx, outBidOrderId); err != nil {
			hlog.CtxErrorf(ctx, "out bid failed: %+v", err)
			return nil, err
		}
	}
	// if not, out bid this bid request [LOST_BID]
	if outBidOrderId != req.YsRefID {
		// TBC: the user maybe confused why it out bid.
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
		TransactionType: consts.TransactionTypeBid,
		AuctionID:       req.AuctionID,
		Price:           nextBidPrice,
		Signature:       req.Signature,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid failed: %+v", err)
		return nil, err
	}

	// update the next bid price
	mock.UpdateNextBidPrice()

	// update bid request status to [BID_PROCESSING]
	if err := db.GetHandler().AddBidRequest(ctx, &model.YahooTransaction{
		BidRequestID:  req.YsRefID,
		Price:         int64(nextBidPrice),
		Status:        model.StatusBiddingInProgress,
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
		result := &model.BiddingResult{
			OrderNumber: bidId,
			Status:      model.StatusLostBid,
		}
		if err := notifyB4UBiddingResult(ctx, result); err != nil {
			hlog.CtxErrorf(ctx, "send message failed: %+v", err)
		}
	}()

	return nil
}

func WinBid(ctx context.Context, bidId string, wonPrice int64) error {
	if err := db.GetHandler().WinBidRequest(ctx, bidId, wonPrice); err != nil {
		hlog.CtxErrorf(ctx, "win bid request failed: %+v", err)
		return err
	}
	go func() {
		result := &model.BiddingResult{
			MetaInfo: model.BiddingMetaInfo{
				WonPrice: wonPrice,
			},
			OrderNumber: bidId,
			Status:      model.StatusWinBid,
		}
		if err := notifyB4UBiddingResult(ctx, result); err != nil {
			hlog.CtxErrorf(ctx, "notifyB4UBiddingResulte failed: %+v", err)
		}
	}()
	return nil
}

func notifyB4UBiddingResult(ctx context.Context, result *model.BiddingResult) error {
	result.Hmac = generateHmac(result.OrderNumber)
	jsonBody, err := json.Marshal(result)
	if err != nil {
		return err
	}
	headers := amqp.Table{}
	headers["x-order-number"] = result.OrderNumber

	return mq.SendMessage(mq.Message{
		Exchange:   config.RetryExchange,
		RoutingKey: config.RetryRoutingKey,
		Publishing: amqp.Publishing{
			Body:    jsonBody,
			Headers: headers,
		},
	})
}
