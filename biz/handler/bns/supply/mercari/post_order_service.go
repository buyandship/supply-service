package mercari

import (
	"context"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostOrderService(ctx context.Context, req *supply.MercariPostOrderReq) (*supply.MercariPostOrderResp, error) {
	h := mercari.GetHandler()
	// 1. GetItem
	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:  req.GetItemID(),
		BuyerId: req.GetBuyerID(),
	})
	if err != nil {
		return nil, err
	}
	// 2. Validation
	// 2.1 check checksum
	if resp.Checksum != req.GetChecksum() {
		hlog.Errorf("invalid checksum, [%s]!=[%s]", resp.Checksum, req.GetChecksum())
		return nil, bizErr.InvalidParameterError
	}
	// 2.2 check price
	if float64(req.GetRefPrice()) < float64(resp.Price)*0.49 {
		hlog.Errorf("invalid ref price [%d]", req.GetRefPrice())
		return nil, bizErr.TooLowReferencePriceError
	}
	// 2.3 check the status
	if resp.Status != "on_sale" {
		hlog.Errorf("invalid status [%s]", resp.Status)
		return nil, bizErr.ItemNotOnSaleError
	}

	// TODO insert into db
	if err := db.GetHandler().InsertTransaction(ctx, &model.Transaction{
		RefID:    req.GetRefID(),
		ItemID:   req.GetItemID(),
		ItemType: "", // TBC
		// ItemDetail: "",  // TBC
		BuyerID:   req.GetBuyerID(),
		Price:     0, // TBC
		PaidPrice: 0, // TBC
		RefPrice:  int(req.GetRefPrice()),
		Checksum:  req.GetChecksum(),
	}); err != nil {

	}

	r, err := h.PurchaseItem(ctx, &mercari.PurchaseItemRequest{
		BuyerId: req.GetBuyerID(),
		ItemId:  req.GetItemID(),
	})
	if err != nil {
		if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
			RefID:         req.GetRefID(),
			FailureReason: err.Error(),
		}); err != nil {
			hlog.Errorf("UpdateTransaction fail, [%s]", err.Error())
		}
		return nil, err
	}

	if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
		RefID:     req.GetRefID(),
		TrxID:     r.TransactionDetails.TrxId,     // TBC
		PaidPrice: r.TransactionDetails.PaidPrice, //
	}); err != nil {
		hlog.Errorf("UpdateTransaction fail, [%s]", err.Error())
	}

	return &supply.MercariPostOrderResp{}, err
}
