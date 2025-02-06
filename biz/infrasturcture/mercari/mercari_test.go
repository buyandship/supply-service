package mercari

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"testing"
)

type ABC struct {
	A string
	B string
	C string
}

func TestMap(t *testing.T) {
	t.Run("test mapping", func(t *testing.T) {
		m := make(map[string]ABC)
		m["abc"] = ABC{
			A: "a",
			B: "b",
			C: "c",
		}
		c := m["abc"]
		c.A = "b"
		hlog.Infof("c: %v", c)
	})
}
