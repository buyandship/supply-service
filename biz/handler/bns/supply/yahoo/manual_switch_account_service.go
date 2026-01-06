package yahoo

import (
	"context"
	"time"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/infrastructure/mq"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func ManualSwitchAccountService(ctx context.Context, req *supply.YahooManualSwitchAccountReq) error {
	if err := mq.SendMessage(mq.Message{
		Topic:      config.OneHourBufferMessageQueue,
		RoutingKey: config.OneHourBufferMessageRoutingKey,
		Msg:        "test message queue",
	}); err != nil {
		hlog.CtxErrorf(ctx, "failed to send message: %v", err)
	}
	// switch account
	if err := db.GetHandler().SwitchYahooAccount(ctx, req.AccountID); err != nil {
		hlog.CtxErrorf(ctx, "failed to switch account: %v", err)
		return err
	}

	// set redis
	if err := cache.GetRedisClient().Set(ctx, config.YahooActiveAccountId, req.AccountID, time.Hour); err != nil {
		hlog.CtxErrorf(ctx, "failed to set active account id: %v", err)
		return err
	}

	return nil
}
