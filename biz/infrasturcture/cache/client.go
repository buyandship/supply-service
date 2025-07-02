package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (result any, err error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) (err error)
	Del(ctx context.Context, key string) (err error)
	TryLock(ctx context.Context, key string) (success bool, err error)
	Unlock(ctx context.Context, key string) (err error)
	Limit(ctx context.Context) bool
	HealthCheck() error
}
