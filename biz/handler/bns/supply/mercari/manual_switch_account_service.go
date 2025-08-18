package mercari

import (
	"context"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/http"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/buyandship/supply-svr/biz/model/mercari"
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

	if err := cache.GetHandler().Set(ctx, config.ActiveAccountId, req.AccountID, time.Hour); err != nil {
		hlog.CtxErrorf(ctx, "failed to set active account id: %v", err)
	}

	var activeAccountId int32
	if err := cache.GetHandler().Get(ctx, config.ActiveAccountId, &activeAccountId); err != nil {
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
