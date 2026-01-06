package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetMyWonListService(ctx context.Context, req *supply.YahooGetMyWonListReq) (*yahoo.MyWonListResponse, error) {
	return yahoo.GetClient().GetMyWonList(ctx, req)
}
