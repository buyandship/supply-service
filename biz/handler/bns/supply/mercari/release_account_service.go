package mercari

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func ReleaseAccountService(ctx context.Context, accountID string) error {
	hlog.CtxInfof(ctx, "ReleaseAccountService is called: %s", accountID)

	if err := db.GetHandler().ReleaseAccount(ctx, accountID); err != nil {
		return err
	}

	return nil
}
