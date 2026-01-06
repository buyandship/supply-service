package yahoo

import (
	"context"

	"github.com/buyandship/supply-service/biz/infrastructure/http"
)

func WebHookService(batchNumber string, body []byte, retryAttempt int) error {

	return http.GetNotifier().NotifyBiddingStatus(context.Background(), batchNumber, body)
}
