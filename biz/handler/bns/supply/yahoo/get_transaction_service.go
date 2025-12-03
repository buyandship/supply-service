package yahoo

import (
	"context"
	"strings"

	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func GetTransactionsService(ctx context.Context, req *supply.YahooGetTransactionsReq) ([]*yahoo.TransactionResult, error) {
	hlog.CtxInfof(ctx, "GetTransactionService: %+v", req)
	yahooClient := yahoo.GetClient()
	// get transaction id list with order id
	// TODO: get bulk transaction
	tx, err := yahooClient.GetTransactions(ctx, req) // TODO: get yahoo account
	if err != nil {
		return nil, err
	}

	transactions := make([]*yahoo.TransactionResult, 0)
	for _, tx := range tx.Transactions {
		apiEndpoint := strings.Split(tx.APIEndpoint, "/")
		eventType := apiEndpoint[len(apiEndpoint)-1]
		if eventType != "placeBidPreview" && eventType != "placeBid" {
			continue
		}
		transactions = append(transactions, &yahoo.TransactionResult{
			TransactionID:   tx.TransactionID,
			YsRefID:         tx.YsRefID,
			AuctionID:       tx.AuctionID,
			CurrentPrice:    tx.CurrentPrice,
			TransactionType: tx.TransactionType,
			Status:          tx.Status,
			ReqPrice:        tx.ReqPrice,
			EventType:       eventType,
			CreatedAt:       tx.CreatedAt,
			UpdatedAt:       tx.UpdatedAt,
			Detail:          tx.ResponseData.ResultSet.Result,
		})
	}

	return transactions, nil
}
