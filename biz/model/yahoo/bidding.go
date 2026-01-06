package yahoo

type BiddingMetaInfo struct {
	WonPrice int64 `json:"won_price"`
}

type BiddingResult struct {
	MetaInfo     BiddingMetaInfo `json:"ss_meta"`
	Status       string          `json:"status"`
	OrderNumber  string          `json:"order_number"`
	ErrorMessage string          `json:"error_message"`
}
