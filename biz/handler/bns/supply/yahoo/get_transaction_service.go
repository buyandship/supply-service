package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
)

func GetTransactionService(ctx context.Context, req *supply.YahooGetTransactionReq) (*supply.YahooTransaction, error) {
	yahooClient := yahoo.GetClient()
	tx, err := yahooClient.MockGetTransaction(ctx, req, "account123") // TODO: get yahoo account
	if err != nil {
		return nil, err
	}

	resp := &supply.YahooTransaction{
		TransactionID:   tx.TransactionID,
		YsRefID:         tx.YsRefID,
		AuctionID:       tx.AuctionID,
		CurrentPrice:    tx.CurrentPrice,
		TransactionType: tx.TransactionType,
		Status:          tx.Status,
		ReqPrice:        tx.ReqPrice,
	}

	return resp, nil
}
