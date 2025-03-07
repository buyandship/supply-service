package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/infrasturcture/redis"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/gorm"
	"time"
)

const (
	DefaultBuyerId = 1
)

func GetBuyer(ctx context.Context, buyerId int32) (*mercari.Account, error) {
	if buyerId == 0 {
		buyerId = DefaultBuyerId
	}
	buyer, err := redis.GetHandler().Get(ctx, fmt.Sprintf("buyer:%d", buyerId))
	if err != nil {
		// degrade to load from
		acc, err := db.GetHandler().GetAccount(ctx, buyerId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, bizErr.InvalidBuyerError
			}
			return nil, bizErr.InternalError
		}
		jsonAcc, err := json.Marshal(acc)
		if err != nil {
			hlog.CtxErrorf(ctx, "json marshal error: %+v", err)
			return acc, nil
		}
		if err := redis.GetHandler().Set(ctx, fmt.Sprintf("buyer:%d", buyerId), jsonAcc, time.Hour); err != nil {
			hlog.CtxErrorf(ctx, "redis set buyer error: %+v", err)
			return acc, nil
		}
		return acc, err
	}
	if bBuyer, ok := buyer.(string); ok {
		acc := &mercari.Account{}
		if err := json.Unmarshal([]byte(bBuyer), acc); err != nil {
			hlog.CtxErrorf(ctx, "json unmarshal error: %+v", err)
			return nil, bizErr.InternalError
		}
		hlog.CtxInfof(ctx, "buyer is %+v", acc)
		return acc, nil
	}
	return nil, bizErr.InternalError
}
