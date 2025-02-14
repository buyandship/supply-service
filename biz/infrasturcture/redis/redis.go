package redis

import (
	"context"
	"errors"
	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/google/uuid"
	"github.com/mennanov/limiters"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
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
	_, err := h.limiter.Limit(ctx)
	if errors.Is(err, limiters.ErrLimitExhausted) {
		return true
	}
	return false
}

func (h *H) HealthCheck() error {
	return h.redisClient.Ping(context.Background()).Err()
}

func (h *H) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return h.redisClient.Set(ctx, key, value, expiration).Err()
}

func (h *H) Get(ctx context.Context, key string) (interface{}, error) {
	return h.redisClient.Get(ctx, key).Result()
}
