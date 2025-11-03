package yahoo

import (
	"context"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func PlaceBidService(ctx context.Context, req *supply.YahooPlaceBidReq) (resp *supply.YahooPlaceBidResp, err error) {

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
		return nil, err
	}

	bidReq := yahoo.PlaceBidRequest{
		YahooAccountID:  "TODO",
		YsRefID:         req.YsRefID,
		TransactionType: req.TransactionType,
		AuctionID:       req.AuctionID,
		Price:           int(req.Price),
		Quantity:        int(req.Quantity),
		Partial:         false,
		Signature:       previewResp.ResultSet.Result.Signature,
	}
	placeBidResp, err := client.MockPlaceBid(ctx, &bidReq)
	if err != nil {
		return nil, err
	}

	resp = &supply.YahooPlaceBidResp{
		Status:     placeBidResp.Result.Status,
		BidID:      placeBidResp.Result.BidID,
		AuctionID:  placeBidResp.Result.AuctionID,
		Price:      int32(placeBidResp.Result.Price),
		Quantity:   int32(placeBidResp.Result.Quantity),
		TotalPrice: int32(placeBidResp.Result.TotalPrice),
		BidTime:    placeBidResp.Result.BidTime,
	}

	return resp, nil
}
