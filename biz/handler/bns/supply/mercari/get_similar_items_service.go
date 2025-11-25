package mercari

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/mercari"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func GetSimilarItemsService(ctx context.Context, req *supply.MercariGetSimilarItemsReq) (*mercari.GetSimilarItemsResponse, error) {
	h := mercari.GetHandler()

	return h.GetSimilarItems(ctx, req)
}
