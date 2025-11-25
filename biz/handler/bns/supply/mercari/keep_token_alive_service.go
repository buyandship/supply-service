package mercari

import (
	"context"
	"time"

	"github.com/buyandship/supply-service/biz/infrasturcture/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func KeepTokenAliveService(ctx context.Context) error {

	// get account list from db
	accounts, err := db.GetHandler().GetAccountList(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		token, err := db.GetHandler().GetToken(ctx, int32(account.ID))
		if err != nil {
			continue
		}

		if token.CreatedAt.Before(time.Now().Add(-85 * time.Hour * 24)) {
			hlog.CtxInfof(ctx, "account: %d, token expired: %v", account.ID, token)
			// refresh token
			// if err := mercari.GetHandler().RefreshToken(ctx, token); err != nil {
			//	hlog.CtxErrorf(ctx, "refresh token error: %v", err)
			//	continue
			// }
			// hlog.CtxInfof(ctx, "refresh token success: %",)
		}

	}
	return nil
}
