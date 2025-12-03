package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetAccount(ctx context.Context, accId int32) (*mercari.Account, error) {

	acc := &mercari.Account{}
	if err := cache.GetRedisClient().Get(ctx, fmt.Sprintf(config.MercariAccountPrefix, accId), acc); err != nil {
		// degrade to load from
		acc, err := db.GetHandler().GetAccount(ctx, accId)
		if err != nil {
			return nil, err
		}
		go func() {
			if err := cache.GetRedisClient().Set(context.Background(), fmt.Sprintf(config.MercariAccountPrefix, accId), acc, time.Hour); err != nil {
				hlog.Warnf("[goroutine] redis set buyer error: %+v", err)
			}
		}()
		return acc, err
	}
	return acc, nil
}
