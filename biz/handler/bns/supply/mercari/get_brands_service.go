package mercari

import (
	"context"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/infrasturcture/cache"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetBrandsService(ctx context.Context) (resp *mercari.GetBrandsResp, err error) {
	resp = &mercari.GetBrandsResp{}
	if err := cache.GetHandler().Get(ctx, config.MercariBrandsKey, resp); err != nil {
		resp, err := mercari.GetHandler().GetBrands(ctx)
		if err != nil {
			return nil, err
		}

		go func() {
			if err := cache.GetHandler().Set(context.Background(), config.MercariBrandsKey, resp, time.Hour); err != nil {
				hlog.Warnf("[goroutine] set mercari_brands err: %v", err)
			}
		}()

		return resp, nil
	}
	return resp, nil
}
