package yahoo

import (
	"context"

	globalConfig "github.com/buyandship/bns-golib/config"
	"github.com/buyandship/supply-service/biz/common/config"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/mock"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// TODO throttle and rate limit.
func PlaceBidService(ctx context.Context, req *supply.YahooPlaceBidReq) (resp *yahoo.PlaceBidResult, err error) {

	if err := mock.MockYahooPlaceBidError(req.AuctionID); err != nil {
		return nil, err
	}

	if req.YsRefID == "" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "ys_ref_id is required",
		}
	}

	client := yahoo.GetClient()
	// TODO: idempotent check.
	tx, err := db.GetHandler().GetBidRequest(ctx, req.YsRefID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get buyout bid request failed: %+v", err)
		return nil, bizErr.BizError{
			Status:  consts.StatusInternalServerError,
			ErrCode: consts.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}

	if tx != nil {
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity,
			ErrMsg:  "order already exists",
		}
	}

	if req.TransactionType != "BID" && req.TransactionType != "BUYOUT" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "transaction_type must be BID or BUYOUT",
		}
	}

	if req.TransactionType == "BID" {
		return nil, bizErr.BizError{
			Status:  consts.StatusBadRequest,
			ErrCode: consts.StatusBadRequest,
			ErrMsg:  "BID type is not supported yet",
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

	// get auction item
	auctionItemResp, err := client.GetAuctionItem(ctx, yahoo.AuctionItemRequest{AuctionID: req.AuctionID})
	if err != nil {
		// Auction item not found
		return nil, bizErr.BizError{
			Status:  consts.StatusNotFound,
			ErrCode: consts.StatusNotFound,
			ErrMsg:  "get auction item failed",
		}
	}

	if globalConfig.GlobalAppConfig.Env == "dev" {
		if auctionItemResp.ResultSet.Result.Seller.AucUserId != "AnzTKsBM5HUpBc3CCQc3dHpETkds1" {
			return nil, bizErr.BizError{
				Status:  consts.StatusUnprocessableEntity,
				ErrCode: consts.StatusUnprocessableEntity,
				ErrMsg:  "this product is not allowed to be placed bid in staging environment",
			}
		}
	}

	// TODO: validation.
	item := auctionItemResp.ResultSet.Result
	if item.Status != "open" {
		// TODO: return Auction Item is not available
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "The auction item is not available",
		}
	}
	if req.TransactionType == "BUYOUT" {
		if item.TaxinBidorbuy != 0 {
			if req.Price != int32(item.TaxinBidorbuy) {
				return nil, bizErr.BizError{
					Status:  consts.StatusUnprocessableEntity,
					ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
					ErrMsg:  "The request price is not same as Buyout price",
				}
			}
		} else {
			if req.Price != int32(item.Bidorbuy) {
				return nil, bizErr.BizError{
					Status:  consts.StatusUnprocessableEntity,
					ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
					ErrMsg:  "The request price is not same as Buyout price",
				}
			}
		}
	}
	if !req.Partial && item.Quantity < int(req.Quantity) {
		// TODO: Requested Quantity is not able to fulfil
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity, // TODO: define error code
			ErrMsg:  "The requested quantity is not able to fulfill",
		}
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
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity,
			ErrMsg:  "order already exists",
		}
	}

	// check if the current price is highest in this auction
	// place bid preview
	previewReq := &yahoo.PlaceBidPreviewRequest{
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           int(req.Price),
		Quantity:        int(req.Quantity),
		Partial:         false,
	}
	previewResp, err := client.PlaceBidPreview(ctx, previewReq)
	if err != nil {
		hlog.CtxErrorf(ctx, "place bid preview failed: %+v", err)
		// TODO: update order status to FAILED
		if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
			OrderID:       req.YsRefID,
			Status:        "FAILED",
			MaxBid:        int64(req.Price),
			ErrorMessage:  err.Error(),
			TransactionID: previewResp.ResultSet.Result.TransactionId,
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}
		return nil, bizErr.BizError{
			Status:  consts.StatusUnprocessableEntity,
			ErrCode: consts.StatusUnprocessableEntity,
			ErrMsg:  "place bid preview failed",
		}
	}

	if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
		OrderID:       req.YsRefID,
		TransactionID: previewResp.ResultSet.Result.TransactionId,
		Status:        "CREATED",
		MaxBid:        int64(req.Price),
	}); err != nil {
		hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
			OrderID:      req.YsRefID,
			Status:       "FAILED",
			MaxBid:       int64(req.Price),
			ErrorMessage: err.Error(),
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}
		return nil, bizErr.BizError{
			Status:  consts.StatusInternalServerError,
			ErrCode: consts.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}

	var yahooAccountID string
	switch globalConfig.GlobalAppConfig.Env {
	case "dev":
		yahooAccountID = config.DevYahoo02AccountID
	case "prod":
		yahooAccountID = config.ProdYahoo02AccountID
	}

	// TODO: check if it's neccessary to update the bid request in database.
	// TODO: determine if it's shopping item.
	bidReq := yahoo.PlaceBidRequest{
		YahooAccountID:  yahooAccountID,
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           int(req.Price),
		Quantity:        int(req.Quantity),
		Partial:         req.Partial,
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
			ErrMsg:  "BID type is not supported yet",
		}
	case "BUYOUT":
		// directly buyout
		placeBidResp, err := client.PlaceBid(ctx, &bidReq)
		if err != nil {
			if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
				OrderID:      req.YsRefID,
				MaxBid:       int64(req.Price),
				Status:       "FAILED",
				ErrorMessage: err.Error(),
			}); err != nil {
				hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
			}
			return nil, bizErr.BizError{
				Status:  consts.StatusUnprocessableEntity,
				ErrCode: consts.StatusUnprocessableEntity,
				ErrMsg:  err.Error(),
			}
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

		if err := db.GetHandler().UpdateBuyoutBidRequest(ctx, &model.BidRequest{
			OrderID:       req.YsRefID,
			Status:        "WIN_BID",
			TransactionID: placeBidResp.ResultSet.Result.TransactionId,
			MaxBid:        int64(req.Price),
		}); err != nil {
			hlog.CtxErrorf(ctx, "update yahoo order failed: %+v", err)
		}

		return &placeBidResp.ResultSet.Result, nil
	}

	return nil, nil
}
