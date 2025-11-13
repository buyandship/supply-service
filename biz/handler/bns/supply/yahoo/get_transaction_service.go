package yahoo

import (
	"context"

	"github.com/buyandship/supply-svr/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/yahoo"
)

func GetTransactionService(ctx context.Context, req *supply.YahooGetTransactionReq) (*model.Transaction, error) {
	yahooClient := yahoo.GetClient()

	// get transaction id list with order id

	// TODO: get bulk transaction
	tx, err := yahooClient.MockGetTransaction(ctx, req, "account123") // TODO: get yahoo account
	if err != nil {
		return nil, err
	}

	resp := &model.Transaction{
		TransactionID: tx.TransactionID,
	}

	return resp, nil
}
