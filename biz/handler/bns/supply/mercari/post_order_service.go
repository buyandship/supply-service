package mercari

import (
	"context"
	"encoding/json"
	"strconv"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/handler/bns/supply/utils"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/mock"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const priceThreshold = 0.49
const CountryCodeHK = "HK"

func validateRequest(ctx context.Context, req *supply.MercariPostOrderReq) error {
	if req.GetItemID() == "" {
		hlog.CtxInfof(ctx, "empty item_id")
		return bizErr.InvalidParameterError
	}
	if req.GetChecksum() == "" {
		hlog.CtxInfof(ctx, "empty checksum")
		return bizErr.InvalidParameterError
	}
	if req.GetRefID() == "" {
		hlog.CtxInfof(ctx, "empty ref_id")
		return bizErr.InvalidParameterError
	}
	if req.GetRefCurrency() == "" {
		hlog.CtxInfof(ctx, "empty ref_currency")
		return bizErr.InvalidParameterError
	}
	if req.GetRefPrice() == 0 {
		hlog.CtxInfof(ctx, "empty ref_price")
		return bizErr.InvalidParameterError
	}
	return nil
}

func getResponse(tx *model.Transaction) *supply.MercariPostOrderResp {
	if tx == nil {
		return nil
	}
	var fee int
	if tx.BuyerShippingFee != "" {
		fee, _ = strconv.Atoi(tx.BuyerShippingFee)
	} else {
		fee = 0
	}

	return &supply.MercariPostOrderResp{
		TrxID:            tx.TrxID,
		CouponID:         int64(tx.CouponID),
		Price:            tx.Price,
		PaidPrice:        tx.PaidPrice,
		BuyerShippingFee: int64(fee),
		ItemID:           tx.ItemID,
		Checksum:         tx.Checksum,
	}
}

func PostOrderService(ctx context.Context, req *supply.MercariPostOrderReq) (*supply.MercariPostOrderResp, error) {
	hlog.CtxInfof(ctx, "PostOrderService is called, req: %+v", req)
	// 1. validation
	if err := validateRequest(ctx, req); err != nil {
		return nil, err
	}

	// Mock
	if err := mock.MockMercariPostOrderError(req.GetItemID()); err != nil {
		return nil, err
	}

	h := mercari.GetHandler()

	// 2. get buyer
	acc := &model.Account{}
	acc, err := utils.GetBuyer(ctx, req.GetBuyerID())
	if err != nil {
		hlog.CtxErrorf(ctx, "GetBuyer error: %v", err)
	}

	// 3. get item by item_id
	resp, err := h.GetItemByID(ctx, &mercari.GetItemByIDRequest{
		ItemId:     req.GetItemID(),
		Prefecture: acc.Prefecture,
	})
	if err != nil {
		return nil, err
	}
	var couponId int
	if resp.ItemDiscount.CouponId != 0 {
		couponId = resp.ItemDiscount.CouponId
	}

	var deliveryId string
	if req.GetDeliveryID() == "" {
		deliveryId = req.GetRefID()
	} else {
		deliveryId = req.GetDeliveryID()
	}

	// 4. check if transaction with ref_id already exists.
	tx, err := db.GetHandler().GetTransaction(ctx, req.GetRefID())
	if err != nil {
		hlog.CtxErrorf(ctx, "get transaction error: %s", err.Error())
		return nil, bizErr.InternalError
	}

	if tx == nil {
		// 4.1 the transaction does not exist
		// 4.1.1 check the checksum
		if resp.Checksum != req.GetChecksum() {
			hlog.CtxErrorf(ctx, "invalid checksum, [%s]!=[%s]", resp.Checksum, req.GetChecksum())
			return nil, bizErr.InvalidCheckSumError
		}
		// 4.1.2 check price
		if float64(req.GetRefPrice()) < float64(resp.Price)*priceThreshold {
			hlog.CtxErrorf(ctx, "invalid ref price [%d]", req.GetRefPrice())
			return nil, bizErr.TooLowReferencePriceError
		}
		// 4.1.3 check the status
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
			BuyerID:    acc.BuyerID,
			RefPrice:   req.GetRefPrice(),
			Checksum:   req.GetChecksum(),
			Currency:   req.GetRefCurrency(),
			CouponID:   couponId,
			DeliveryId: deliveryId,
		}); err != nil {
			return nil, bizErr.InternalError
		}
	}

	if tx != nil {
		// 4.2 the transaction does exist.
		if tx.TrxID != "" && tx.TrxID != "0" {
			// 4.2.1 the transaction does exist, and the database record the trx_id, which means this order purchased
			// successfully.
			hlog.CtxInfof(ctx, "transaction:[%s] exists. ", req.GetRefID())
			return getResponse(tx), nil
		}

		if tx.FailureReason != "" {
			// 4.2.2 the transaction does exist, but the failure_reason is not empty, which means this order purchased
			// failure
			return nil, bizErr.BizError{
				Status:  500,
				ErrCode: 500,
				ErrMsg:  tx.FailureReason,
			}
		}

		// 4.2.3 check if the parameters is changed
		if tx.ItemID != req.GetItemID() {
			hlog.CtxErrorf(ctx, "item_id does not match")
			return nil, bizErr.InvalidParameterError
		}
		if tx.RefPrice != req.GetRefPrice() {
			hlog.CtxErrorf(ctx, "ref_price does not match")
			return nil, bizErr.InvalidParameterError
		}
		if tx.Currency != req.GetRefCurrency() {
			hlog.CtxErrorf(ctx, "ref_currency does not match")
			return nil, bizErr.InvalidParameterError
		}

	}

	// 5. purchase item
	if err := h.PurchaseItem(ctx, req.GetRefID(), &mercari.PurchaseItemRequest{
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
		DeliveryIdentifier: deliveryId,
		CountryCode:        CountryCodeHK,
	}); err != nil {
		return nil, err
	}

	tx, err = db.GetHandler().GetTransaction(ctx, req.GetRefID())
	if err != nil {
		hlog.CtxErrorf(ctx, "get transaction error: %s", err.Error())
		return nil, bizErr.InternalError
	}
	return getResponse(tx), err
}
