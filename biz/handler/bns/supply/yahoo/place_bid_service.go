package yahoo

import (
	"context"

	globalConfig "github.com/buyandship/bns-golib/config"
	"github.com/buyandship/supply-service/biz/common/consts"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/handler/bns/supply/yahoo/utils"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	httpConsts "github.com/cloudwego/hertz/pkg/protocol/consts"
)

// TODO throttle and rate limit.
func PlaceBidService(ctx context.Context, req *supply.YahooPlaceBidReq) (resp *yahoo.PlaceBidResult, err error) {

	if err := mock.MockYahooPlaceBidError(req.AuctionID); err != nil {
		return nil, err
	}

	if req.YsRefID == "" {
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusBadRequest,
			ErrCode: httpConsts.StatusBadRequest,
			ErrMsg:  "ys_ref_id is required",
		}
	}

	tx, err := db.GetHandler().GetBidRequest(ctx, req.YsRefID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get buyout bid request failed: %+v", err)
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusInternalServerError,
			ErrCode: httpConsts.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}

	if tx != nil {
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusUnprocessableEntity,
			ErrCode: httpConsts.StatusUnprocessableEntity,
			ErrMsg:  "order already exists",
		}
	}

	if req.TransactionType != consts.TransactionTypeBid && req.TransactionType != consts.TransactionTypeBuyout {
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusBadRequest,
			ErrCode: httpConsts.StatusBadRequest,
			ErrMsg:  "transaction_type must be BID or BUYOUT",
		}
	}

	if req.AuctionID == "" {
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusBadRequest,
			ErrCode: httpConsts.StatusBadRequest,
			ErrMsg:  "auction_id is required",
		}
	}

	if req.Price <= 0 {
		// for BID type, price is the max bid price
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusBadRequest,
			ErrCode: httpConsts.StatusBadRequest,
			ErrMsg:  "price must be greater than 0",
		}
	}

	if req.Quantity <= 0 {
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusBadRequest,
			ErrCode: httpConsts.StatusBadRequest,
			ErrMsg:  "quantity must be greater than 0",
		}
	}

	// get auction item
	auctionItemResp, err := utils.GetAuctionItem(ctx, req.AuctionID)
	if err != nil {
		// Auction item not found
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusNotFound,
			ErrCode: httpConsts.StatusNotFound,
			ErrMsg:  "get auction item failed",
		}
	}

	if globalConfig.GlobalAppConfig.Env == "dev" {
		if auctionItemResp.ResultSet.Result.Seller.AucUserId != "AnzTKsBM5HUpBc3CCQc3dHpETkds1" {
			return nil, bizErr.BizError{
				Status:  httpConsts.StatusUnprocessableEntity,
				ErrCode: httpConsts.StatusUnprocessableEntity,
				ErrMsg:  "this product is not allowed to be placed bid in staging environment",
			}
		}
	}

	// TODO: validation.
	item := auctionItemResp.ResultSet.Result
	if item.Status != "open" {
		// TODO: return Auction Item is not available
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusUnprocessableEntity,
			ErrCode: httpConsts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "The auction item is not available",
		}
	}

	price := 0
	if req.TransactionType == consts.TransactionTypeBuyout {
		if item.TaxinBidorbuy != 0 {
			if req.Price != int32(item.TaxinBidorbuy) {
				return nil, bizErr.BizError{
					Status:  httpConsts.StatusUnprocessableEntity,
					ErrCode: httpConsts.StatusUnprocessableEntity, // TODO: define error code
					ErrMsg:  "The request price is not same as Buyout price",
				}
			}
		} else {
			if req.Price != int32(item.Bidorbuy) {
				return nil, bizErr.BizError{
					Status:  httpConsts.StatusUnprocessableEntity,
					ErrCode: httpConsts.StatusUnprocessableEntity, // TODO: define error code
					ErrMsg:  "The request price is not same as Buyout price",
				}
			}
		}
		price = int(item.Bidorbuy)
	} else {
		price = int(req.Price) // TBC: should display next bid price?
		if int(req.Price) < int(item.BidInfo.NextBid.Price) {
			return nil, bizErr.BizError{
				Status:  httpConsts.StatusUnprocessableEntity,
				ErrCode: httpConsts.StatusUnprocessableEntity, // TODO: define error code
				ErrMsg:  "The request price is not greater than the current price",
			}
		}
	}

	if !req.Partial && item.Quantity < int(req.Quantity) {
		// TODO: Requested Quantity is not able to fulfil
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusUnprocessableEntity,
			ErrCode: httpConsts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "The requested quantity is not able to fulfill",
		}
	}

	// save order into database
	order := &model.BidRequest{
		RequestType: req.TransactionType,
		OrderID:     req.YsRefID,
		AuctionID:   req.AuctionID,
		MaxBid:      int64(price),
		Quantity:    int32(req.Quantity),
		Partial:     false,
		Status:      model.StatusCreated,
	}
	if err := db.GetHandler().InsertBidRequest(ctx, order); err != nil {
		hlog.CtxErrorf(ctx, "insert yahoo order failed: %+v", err)
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

	if order.Status != "CREATED" {
		hlog.CtxErrorf(ctx, "order already exists: %+v", order.Status)
		// order already exists
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusUnprocessableEntity,
			ErrCode: httpConsts.StatusUnprocessableEntity,
			ErrMsg:  "order already exists",
		}
	}

	// check if the current price is highest in this auction
	// place bid preview
	previewReq := &yahoo.PlaceBidPreviewRequest{
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           price,
		Quantity:        int(req.Quantity),
		Partial:         false,
	}
	previewResp, err := yahoo.GetClient().PlaceBidPreview(ctx, previewReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid preview failed: %+v", err)
		// TODO: update order status to FAILED
		if err := db.GetHandler().UpdateBuyoutRequest(ctx, &model.BidRequest{
			OrderID:       req.YsRefID,
			Status:        "FAILED",
			MaxBid:        int64(req.Price),
			ErrorMessage:  err.Error(),
			TransactionID: previewResp.ResultSet.Result.TransactionId,
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}
		return nil, bizErr.BizError{
			Status:  httpConsts.StatusUnprocessableEntity,
			ErrCode: httpConsts.StatusUnprocessableEntity,
			ErrMsg:  "place bid preview failed",
		}
	}

	/*
		if err := db.GetHandler().UpdateBuyoutRequest(ctx, &model.BidRequest{
			OrderID:       req.YsRefID,
			TransactionID: previewResp.ResultSet.Result.TransactionId,
			Status:        model.StatusCreated,
			MaxBid:        int64(req.Price),
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
			if err := db.GetHandler().UpdateBuyoutRequest(ctx, &model.BidRequest{
				OrderID:      req.YsRefID,
				Status:       "FAILED",
				MaxBid:       int64(req.Price),
				ErrorMessage: err.Error(),
			}); err != nil {
				hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
			}
			return nil, bizErr.BizError{
				Status:  httpConsts.StatusInternalServerError,
				ErrCode: httpConsts.StatusInternalServerError,
				ErrMsg:  "internal server error",
			}
		}
	*/

	// TODO: check if it's neccessary to update the bid request in database.
	bidReq := yahoo.PlaceBidRequest{
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           price,
		Quantity:        int(req.Quantity),
		Partial:         req.Partial,
		Signature:       previewResp.ResultSet.Result.Signature,
		IsShoppingItem:  isShoppingItem(&item),
	}

	// check if it's buyout
	switch req.TransactionType {
	case consts.TransactionTypeBid:
		// check if the current price is highest in this auction
		return utils.Bid(ctx, &bidReq, &item)
	case consts.TransactionTypeBuyout:
		return utils.Buyout(ctx, &bidReq, &item)
	}

	return nil, nil
}

func isShoppingItem(item *yahoo.AuctionItemDetail) bool {
	if item == nil {
		return false
	}
	if item.ShoppingItemCode != "" && item.ShoppingItem != nil && !item.ShoppingItem.IsOptionEnabled {
		return true
	}
	return false
}
