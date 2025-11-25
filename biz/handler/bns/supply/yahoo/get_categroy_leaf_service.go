package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetCategoryLeafService(ctx context.Context, req *supply.YahooGetCategoryLeafReq) (*yahoo.CategoryLeafResponse, error) {
	client := yahoo.GetClient()
	resp, err := client.GetCategoryLeaf(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
