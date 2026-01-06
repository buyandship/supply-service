package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func DeleteMyWonListService(ctx context.Context, req *supply.YahooDeleteMyWonListReq) error {
	_, err := yahoo.GetClient().DeleteMyWonList(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
