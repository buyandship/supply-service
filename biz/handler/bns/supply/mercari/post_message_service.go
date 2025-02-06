package mercari

import (
	"context"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
)

func PostMessageService(ctx context.Context, req *supply.MercariPostMessageReq) (*supply.MercariPostMessageResp, error) {
	h := mercari.GetHandler()

	_, err := h.PostTransactionMessage(ctx, &mercari.PostTransactionMessageRequest{
		TransactionId: req.GetTrxID(),
		BuyerId:       req.GetBuyerID(),
		Message:       req.GetMsg(),
	})

	if err := db.GetHandler().InsertMessage(ctx, &model.Message{
		TrxID:   req.GetTrxID(),
		Message: req.GetMsg(),
		BuyerID: req.GetBuyerID(),
	}); err != nil {
		return nil, err
	}

	return &supply.MercariPostMessageResp{}, err
}
