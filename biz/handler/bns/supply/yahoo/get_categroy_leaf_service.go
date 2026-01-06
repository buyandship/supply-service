package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetCategoryLeafService(ctx context.Context, req *supply.YahooGetCategoryLeafReq) (*yahoo.CategoryLeafResponse, error) {
	return yahoo.GetClient().GetCategoryLeaf(ctx, req)
}
