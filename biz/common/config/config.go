package config

const (
	MasterYahooAccountID = "chkyj_cp_evjr2p2v" // For product list
	Yahoo02AccountID     = "chkyj_cp_by4d1vej" // For bidding
	Yahoo03AccountID     = "chkyj_cp_c0rufa99" // For backuup bidding
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
