package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func BuyoutService(ctx context.Context, req *supply.YahooBuyoutReq) (string, error) {
	return "ok", nil
}
