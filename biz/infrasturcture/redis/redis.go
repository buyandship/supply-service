package redis

import (
	"context"
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
	TokenRedisKey = "mercari_token"
	LockKeyPrefix = "lock:"
	LockTTL       = 2 * time.Second
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

	err = h.redisClient.Set(ctx, key, value, expiration).Err()
	return err
}

func (h *H) Get(ctx context.Context, key string) (result any, err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Get", key)
	defer trace.EndSpan(span, nil)

	result, err = h.redisClient.Get(ctx, key).Result()
	return
}

func (h *H) Del(ctx context.Context, key string) (err error) {
	ctx, span := trace.StartRedisOperation(ctx, "Del", key)
	defer trace.EndSpan(span, nil)

	err = h.redisClient.Del(ctx, key).Err()
	return
}

func (h *H) TryLock(ctx context.Context, key string) (success bool, err error) {
	ctx, span := trace.StartRedisOperation(ctx, "TryLock", key)
	defer trace.EndSpan(span, nil)

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
