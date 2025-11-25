package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
)

func GetCategoryTreeService(ctx context.Context, req *supply.YahooGetCategoryTreeReq) (*model.Category, error) {
	client := yahoo.GetClient()
	resp, err := client.GetCategoryTree(ctx, req)
	if err != nil {
		return nil, err
	}
	return &resp.ResultSet.Result, nil
}
