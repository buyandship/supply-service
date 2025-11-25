package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrasturcture/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/yahoo"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTransactionsService(ctx context.Context, req *supply.YahooGetTransactionsReq) ([]*model.Transaction, error) {
	hlog.CtxInfof(ctx, "GetTransactionService: %+v", req)
	yahooClient := yahoo.GetClient()
	// get transaction id list with order id
	// TODO: get bulk transaction
	tx, err := yahooClient.GetTransactions(ctx, req) // TODO: get yahoo account
	if err != nil {
		return nil, err
	}

	transactions := make([]*model.Transaction, 0)
	for _, tx := range tx.Transactions {
		transactions = append(transactions, &model.Transaction{
			TransactionID:   tx.TransactionID,
			YsRefID:         tx.YsRefID,
			AuctionID:       tx.AuctionID,
			CurrentPrice:    tx.CurrentPrice,
			TransactionType: tx.TransactionType,
			Status:          tx.Status,
			ReqPrice:        tx.ReqPrice,
		})
	}

	return transactions, nil
}
