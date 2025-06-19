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

    MercariRegisterAccountResp MercariRegisterAccountService(1: MercariRegisterAccountReq req) (api.post="/v1/supplysrv/internal/mercari/register");
    MercariPostOrderResp MercariPostOrderService(1: MercariPostOrderReq req) (api.post="/v1/supplysrv/internal/mercari/order")
    MercariPostMessageResp MercariPostMessageService(1: MercariPostMessageReq req) (api.post="/v1/supplysrv/internal/mercari/message")
    MercariGetTokenResp MercariGetTokenService() (api.get="/v1/supplysrv/internal/mercari/token")
    MercariGetAccountResp MercariGetAccountService() (api.get="/v1/supplysrv/internal/mercari/account/list")

}