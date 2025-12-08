// idl/hello.thrift
namespace go bns.supply

struct MercariGetItemReq {
    1: string item_id (api.json="item_id");
}

struct MercariGetSellerReq {
    1: string seller_id (api.json="seller_id");
}

struct MercariPostOrderReq {
    1: string item_id (api.json="item_id")
    2: i64 ref_price (api.json="ref_price")
    3: string ref_currency (api.json="ref_currency")
    4: string checksum (api.json="checksum")
    5: string ref_id (api.json="ref_id")
    6: string delivery_id (api.json = "delivery_id")
}

struct MercariPostMessageReq {
    1: string trx_id (api.json="trx_id");
    2: string msg (api.json="msg");
}

struct MercariPostMessageResp {
    1: string trx_id
    2: string body
    3: string id
    4: i32 account_id
}

struct MercariRegisterAccountReq {
    1: string email (api.json="email");
    2: string family_name (api.json="family_name");
    3: string first_name (api.json="first_name");
    4: string family_name_kana (api.json="family_name_kana");
    5: string first_name_kana (api.json="first_name_kana");
    6: string telephone (api.json="telephone");
    7: string zip_code1 (api.json="zip_code1");
    8: string zip_code2 (api.json="zip_code2");
    9: string prefecture (api.json="prefecture");
    10: string city (api.json="city");
    11: string address1 (api.json="address1");
    12: string address2 (api.json="address2");
    13: string status (api.json="status");
    14: i32 priority (api.json="priority");
    15: string banned_at (api.json="banned_at");
    16: string active_at (api.json="active_at");
}

struct Account {
    1: i32 id
    2: string email
    3: string family_name
    4: string first_name
    5: string family_name_kana
    6: string first_name_kana
    7: string telephone
    8: string zip_code1
    9: string zip_code2
    10: string prefecture
    11: string city
    12: string address1
    13: string address2
    14: string status
    15: i32 priority
    16: optional string banned_at
    17: optional string active_at
}

struct MercariRegisterAccountResp {
    1: Account account
}

struct MercariLoginCallBackReq {
    1: string code (api.query="code")
    2: string scope (api.query="scope")
    3: string state (api.query="state")
    4: string redirectUrl
}

struct MercariGetTokenResp {
    1: string token
}

struct MercariGetTransactionByItemIdReq {
    1: string item_id (api.json="item_id");
}

struct MercariPostOrderResp {
    1: string trx_id
    2: i64 coupon_id
    3: i64 price
    4: i64 paid_price
    5: i64 buyer_shipping_fee
    6: string item_id
    7: string checksum
    8: i32 account_id
}

struct MercariPostTransactionReviewReq {
    1: string trx_id
    2: string fame
    3: string review
    4: i32 account_id
}

struct MercariGetTodoListReq {
    1: i32 limit
    2: string page_token
}

struct MercariSearchItemsReq {
    1: optional string keyword
    2: optional string exclude_keyword
    3: optional string category_id
    4: optional string brand_id
    5: optional string seller_id
    6: optional string shop_id
    7: optional string size_id
    8: optional string color_id
    9: optional i32 price_min
    10: optional i32 price_max
    11: optional i32 created_before_date
    12: optional i32 created_after_date
    13: optional string item_condition_id
    14: optional string shipping_payer_id
    15: optional string status
    16: optional string marketplace
    17: optional string sort
    18: optional string order
    19: optional i32 page
    20: optional i32 limit
    21: optional bool item_authentication
    22: optional bool time_sale
    23: optional bool with_offer_price_promotion
}


struct MercariGetAccountResp {
    1: list<Account> accounts
}

struct MercariManualSwitchAccountReq {
    1: i32 account_id (api.json="account_id");
}

struct MercariFetchItemsReq {
    1: list<string> item_ids (api.json="item_ids");
}

struct MercariGetSimilarItemsReq {
    1: string item_id (api.json="item_id");
}


struct YahooPlaceBidReq {
    1: string ys_ref_id (api.json="ys_ref_id")
    2: string transaction_type (api.json="transaction_type")
    3: string auction_id (api.json="auction_id")
    4: i32 price (api.json="price")
    5: i32 quantity (api.json="quantity")
    6: bool partial (api.json="partial")
}

/*
struct YahooPlaceBidResp {
    1: string status (api.json="status" go.tag="example:\"Success\"");
    2: string bid_id (api.json="bid_id" go.tag="example:\"bid_12345\"");
    3: string auction_id (api.json="auction_id" go.tag="example:\"x12345\"");
    4: i32 price (api.json="price" go.tag="example:\"1000\"");
    5: i32 quantity (api.json="quantity" go.tag="example:\"1\"");
    6: i32 total_price (api.json="total_price" go.tag="example:\"1100\"");
    7: string bid_time (api.json="bid_time" go.tag="example:\"2025-10-22T12:00:00Z\"");
}
*/

struct YahooGetTransactionReq {
    1: string transaction_id (api.query="transaction_id" go.tag="validate:\"required\"");
}


/*
struct YahooTransaction {
    1: string transaction_id (api.json="transaction_id" go.tag="example:\"txn_abc123\"");
    2: string ys_ref_id (api.json="ys_ref_id" go.tag="example:\"YS-REF-001\"");
    3: string auction_id (api.json="auction_id" go.tag="example:\"x12345\"");
    4: i64 current_price (api.json="current_price" go.tag="example:\"1000\"");
    5: string transaction_type (api.json="transaction_type" go.tag="example:\"BID\"");
    6: string status (api.json="status" go.tag="example:\"completed\"");
    7: i64 req_price (api.json="req_price" go.tag="example:\"1000\"");
}
*/

struct YahooGetTransactionsReq {
   1: string transaction_id (api.query="transaction_id");
   2: string ys_ref_id (api.query="ys_ref_id");
   3: string auction_id (api.query="auction_id");
}

/*
struct YahooGetTransactionsResp {
    1: list<YahooTransaction> transactions (api.json="transactions");
}
*/

struct YahooGetAuctionItemReq {
    1: string auction_id (api.query="auction_id" go.tag="validate:\"required\"");
}

struct Seller {
    1: string id (api.json="id" go.tag="example:\"seller123\"");
    2: double rating (api.json="rating" go.tag="example:\"98.5\"");
    3: bool is_suspended (api.json="is_suspended" go.tag="example:\"false\"");
    4: bool is_deleted (api.json="is_deleted" go.tag="example:\"false\"");
}

struct ShoppingItem {
    1: bool is_option_enabled (api.json="is_option_enabled" go.tag="example:\"true\"");
}

struct YahooGetAuctionItemResp {
    1: string auction_id (api.json="auction_id" go.tag="example:\"x12345\"");
    2: string title (api.json="title" go.tag="example:\"Sample Item Title\"");
    3: string description (api.json="description" go.tag="example:\"Item description...\"");
    4: i64 current_price (api.json="current_price" go.tag="example:\"1000\"");
    5: i64 start_price (api.json="start_price" go.tag="example:\"500\"");
    6: i32 bids (api.json="bids" go.tag="example:\"5\"");
    7: string item_status (api.json="item_status" go.tag="example:\"open\"");
    8: string end_time (api.json="end_time" go.tag="example:\"2025-10-30T23:59:59Z\"");
    9: string start_time (api.json="start_time" go.tag="example:\"2025-10-22T00:00:00Z\"");
    10: Seller seller (api.json="seller");
    11: string image (api.json="image" go.tag="example:\"https://example.com/image.jpg\"");
    12: i32 quantity (api.json="quantity" go.tag="example:\"1\"");
    13: string shopping_item_code (api.json="shopping_item_code" go.tag="example:\"1234567890\"");
    14: ShoppingItem shopping_item (api.json="shopping_item");
    15: string bidorbuy (api.json="bidorbuy" go.tag="example:\"1000\"");
}

struct YahooGetAuctionItemAuthReq {
    1: string auction_id (api.query="auction_id" go.tag="validate:\"required\"");
}

struct BidStatus {
    1: bool has_bid (api.json="has_bid" go.tag="example:\"true\"");
    2: i64 my_highest_bid (api.json="my_highest_bid" go.tag="example:\"950\"");
    3: bool is_winning (api.json="is_winning" go.tag="example:\"false\"");
}

struct YahooGetAuctionItemAuthResp {
    1: YahooGetAuctionItemResp auction_item (api.json="auction_item")
    2: bool is_watching (api.json="is_watching" go.tag="example:\"true\"");
    3: BidStatus bid_status (api.json="bid_status");
}



struct YahooGetCategoryTreeReq {
    1: i32 category (api.query="category")
    2: i32 adf (api.query="adf")
    3: i32 is_fnavi_only (api.query="is_fnavi_only")
}

struct YahooSearchAuctionsReq {
    1: string keyword (api.query="keyword")
    2: string type (api.query="type")
    3: i32 category (api.query="category")
    4: string expect_category (api.query="expect_category")
    5: i32 page (api.query="page")
    6: string sort (api.query="sort")
    7: string order (api.query="order")
    8: i32 store (api.query="store")
    9: i32 aucminprice (api.query="aucminprice")
    10: i32 aucmaxprice (api.query="aucmaxprice")
    11: i32 aucmin_bidorbuy_price (api.query="aucmin_bidorbuy_price")
    12: i32 aucmax_bidorbuy_price (api.query="aucmax_bidorbuy_price")
    13: i32 loc_cd (api.query="loc_cd")
    14: i32 easypayment (api.query="easypayment")
    15: i32 new (api.query="new")
    16: i32 freeshipping (api.query="freeshipping")
    17: i32 wrappingicon (api.query="wrappingicon")
    18: i32 buynow (api.query="buynow")
    19: i32 thumbnail (api.query="thumbnail")
    20: i32 attn (api.query="attn")
    21: i32 point (api.query="point")
    22: i32 item_status (api.query="item_status")
    23: i32 adf (api.query="adf")
    24: string seller_auc_user_id (api.query="seller_auc_user_id")
    25: string f (api.query="f")
    26: i32 ngram (api.query="ngram")
    27: i32 fixed (api.query="fixed")
    28: i32 min_charity (api.query="min_charity")
    29: i32 max_charity (api.query="max_charity")
    30: i32 min_affiliate (api.query="min_affiliate")
    31: i32 max_affiliate (api.query="max_affiliate")
    32: i32 timebuf (api.query="timebuf")
    33: string ranking (api.query="ranking")
    34: string black_seller_auc_user_id (api.query="black_seller_auc_user_id")
    35: string featured (api.query="featured")
    36: string sort2 (api.query="sort2")
    37: string order2 (api.query="order2")
    38: i32 min_start (api.query="min_start")
    39: i32 max_start (api.query="max_start")
    40: bool except_shoppingitem (api.query="except_shoppingitem")
} 

struct YahooGetCategoryLeafReq {
    1: optional string yahoo_account_id (api.query="yahoo_account_id")
    2: optional string ys_ref_id (api.query="ys_ref_id")
    3: i32 category (api.query="category")
    4: optional string except_category (api.query="except_category")
    5: optional i32 featured (api.query="featured")
    6: optional i32 page (api.query="page")
    7: optional string sort (api.query="sort")
    8: optional string order (api.query="order")
    9: optional i32 store (api.query="store")
    10: optional i32 aucminprice (api.query="aucminprice")
    11: optional i32 aucmaxprice (api.query="aucmaxprice")
    12: optional i32 aucmin_bidorbuy_price (api.query="aucmin_bidorbuy_price")
    13: optional i32 aucmax_bidorbuy_price (api.query="aucmax_bidorbuy_price")
    14: optional i32 easypayment (api.query="easypayment")
    15: optional i32 new (api.query="new")
    16: optional i32 freeshipping (api.query="freeshipping")
    17: optional i32 wrappingicon (api.query="wrappingicon")
    18: optional i32 buynow (api.query="buynow")
    19: optional i32 thumbnail (api.query="thumbnail")
    20: optional i32 attn (api.query="attn")
    21: optional i32 point (api.query="point")
    22: optional string item_status (api.query="item_status")
    23: optional i32 adf (api.query="adf")
    24: optional i32 min_charity (api.query="min_charity")
    25: optional i32 max_charity (api.query="max_charity")
    26: optional i32 min_affiliate (api.query="min_affiliate")
    27: optional i32 max_affiliate (api.query="max_affiliate")
    28: optional i32 timebuf (api.query="timebuf")
    29: optional i32 ranking (api.query="ranking")
    30: optional string seller_auc_user_id (api.query="seller_auc_user_id")
    31: optional string black_seller_auc_user_id (api.query="black_seller_auc_user_id")
    32: optional string sort2 (api.query="sort2")
    33: optional string order2 (api.query="order2")
    34: optional string loc_cd (api.query="loc_cd")
    35: optional i32 fixed (api.query="fixed")
    36: optional i64 max_start (api.query="max_start")
    37: optional i64 min_start (api.query="min_start")
    38: optional i32 except_shoppingitem (api.query="except_shoppingitem")
    39: optional string callback (api.query="callback")
}

struct YahooGetMyWonListReq {
    1: optional string ys_ref_id (api.query="ys_ref_id")
    2: optional i32 start (api.query="start")
    3: optional string contact_progress (api.query="contact_progress")
    4: optional string auction_id (api.query="auction_id")
}

struct YahooGetSellingListReq {
    1: string sellerAucUserId (api.query="sellerAucUserId")
    2: optional string ys_ref_id (api.query="ys_ref_id")
    3: optional i32 start (api.query="start")
    4: optional string status (api.query="status")
}

service SupplyService {
    string MercariGetItemService(1: MercariGetItemReq req) (api.get="/v1/supplysrv/internal/mercari/item");
    string MercariGetCategoriesService() (api.get="/v1/supplysrv/internal/mercari/categories")
    string MercariGetSellerService(1: MercariGetSellerReq req) (api.get="/v1/supplysrv/internal/mercari/seller")
    string MercariLoginCallBackService(1: MercariLoginCallBackReq req) (api.get="/xb/login_callback")
    string MercariGetTransactionByItemIdService(1: MercariGetTransactionByItemIdReq req) (api.get="/v1/supplysrv/internal/mercari/tx")
    string MercariPostTransactionReviewService(1: MercariPostTransactionReviewReq req) (api.post="/v1/supplysrv/internal/mercari/review")
    string MercariGetTodoListService(1: MercariPostTransactionReviewReq req) (api.get="/v1/supplysrv/internal/mercari/todo")
    string MercariSearchItemsService(1: MercariSearchItemsReq req) (api.get="/v1/supplysrv/internal/mercari/search")
    string MercariGetBrandsService() (api.get="/v1/supplysrv/internal/mercari/brands")
    string MercariManualSwitchAccountService(1: MercariManualSwitchAccountReq req) (api.post="/v1/supplysrv/public/mercari/switch_account")
    string MercariReleaseAccountService(1: string account_id) (api.post="/v1/supplysrv/public/mercari/release_account")
    string KeepTokenAliveService() (api.post="/v1/supplysrv/internal/mercari/keep_token_alive")

    MercariRegisterAccountResp MercariRegisterAccountService(1: MercariRegisterAccountReq req) (api.post="/v1/supplysrv/public/mercari/register")
    MercariPostOrderResp MercariPostOrderService(1: MercariPostOrderReq req) (api.post="/v1/supplysrv/internal/mercari/order")
    MercariPostMessageResp MercariPostMessageService(1: MercariPostMessageReq req) (api.post="/v1/supplysrv/internal/mercari/message")
    MercariGetTokenResp MercariGetTokenService() (api.get="/v1/supplysrv/internal/mercari/token")
    MercariGetAccountResp MercariGetAccountService() (api.get="/v1/supplysrv/internal/mercari/account/list")
    string MercariFetchItemsService(1: MercariFetchItemsReq req) (api.get="/v1/supplysrv/internal/mercari/fetch_items")
    string MercariGetSimilarItemsService(1: MercariGetSimilarItemsReq req) (api.get="/v1/supplysrv/internal/mercari/similar_items")


    // YahooPlaceBidPreviewResp YahooPlaceBidPreviewService(1: YahooPlaceBidPreviewReq req) (api.post="/v1/supplysrv/internal/yahoo/placeBidPreview")
    string YahooPlaceBidService(1: YahooPlaceBidReq req) (api.post="/v1/supplysrv/internal/yahoo/placeBid")
    string YahooGetTransactionService(1: YahooGetTransactionReq req) (api.get="/v1/supplysrv/internal/yahoo/transaction")
    // YahooGetTransactionsResp YahooGetTransactionsService(1: YahooGetTransactionsReq req) (api.get="/v1/supplysrv/internal/yahoo/transactions")
    string YahooGetAuctionItemService(1: YahooGetAuctionItemReq req) (api.get="/v1/supplysrv/internal/yahoo/auctionItem")
    // YahooGetAuctionItemResp YahooGetAuctionItemAuthService(1: YahooGetAuctionItemAuthReq req) (api.get="/v1/supplysrv/internal/yahoo/auctionItemAuth")
    string YahooGetCategoryTreeService(1: YahooGetCategoryTreeReq req) (api.get="/v1/supplysrv/internal/yahoo/categoryTree")
    string YahooSearchAuctionsService(1: YahooSearchAuctionsReq req) (api.get="/v1/supplysrv/internal/yahoo/search")


    string YahooGetCategoryLeafService(1: YahooGetCategoryLeafReq req) (api.get="/v1/supplysrv/internal/yahoo/categoryLeaf")
    string YahooGetMyWonListService(1: YahooGetMyWonListReq req) (api.get="/v1/supplysrv/internal/yahoo/myWonList")
    string YahooGetSellingListService(1: YahooGetSellingListReq req) (api.get="/v1/supplysrv/internal/yahoo/sellingList")
}