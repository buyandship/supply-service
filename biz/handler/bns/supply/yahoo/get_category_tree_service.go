package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/yahoo"
)

func GetCategoryTreeService(ctx context.Context, req *supply.YahooGetCategoryTreeReq) (*model.Category, error) {
	client := yahoo.GetClient()
	resp, err := client.MockGetCategoryTree(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp.ResultSet.Result, nil
}
