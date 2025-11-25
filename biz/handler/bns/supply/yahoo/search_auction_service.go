package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
)

func SearchAuctionService(ctx context.Context, req *supply.YahooSearchAuctionsReq) (*yahoo.SearchAuctionsResponse, error) {
	client := yahoo.GetClient()
	searchAuctionsResp, err := client.SearchAuctions(ctx, &yahoo.SearchAuctionsRequest{
		Query:                req.Keyword,
		Type:                 req.Type,
		Category:             int(req.Category),
		ExceptCategory:       req.ExpectCategory,
		Page:                 int(req.Page),
		Sort:                 req.Sort,
		Order:                req.Order,
		Store:                int(req.Store),
		AucMinPrice:          int(req.Aucminprice),
		AucMaxPrice:          int(req.Aucmaxprice),
		AucMinBidorbuyPrice:  int(req.AucminBidorbuyPrice),
		AucMaxBidorbuyPrice:  int(req.AucmaxBidorbuyPrice),
		LocCd:                int(req.LocCd),
		EasyPayment:          int(req.Easypayment),
		New:                  int(req.New),
		FreeShipping:         int(req.Freeshipping),
		WrappingIcon:         int(req.Wrappingicon),
		BuyNow:               int(req.Buynow),
		Thumbnail:            int(req.Thumbnail),
		Attn:                 int(req.Attn),
		Point:                int(req.Point),
		ItemStatus:           int(req.ItemStatus),
		Adf:                  int(req.Adf),
		SellerAucUserID:      req.SellerAucUserID,
		F:                    req.F,
		Ngram:                int(req.Ngram),
		Fixed:                int(req.Fixed),
		MinCharity:           int(req.MinCharity),
		MaxCharity:           int(req.MaxCharity),
		MinAffiliate:         int(req.MinAffiliate),
		MaxAffiliate:         int(req.MaxAffiliate),
		Timebuf:              int(req.Timebuf),
		Ranking:              req.Ranking,
		BlackSellerAucUserID: req.BlackSellerAucUserID,
		Featured:             req.Featured,
		Sort2:                req.Sort2,
		Order2:               req.Order2,
		MinStart:             int64(req.MinStart),
		MaxStart:             int64(req.MaxStart),
		ExceptShoppingItem:   true,
	})
	if err != nil {
		return nil, err
	}
	return searchAuctionsResp, nil
}
