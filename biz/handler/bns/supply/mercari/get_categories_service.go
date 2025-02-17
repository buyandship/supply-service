package mercari

import (
	"context"
	"encoding/json"
	"errors"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
	cache "github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/redis/go-redis/v9"
	"time"
)

func GetCategoriesService(ctx context.Context) (resp *mercari.GetItemCategoriesResp, err error) {
	cat, err := cache.GetHandler().Get(ctx, "mercari_categories")
	if errors.Is(err, redis.Nil) {
		// 1.2.3
		h := mercari.GetHandler()
		resp, err := h.GetItemCategories(ctx)
		if err != nil {
			return nil, err
		}
		r, err := json.Marshal(resp)
		if err != nil {
			hlog.CtxErrorf(ctx, "json marshal err: %v", err)
			return nil, bizErr.InternalError
		}
		if err := cache.GetHandler().Set(ctx, "mercari_categories", r, time.Hour); err != nil {
			hlog.CtxErrorf(ctx, "set mercari_categories err: %v", err)
		}
		return resp, nil
	}
	if err != nil {
		return nil, bizErr.InternalError
	}

	resp = &mercari.GetItemCategoriesResp{}
	if jsonStr, ok := cat.(string); ok {
		if err := json.Unmarshal([]byte(jsonStr), resp); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal err: %v", err)
			return nil, bizErr.InternalError
		}
		return resp, nil
	}

	return resp, bizErr.InternalError
}
