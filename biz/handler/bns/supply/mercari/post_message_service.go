package mercari

import (
	"context"
	"unicode/utf8"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostMessageService(ctx context.Context, req *supply.MercariPostMessageReq) (*supply.MercariPostMessageResp, error) {
	hlog.CtxInfof(ctx, "PostMessageService is called, req: %+v", req)
	h := mercari.GetHandler()

	/* 	if config.GlobalServerConfig.Env == "development" {
		return nil, bizErr.ConflictError
	} */

	charCount := utf8.RuneCountInString(req.GetMsg())

	if req.GetTrxID() == "" {
		hlog.CtxErrorf(ctx, "empty trx_id")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetMsg() == "" || charCount > 1000 {
		hlog.CtxErrorf(ctx, "msg is empty or length exceeds 1000")
		return nil, bizErr.InvalidParameterError
	}

	trx, err := db.GetHandler().GetTransaction(ctx, &model.Transaction{
		TrxID: req.GetTrxID(),
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Get transaction failed: %v", err)
		return nil, bizErr.InternalError
	}

	mResp, err := h.PostTransactionMessage(ctx, &mercari.PostTransactionMessageRequest{
		TransactionId: req.GetTrxID(),
		Message:       req.GetMsg(),
		AccountID:     trx.AccountID,
	})
	if err != nil {
		return nil, err
	}
	if err := db.GetHandler().InsertMessage(ctx, &model.Message{
		TrxID:   req.GetTrxID(),
		Message: req.GetMsg(),
	}); err != nil {
		hlog.CtxErrorf(ctx, "Insert message failed: %v", err)
		return nil, bizErr.InternalError
	}

	return &supply.MercariPostMessageResp{
		TrxID:     req.GetTrxID(),
		Body:      mResp.Body,
		ID:        mResp.Id,
		AccountID: mResp.AccountID,
	}, nil
}
