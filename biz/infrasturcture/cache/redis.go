package cache

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/common/trace"
	"github.com/google/uuid"
	"github.com/mennanov/limiters"
	"github.com/redis/go-redis/v9"
)

const (
	TokenRedisKeyPrefix = "supplysrv:v1:token:%d"
	LockKeyPrefix       = "lock:"
	LockTTL             = 2 * time.Second
)

var (
	once    sync.Once
	Handler *H
)

type H struct {
	limiter     *limiters.SlidingWindow
	redisClient *redis.Client
}

func GetHandler() *H {
	once.Do(func() {
		redisClient := redis.NewClient(&redis.Options{
			Addr: config.GlobalServerConfig.Redis.Address,
		})
		backend := limiters.NewSlidingWindowRedis(redisClient, uuid.NewString())
		limiter := limiters.NewSlidingWindow(40, time.Second, backend, limiters.NewSystemClock(), 3e-9)

		Handler = &H{
			redisClient: redisClient,
			limiter:     limiter,
		}
	})
	return Handler
}

func (h *H) Limit(ctx context.Context) bool {
	ctx, span := trace.StartRedisOperation(ctx, "Limit", "")
	var err error
	defer trace.EndSpan(span, err)

	_, err = h.limiter.Limit(ctx)
	if errors.Is(err, limiters.ErrLimitExhausted) {
		return true
	}
	return false
}

func (h *H) HealthCheck() error {
	return h.redisClient.Ping(context.Background()).Err()
}

func (h *H) Set(ctx context.Context, key string, value any, expiration time.Duration) (err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Set", key)
	defer trace.EndSpan(span, err)

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return h.redisClient.Set(ctx, key, jsonValue, expiration).Err()
}

func (h *H) Get(ctx context.Context, key string, value any) (err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Get", key)
	defer trace.EndSpan(span, err)

	result, err := h.redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(result), value)
}

func (h *H) Del(ctx context.Context, key string) (err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Del", key)
	defer trace.EndSpan(span, err)

	return h.redisClient.Del(ctx, key).Err()
}

func (h *H) TryLock(ctx context.Context, key string) (success bool, err error) {
	ctx, span := trace.StartRedisOperation(ctx, "TryLock", key)
	defer trace.EndSpan(span, err)

	lockKey := LockKeyPrefix + key
	success, err = h.redisClient.SetNX(ctx, lockKey, 1, LockTTL).Result()
	if err != nil {
		return false, err
	}
	return
}

func (h *H) Unlock(ctx context.Context, key string) (err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Unlock", key)
	defer trace.EndSpan(span, err)

	lockKey := LockKeyPrefix + key
	err = h.redisClient.Del(ctx, lockKey).Err()
	return
}
