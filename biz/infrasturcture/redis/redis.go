package redis

import (
	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"sync"
)

var (
	once    sync.Once
	Handler *H
)

type H struct {
	cli     *goredislib.Client
	redSync *redsync.Redsync
}

func GetHandler() *H {
	once.Do(func() {
		cli := goredislib.NewClient(&goredislib.Options{
			Addr: config.GlobalServerConfig.Redis.Address,
		})
		pool := goredis.NewPool(cli) // or, pool := redigo.NewPool(...)
		redSync := redsync.New(pool)
		Handler = &H{
			cli:     cli,
			redSync: redSync,
		}
	})
	return Handler
}

func (h *H) HealthCheck() error {
	return nil
}
