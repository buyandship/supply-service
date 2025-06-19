package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

const (
	DefaultBuyerId = 1
)

func GetAccount(ctx context.Context, buyerId int32) (*mercari.Account, error) {

	acc := &mercari.Account{}
	if buyerId == 0 {
		buyerId = DefaultBuyerId
	}
	if err := cache.GetHandler().Get(ctx, fmt.Sprintf(config.MercariAccountPrefix, buyerId), acc); err != nil {
		// degrade to load from
		acc, err := db.GetHandler().GetAccount(ctx, buyerId)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := cache.GetHandler().Set(context.Background(), fmt.Sprintf(config.MercariAccountPrefix, buyerId), acc, time.Hour); err != nil {
				hlog.CtxWarnf(ctx, "redis set buyer error: %+v", err)
			}
		}()
		return acc, err
	}
	return acc, nil
}
