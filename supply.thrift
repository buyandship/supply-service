// idl/hello.thrift
namespace go bns.supply

struct MercariGetItemReq {
    1: string buyer_id (api.json="buyer_id");
    2: string item_id (api.json="item_id"); // 添加 api 注解为方便进行参数绑定
}

struct MercariGetItemResp {
}

struct MercariGetCategoriesReq {
    1: string buyer_id (api.json="buyer_id");
}

struct MercariGetCategoriesResp {

}

struct MercariGetSellerReq {
    1: string buyer_id (api.json="buyer_id");
    2: string seller_id (api.json="seller_id");
}

struct MercariGetSellerResp {

}

struct MercariPostOrderReq {
    1: string buyer_id (api.json="buyer_id");
    2: string item_id (api.json="item_id")
    3: i32 ref_price (api.json="ref_price")
    4: string ref_currency (api.json="ref_currency")
    5: string checksum (api.json="checksum")
    6: string ref_id (api.json="ref_id")
}

struct MercariPostOrderResp {
}

struct MercariPostMessageReq {
    1: string buyer_id (api.json="buyer_id");
    2: string trx_id (api.json="trx_id");
    3: string msg (api.json="msg");
}

struct MercariPostMessageResp {
}


struct MercariRegisterAccountReq {
    1: string buyer_id (api.json="buyer_id");
    2: string email (api.json="email");
    3: string redirectUrl (api.json="redirect_url");
    4: string family_name (api.json="family_name");
    5: string first_name (api.json="first_name");
    6: string family_name_kana (api.json="family_name_kana");
    7: string first_name_kana (api.json="first_name_kana");
    8: string telephone (api.json="telephone");
    9: string zip_code1 (api.json="zip_code1");
    10: string zip_code2 (api.json="zip_code2");
    11: string prefecture (api.json="prefecture");
    12: string city (api.json="city");
    13: string address1 (api.json="address1");
    14: string address2 (api.json="address2");
    15: string client_id (api.json="client_id");
    16: string client_secret (api.json="client_secret")
}

struct MercariRegisterAccountResp {

}


service SupplyService {
    MercariRegisterAccountResp MercariRegisterAccountService(1: MercariRegisterAccountReq req) (api.post="/mercari/register");
    MercariGetItemResp MercariGetItemService(1: MercariGetItemReq req) (api.get="/mercari/item");
    MercariGetCategoriesResp MercariGetCategoriesService(1: MercariGetCategoriesReq req) (api.get="/mercari/categories")
    MercariGetSellerResp MercariGetSellerService(1: MercariGetSellerReq req) (api.get="/mercari/seller")
    MercariPostOrderResp MercariPostOrderService(1: MercariPostOrderReq req) (api.post="/mercari/order")
    MercariPostMessageResp MercariPostMessageService(1: MercariPostMessageReq req) (api.post="/mercari/message")
}