package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"net/http"
)

func HmacValidator() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ts := string(c.GetHeader("timestamp"))
		if ts == "" {
			hlog.CtxErrorf(ctx, "timestamp header is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, bizErr.UnauthorisedError)
			return
		}
		src := string(c.GetHeader("hmac"))
		s := hmac.New(sha256.New, []byte(config.GlobalServerConfig.HmacSecret))
		s.Write([]byte(ts))
		target := hex.EncodeToString(s.Sum(nil))
		if src != target {
			hlog.CtxErrorf(ctx, "invalid hmac, [%s]!=[%s]", src, target)
			c.AbortWithStatusJSON(http.StatusUnauthorized, bizErr.UnauthorisedError)
			return
		}
	}
}
