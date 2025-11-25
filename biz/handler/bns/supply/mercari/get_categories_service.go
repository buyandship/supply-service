package mercari

import (
	"context"
	"time"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrasturcture/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetCategoriesService(ctx context.Context) (resp *mercari.GetItemCategoriesResp, err error) {
	resp = &mercari.GetItemCategoriesResp{}
	if err := cache.GetRedisClient().Get(ctx, config.MercariCategoriesKey, resp); err != nil {

		resp, err := mercari.GetHandler().GetItemCategories(ctx)
		if err != nil {
			return nil, err
		}
		if err := cache.GetRedisClient().Set(context.Background(), config.MercariCategoriesKey, resp, time.Hour); err != nil {
			hlog.Warnf("[goroutine] set mercari_categories err: %v", err)
		}
		return resp, nil
	}
	return resp, nil
}
