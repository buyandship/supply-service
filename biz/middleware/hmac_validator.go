package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
)

func HmacValidator() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := string(c.GetHeader("timestamp"))
		if ts == "" {
			hlog.Errorf("timestamp header is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, bizErr.UnauthorisedError)
		}
		src := string(c.GetHeader("hmac"))
		s := hmac.New(sha256.New, []byte("utmtuacz-hhme-fht4-gba6-3lg4vi8fzpu6"))
		s.Write([]byte(ts))
		target := hex.EncodeToString(s.Sum(nil))
		if src != target {
			hlog.CtxErrorf(ctx, "invalid hmac, [%s]!=[%s]", src, target)
			c.AbortWithStatusJSON(http.StatusUnauthorized, bizErr.UnauthorisedError)
		}
	}
}
