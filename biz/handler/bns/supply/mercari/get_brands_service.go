package mercari

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	cache "github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
)

func GetBrandsService(ctx context.Context) (resp *mercari.GetBrandsResp, err error) {
	cat, err := cache.GetHandler().Get(ctx, "mercari_brands")
	if errors.Is(err, redis.Nil) {
		// 1.2.3
		h := mercari.GetHandler()
		resp, err := h.GetBrands(ctx)
		if err != nil {
			return nil, err
		}

		go func() {
			r, err := json.Marshal(resp)
			if err != nil || string(r) == "null" || string(r) == "" {
				hlog.CtxErrorf(ctx, "json marshal err: %v", err)
				return
			}
			if err := cache.GetHandler().Set(ctx, "mercari_brands", r, time.Hour); err != nil {
				hlog.CtxErrorf(ctx, "set mercari_brands err: %v", err)
			}
		}()
		return resp, nil
	}
	if err != nil {
		return nil, bizErr.InternalError
	}

	resp = &mercari.GetBrandsResp{}
	if jsonStr, ok := cat.(string); ok {
		if err := json.Unmarshal([]byte(jsonStr), resp); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal err: %v", err)
			return nil, bizErr.InternalError
		}
		return resp, nil
	}

	return resp, bizErr.InternalError // cache

}
