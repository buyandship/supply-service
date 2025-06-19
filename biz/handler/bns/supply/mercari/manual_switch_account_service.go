package mercari

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func ManualSwitchAccountService(ctx context.Context, req *supply.MercariManualSwitchAccountReq) error {

	if err := db.GetHandler().SwitchAccount(ctx, req.AccountID); err != nil {
		return err
	}

	if err := cache.GetHandler().Del(ctx, cache.ActiveAccountId); err != nil {
		hlog.CtxErrorf(ctx, "failed to del active account id: %v", err)
	}

	return nil
}
