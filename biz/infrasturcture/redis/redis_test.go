package redis

import (
	"context"
	"github.com/buyandship/supply-svr/biz/common/config"
	"testing"
	"time"
)

func TestLimit(t *testing.T) {
	t.Run("test", func(t *testing.T) {
		config.LoadTestConfig()

		for i := 0; i < 10; i++ {
			err := GetHandler().Limit(context.Background())
			t.Logf("%d: err:%v", i, err)
		}
		time.Sleep(2 * time.Second)
		for i := 0; i < 10; i++ {
			err := GetHandler().Limit(context.Background())
			t.Logf("%d: err:%v", i, err)
		}
	})
}
