package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetSellingListService(ctx context.Context, req *supply.YahooGetSellingListReq) (*yahoo.SellingListResponse, error) {
	client := yahoo.GetClient()
	resp, err := client.GetSellingList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
