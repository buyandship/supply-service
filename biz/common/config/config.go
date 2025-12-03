package config

const (
	DevMasterYahooAccountID = "chkyj_cp_evjr2p2v" // For product list
	DevYahoo02AccountID     = "chkyj_cp_by4d1vej" // For bidding
	DevYahoo03AccountID     = "chkyj_cp_c0rufa99" // For backuup bidding

	ProdMasterYahooAccountID = "juhyq37695" // For product list
	ProdYahoo02AccountID     = "wfwnd94596" // For bidding
	ProdYahoo03AccountID     = "vfkwo56810" // For backuup bidding
	ProdYahoo04AccountID     = "poxoe05998" // For backup bidding
)

// Cache keys
const (
	MercariBrandsKey     = "supply_srv:v1:mercari_brands"
	MercariCategoriesKey = "supply_srv:v1:mercari_categories"
	MercariAccountPrefix = "supply_srv:v1:account:%d"

	MercariRefreshTokenLock = "supply_srv:v1:mercari_refresh_token"
	MercariFailoverLock     = "supply_srv:v1:mercari_failover"

	ActiveAccountId     = "supply_srv:v1:active_account_id"
	TokenRedisKeyPrefix = "supply_srv:v1:token:%d"
)
