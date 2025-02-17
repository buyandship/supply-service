package mercari

import (
	"context"
	"encoding/json"
	"errors"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
)

func PostOrderService(ctx context.Context, req *supply.MercariPostOrderReq) (*mercari.PurchaseItemResponse, error) {
	hlog.CtxInfof(ctx, "PostOrderService is called, req: %+v", req)

	if req.GetItemID() == "" {
		hlog.CtxInfof(ctx, "empty item_id")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetChecksum() == "" {
		hlog.CtxInfof(ctx, "empty checksum")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetRefID() == "" {
		hlog.CtxInfof(ctx, "empty ref_id")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetRefCurrency() == "" {
		hlog.CtxInfof(ctx, "empty ref_currency")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetRefPrice() == 0 {
		hlog.CtxInfof(ctx, "empty ref_price")
		return nil, bizErr.InvalidParameterError
	}

	h := mercari.GetHandler()
	var buyerId int32 = 1
	if req.GetBuyerID() != 0 {
		buyerId = req.GetBuyerID()
	}
	// check buyer_id
	acc, err := db.GetHandler().GetAccount(ctx, buyerId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, bizErr.InvalidBuyerError
	}
	if err != nil {
		hlog.CtxErrorf(ctx, "Get account error: %s", err.Error())
		return nil, bizErr.InternalError
	}

	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:     req.GetItemID(),
		BuyerId:    req.GetBuyerID(),
		Prefecture: acc.Prefecture,
	})
	if err != nil {
		return nil, err
	}
	// check ref_id
	ok, err := db.GetHandler().CheckTransactionExist(ctx, req.GetRefID())
	if err != nil {
		hlog.CtxErrorf(ctx, "CheckTransactionExist error: %s", err.Error())
		return nil, bizErr.InternalError
	}
	if ok {
		hlog.CtxErrorf(ctx, "duplicated ref_id: %s", req.GetRefID())
		return nil, bizErr.RefIdDuplicatedError
	}

	// 2. Validation
	// 2.1 check checksum
	if resp.Checksum != req.GetChecksum() {
		hlog.CtxErrorf(ctx, "invalid checksum, [%s]!=[%s]", resp.Checksum, req.GetChecksum())
		return nil, bizErr.InvalidCheckSumError
	}
	// 2.2 check price
	if float64(req.GetRefPrice()) < float64(resp.Price)*0.49 {
		hlog.CtxErrorf(ctx, "invalid ref price [%d]", req.GetRefPrice())
		return nil, bizErr.TooLowReferencePriceError
	}
	// 2.3 check the status
	if resp.Status != "on_sale" {
		hlog.CtxErrorf(ctx, "invalid status [%s]", resp.Status)
		return nil, bizErr.ItemNotOnSaleError
	}
	jsonItemDetail, err := json.Marshal(resp)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to marshal json item detail [%s]", err.Error())
		return nil, bizErr.InternalError
	}
	if err := db.GetHandler().InsertTransaction(ctx, &model.Transaction{
		RefID:      req.GetRefID(),
		ItemID:     req.GetItemID(),
		ItemType:   resp.ItemType,
		ItemDetail: jsonItemDetail,
		BuyerID:    buyerId,
		RefPrice:   req.GetRefPrice(),
		Checksum:   req.GetChecksum(),
		Currency:   req.GetRefCurrency(),
	}); err != nil {
		return nil, bizErr.InternalError
	}
	var couponId int
	if resp.ItemDiscount.CouponId != 0 {
		couponId = resp.ItemDiscount.CouponId
	}
	r, err := h.PurchaseItem(ctx, &mercari.PurchaseItemRequest{
		BuyerId:            req.GetBuyerID(),
		ItemId:             req.GetItemID(),
		FamilyName:         acc.FamilyName,
		FirstName:          acc.FirstName,
		FamilyNameKana:     acc.FamilyNameKana,
		FirstNameKana:      acc.FirstNameKana,
		Telephone:          acc.Telephone,
		ZipCode1:           acc.ZipCode1,
		ZipCode2:           acc.ZipCode2,
		Prefecture:         acc.Prefecture,
		City:               acc.City,
		Address1:           acc.Address1,
		Address2:           acc.Address2,
		Checksum:           resp.Checksum,
		CouponId:           couponId,
		DeliveryIdentifier: req.GetItemID(),
	})
	if err != nil {
		if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
			RefID:         req.GetRefID(),
			FailureReason: err.Error(),
		}); err != nil {
			hlog.CtxErrorf(ctx, "UpdateTransaction fail, [%s]", err.Error())
		}
		return nil, err
	}
	if err := db.GetHandler().UpdateTransaction(ctx, &model.Transaction{
		RefID:     req.GetRefID(),
		TrxID:     r.TransactionDetails.TrxId,
		PaidPrice: r.TransactionDetails.PaidPrice,
		Price:     r.TransactionDetails.Price,
	}); err != nil {
		hlog.CtxErrorf(ctx, "UpdateTransaction fail, [%s]", err.Error())
	}
	r.CouponId = couponId
	return r, err
}
