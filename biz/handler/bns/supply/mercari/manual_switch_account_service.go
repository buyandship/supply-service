package mercari

import (
	"context"
	"time"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrasturcture/db"
	"github.com/buyandship/supply-service/biz/infrasturcture/http"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/buyandship/supply-service/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func ManualSwitchAccountService(ctx context.Context, req *supply.MercariManualSwitchAccountReq) error {

	account, err := db.GetHandler().GetAccount(ctx, req.AccountID)
	if err != nil {
		return err
	}

	if account.BannedAt != nil {
		return bizErr.AccountBannedError
	}

	if err := db.GetHandler().SwitchAccount(ctx, req.AccountID); err != nil {
		return err
	}

	if err := cache.GetRedisClient().Set(ctx, config.ActiveAccountId, req.AccountID, time.Hour); err != nil {
		hlog.CtxErrorf(ctx, "failed to set active account id: %v", err)
	}

	var activeAccountId int32
	if err := cache.GetRedisClient().Get(ctx, config.ActiveAccountId, &activeAccountId); err != nil {
		hlog.CtxErrorf(ctx, "failed to get active account id: %v", err)
		return err
	}

	if err := http.GetNotifier().Notify(ctx, mercari.SwitchAccountInfo{
		FromAccountID: activeAccountId,
		ToAccountID:   req.AccountID,
		Reason:        "manual switch account",
	}); err != nil {
		hlog.CtxErrorf(ctx, "failed to notify b4u: %v", err)
	}

	return nil
}
