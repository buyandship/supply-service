package mercari

import (
	"context"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/mercari"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func PostTransactionReviewService(ctx context.Context, req *supply.MercariPostTransactionReviewReq) (*mercari.PostTransactionReviewResponse, error) {
	hlog.CtxInfof(ctx, "PostTransactionReviewService is called, %+v", req)
	h := mercari.GetHandler()

	/* 	if config.GlobalServerConfig.Env == "dev" {
		return nil, bizErr.NotFoundError
	} */

	/* 	if config.GlobalServerConfig.Env == "dev" {
		return &mercari.PostTransactionReviewResponse{
			ReviewStatus: "success",
			RequestId:    uuid.NewString(),
		}, nil
	} */

	if req.GetTrxID() == "" {
		hlog.CtxWarnf(ctx, "empty trx_id")
		return nil, bizErr.InvalidParameterError
	}

	if req.GetFame() != "good" && req.GetFame() != "bad" {
		hlog.CtxErrorf(ctx, "invalid fame: %s", req.GetFame())
		return nil, bizErr.InvalidParameterError
	}

	if len(req.GetReview()) >= 140 {
		hlog.CtxErrorf(ctx, "the length of review reach maximum: %d", len(req.GetReview()))
		return nil, bizErr.InvalidParameterError
	}

	trx, err := db.GetHandler().GetTransaction(ctx, &model.Transaction{
		TrxID: req.GetTrxID(),
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "get transaction failed: %v", err)
		return nil, bizErr.NotFoundError
	}

	mResp, err := h.PostTransactionReview(ctx, &mercari.PostTransactionReviewRequest{
		TrxId:     req.GetTrxID(),
		Fame:      req.GetFame(),
		Message:   req.GetReview(),
		AccountID: trx.AccountID,
	})
	if err != nil {
		return nil, err
	}

	if err := db.GetHandler().InsertReview(ctx, &model.Review{
		TrxID:     req.GetTrxID(),
		Fame:      req.GetFame(),
		Review:    req.GetReview(),
		AccountID: trx.AccountID,
	}); err != nil {
		hlog.CtxErrorf(ctx, "Insert review failed: %v", err)
		return nil, bizErr.InternalError
	}
	return mResp, nil
}
