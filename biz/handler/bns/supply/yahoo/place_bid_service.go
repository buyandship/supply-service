package yahoo

import (
	"context"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// TODO throttle and rate limit.
func PlaceBidService(ctx context.Context, req *supply.YahooPlaceBidReq) (resp *yahoo.PlaceBidResult, err error) {

	if req.YsRefID == "" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "ys_ref_id is required",
		}
	}

	if req.TransactionType != "BID" && req.TransactionType != "BUYOUT" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "transaction_type must be BID or BUYOUT",
		}
	}

	if req.AuctionID == "" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "auction_id is required",
		}
	}

	if req.Price <= 0 {
		// for BID type, price is the max bid price
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "price must be greater than 0",
		}
	}

	if req.Quantity <= 0 {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "quantity must be greater than 0",
		}
	}

	client := yahoo.GetClient()

	// get auction item
	auctionItemResp, err := client.MockGetAuctionItem(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID})
	if err != nil {
		// Auction item not found
		return nil, err
	}
	// TODO: validation.
	item := auctionItemResp.ResultSet.Result
	if item.Status != "OPEN" {
		// TODO: return Auction Item is not available
	}
	if req.TransactionType == "BUYOUT" && req.Price != int32(item.Bidorbuy) {
		// TODO: return Request price is not same as Buyout price
	}
	if !req.Partial && item.Quantity < int(req.Quantity) {
		// TODO: Requested Quantity is not able to fulfil
	}

	// save order into database
	order := &model.BidRequest{
		RequestType: req.TransactionType,
		OrderID:     req.YsRefID,
		AuctionID:   req.AuctionID,
		MaxBid:      int64(req.Price),
		Quantity:    int32(req.Quantity),
		Partial:     false,
		Status:      "CREATED",
	}
	if err := db.GetHandler().InsertBuyoutBidRequest(ctx, order); err != nil {
		hlog.CtxErrorf(ctx, "insert yahoo order failed: %+v", err)
		return nil, err
	}

	if order.Status != "CREATED" {
		hlog.CtxErrorf(ctx, "order already exists: %+v", order.Status)
		// order already exists
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "order already exists",
		}
	}

	// check if the current price is highest in this auction
	// place bid preview
	previewReq := &yahoo.PlaceBidPreviewRequest{
		YahooAccountID:  "TODO",
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           int(req.Price),
		Quantity:        int(req.Quantity),
		Partial:         false,
	}
	previewResp, err := client.MockPlaceBidPreview(ctx, previewReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid preview failed: %+v", err)
		// TODO: update order status to FAILED
		if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
			OrderID:      req.YsRefID,
			Status:       "FAILED",
			ErrorMessage: err.Error(),
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}
		return nil, err
	}

	// TODO: check if it's neccessary to update the bid request in database.
	// TODO: determine if it's shopping item.
	bidReq := yahoo.PlaceBidRequest{
		YahooAccountID:  "TODO",
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           int(req.Price),
		Quantity:        int(req.Quantity),
		Partial:         false,
		Signature:       previewResp.ResultSet.Result.Signature,
		// TODO: IsShoppingItem
	}

	// check if it's buyout
	switch req.TransactionType {
	case "BID":
		// check if the current price is highest in this auction
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "BID type is not supported",
		}
	case "BUYOUT":
		// directly buyout
		placeBidResp, err := client.MockPlaceBid(ctx, &bidReq)
		if err != nil {
			// TODO: update order status to FAILED
			if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
				OrderID:      req.YsRefID,
				Status:       "FAILED",
				ErrorMessage: err.Error(),
			}); err != nil {
				hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
			}
			return nil, bizErr.BizError{
				Status:  consts.StatusBadRequest,
				ErrCode: consts.StatusBadRequest,
				ErrMsg:  err.Error(),
			}
		}

		// TODO: update bid request to WIN_BID
		if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
			OrderID: req.YsRefID,
			Status:  "WIN_BID",
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}

		return &placeBidResp.ResultSet.Result, nil
	}

	return nil, nil
}
