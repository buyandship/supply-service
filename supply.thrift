// idl/hello.thrift
namespace go bns.supply

struct MercariGetItemReq {
    1: i32 buyer_id (api.json="buyer_id");
    2: string item_id (api.json="item_id");
}

struct MercariGetSellerReq {
    1: string seller_id (api.json="seller_id");
}

struct MercariPostOrderReq {
    1: i32 buyer_id (api.json="buyer_id");
    2: string item_id (api.json="item_id")
    3: i64 ref_price (api.json="ref_price")
    4: string ref_currency (api.json="ref_currency")
    5: string checksum (api.json="checksum")
    6: string ref_id (api.json="ref_id")
}

struct MercariPostMessageReq {
    1: i64 trx_id (api.json="trx_id");
    2: string msg (api.json="msg");
}

struct MercariPostMessageResp {
    1: i64 trx_id
    2: string body
    3: string id
}

struct MercariRegisterAccountReq {
    1: i32 buyer_id (api.json="buyer_id");
    2: string email (api.json="email");
    3: string family_name (api.json="family_name");
    4: string first_name (api.json="first_name");
    5: string family_name_kana (api.json="family_name_kana");
    6: string first_name_kana (api.json="first_name_kana");
    7: string telephone (api.json="telephone");
    8: string zip_code1 (api.json="zip_code1");
    9: string zip_code2 (api.json="zip_code2");
    10: string prefecture (api.json="prefecture");
    11: string city (api.json="city");
    12: string address1 (api.json="address1");
    13: string address2 (api.json="address2");
}

struct MercariRegisterAccountResp {
    1: i32 buyer_id
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

service SupplyService {
    MercariRegisterAccountResp MercariRegisterAccountService(1: MercariRegisterAccountReq req) (api.post="/v1/supplysrv/mercari/register");
    string MercariGetItemService(1: MercariGetItemReq req) (api.get="/v1/supplysrv/mercari/item");
    string MercariGetCategoriesService() (api.get="/v1/supplysrv/mercari/categories")
    string MercariGetSellerService(1: MercariGetSellerReq req) (api.get="/v1/supplysrv/mercari/seller")
    string MercariPostOrderService(1: MercariPostOrderReq req) (api.post="/v1/supplysrv/mercari/order")
    MercariPostMessageResp MercariPostMessageService(1: MercariPostMessageReq req) (api.post="/v1/supplysrv/mercari/message")
    string MercariLoginCallBackService(1: MercariLoginCallBackReq req) (api.get="/xb/login_callback")
    MercariGetTokenResp MercariGetTokenService() (api.get="/v1/supplysrv/mercari/token")
}