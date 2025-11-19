package yahoo

import (
	"context"
	"net/http"
	"net/url"

	bizErr "github.com/buyandship/supply-svr/biz/common/err"
)

// AuctionItemDetail represents detailed Yahoo Auction item information
type AuctionItemDetail struct {
	AuctionID                  string                  `json:"AuctionID,omitempty" example:"x123456789"`
	CategoryID                 int                     `json:"CategoryID" example:"22216"`
	CategoryFarm               int                     `json:"CategoryFarm,omitempty" example:"2"`
	CategoryIdPath             string                  `json:"CategoryIdPath,omitempty" example:"0,2084005403,22216"`
	CategoryPath               string                  `json:"CategoryPath,omitempty" example:"オークション > 音楽 > CD > R&B、ソウル"`
	Title                      string                  `json:"Title,omitempty" example:"【新品未開封】サンプルCD アルバム 限定版"`
	SeoKeywords                string                  `json:"SeoKeywords,omitempty" example:"CD,R&B,ソウル,新品,未開封"`
	Seller                     SellerInfo              `json:"Seller,omitempty"`
	ShoppingItemCode           string                  `json:"ShoppingItemCode,omitempty" example:"shopping_item_abc123"`
	AuctionItemUrl             string                  `json:"AuctionItemUrl,omitempty" example:"https://page.auctions.yahoo.co.jp/jp/auction/x123456789"`
	Img                        Images                  `json:"Img,omitempty"`
	ImgColor                   string                  `json:"ImgColor,omitempty" example:"red"`
	Thumbnails                 Thumbnails              `json:"Thumbnails,omitempty"`
	Initprice                  float64                 `json:"Initprice,omitempty" example:"1000"`
	LastInitprice              float64                 `json:"LastInitprice,omitempty" example:"1200"`
	Price                      float64                 `json:"Price,omitempty" example:"2480"`
	TaxinStartPrice            float64                 `json:"TaxinStartPrice,omitempty" example:"1080"`
	TaxinPrice                 float64                 `json:"TaxinPrice,omitempty" example:"2678.4"`
	TaxinBidorbuy              float64                 `json:"TaxinBidorbuy,omitempty" example:"5400"`
	Bidorbuy                   float64                 `json:"Bidorbuy,omitempty" example:"5000"`
	TaxRate                    int                     `json:"TaxRate,omitempty" example:"8"`
	Quantity                   int                     `json:"Quantity,omitempty" example:"2"`
	AvailableQuantity          int                     `json:"AvailableQuantity,omitempty" example:"1"`
	WatchListNum               int                     `json:"WatchListNum,omitempty" example:"42"`
	Bids                       int                     `json:"Bids,omitempty" example:"5"`
	HighestBidders             *HighestBidders         `json:"HighestBidders,omitempty"`
	YPoint                     int                     `json:"YPoint,omitempty" example:"10"`
	ItemStatus                 ItemStatus              `json:"ItemStatus,omitempty"`
	ItemReturnable             ItemReturnable          `json:"ItemReturnable,omitempty"`
	StartTime                  string                  `json:"StartTime,omitempty" example:"2025-01-15T10:00:00+09:00"`
	EndTime                    string                  `json:"EndTime,omitempty" example:"2025-02-15T23:59:59+09:00"`
	IsBidCreditRestrictions    bool                    `json:"IsBidCreditRestrictions,omitempty" example:"true"`
	IsBidderRestrictions       bool                    `json:"IsBidderRestrictions,omitempty" example:"true"`
	IsBidderRatioRestrictions  bool                    `json:"isBidderRatioRestrictions,omitempty" example:"false"`
	IsEarlyClosing             bool                    `json:"IsEarlyClosing,omitempty" example:"false"`
	IsAutomaticExtension       bool                    `json:"IsAutomaticExtension,omitempty" example:"true"`
	IsOffer                    bool                    `json:"IsOffer,omitempty" example:"true"`
	IsCharity                  bool                    `json:"IsCharity,omitempty" example:"false"`
	Option                     Option                  `json:"Option,omitempty"`
	Description                string                  `json:"Description,omitempty" example:"<![CDATA[新品未開封のCDアルバムです。限定版となります。<br>送料無料でお届けします。]]>"`
	ItemDescriptionURL         string                  `json:"ItemDescriptionURL,omitempty" example:"https://pageX.auctions.yahoo.co.jp/jp/show/description?aID=x123456789&plainview=1"`
	DescriptionInputType       string                  `json:"Description_input_type,omitempty" example:"html"`
	Payment                    Payment                 `json:"Payment,omitempty"`
	BlindBusiness              string                  `json:"BlindBusiness,omitempty" example:"impossible"`
	SevenElevenReceive         string                  `json:"SevenElevenReceive,omitempty" example:"impossible"`
	ChargeForShipping          string                  `json:"ChargeForShipping,omitempty" example:"seller"`
	Location                   string                  `json:"Location,omitempty" example:"東京都"`
	IsWorldwide                bool                    `json:"IsWorldwide,omitempty" example:"true"`
	ShipTime                   string                  `json:"ShipTime,omitempty" example:"after"`
	ShippingInput              string                  `json:"ShippingInput,omitempty" example:"now"`
	IsYahunekoPack             bool                    `json:"IsYahunekoPack,omitempty" example:"true"`
	IsJPOfficialDelivery       bool                    `json:"IsJPOfficialDelivery,omitempty" example:"true"`
	IsPrivacyDeliveryAvailable bool                    `json:"IsPrivacyDeliveryAvailable,omitempty" example:"true"`
	ShipSchedule               int                     `json:"ShipSchedule,omitempty" example:"1"`
	ManualStartTime            string                  `json:"ManualStartTime,omitempty" example:"2025-01-15T10:00:00+09:00"`
	Shipping                   *Shipping               `json:"Shipping,omitempty"`
	BaggageInfo                BaggageInfo             `json:"BaggageInfo,omitempty"`
	IsAdult                    bool                    `json:"IsAdult,omitempty" example:"false"`
	IsCreature                 bool                    `json:"IsCreature,omitempty" example:"false"`
	IsSpecificCategory         bool                    `json:"IsSpecificCategory,omitempty" example:"false"`
	IsCharityCategory          bool                    `json:"IsCharityCategory,omitempty" example:"false"`
	CharityOption              *CharityOption          `json:"CharityOption,omitempty"`
	AnsweredQAndANum           int                     `json:"AnsweredQAndANum,omitempty" example:"3"`
	Status                     string                  `json:"Status" example:"open"`
	CpaRate                    int                     `json:"CpaRate,omitempty" example:"5"`
	BiddingViaCpa              bool                    `json:"BiddingViaCpa,omitempty" example:"true"`
	BrandLineIDPath            string                  `json:"BrandLineIDPath,omitempty" example:"brand123|line456"`
	BrandLineNamePath          string                  `json:"BrandLineNamePath,omitempty" example:"サンプルブランド|サンプルライン"`
	ItemSpec                   ItemSpec                `json:"ItemSpec,omitempty"`
	CatalogId                  string                  `json:"CatalogId,omitempty" example:"catalog_12345"`
	ProductName                string                  `json:"ProductName,omitempty" example:"サンプルCD アルバム"`
	Car                        *Car                    `json:"Car,omitempty"`
	OfferNum                   int                     `json:"OfferNum,omitempty" example:"2"`
	HasOfferAccept             bool                    `json:"HasOfferAccept,omitempty" example:"false"`
	ArticleNumber              string                  `json:"ArticleNumber,omitempty" example:"1234567890123"`
	IsDsk                      bool                    `json:"IsDsk,omitempty" example:"true"`
	CategoryInsuranceType      int                     `json:"CategoryInsuranceType,omitempty" example:"1"`
	ExternalFleaMarketInfo     *ExternalFleaMarketInfo `json:"ExternalFleaMarketInfo,omitempty"`
	ShoppingSpecs              *ShoppingSpecs          `json:"ShoppingSpecs,omitempty"`
	ItemTagList                *ItemTagList            `json:"ItemTagList,omitempty"`
	ShoppingItem               ShoppingItem            `json:"ShoppingItem,omitempty"`
	IsWatched                  bool                    `json:"IsWatched,omitempty" example:"true"`
	NotifyID                   string                  `json:"NotifyID,omitempty" example:"notify_abc123"`
	StoreSearchKeywords        string                  `json:"StoreSearchKeywords,omitempty" example:"CD,音楽,限定版"`
	SellingInfo                *SellingInfo            `json:"SellingInfo,omitempty"`
	AucUserIdContactUrl        string                  `json:"AucUserIdContactUrl,omitempty" example:"https://auctions.yahoo.co.jp/jp/show/contact?aID=x123456789"`
	WinnersInfo                *WinnersInfo            `json:"WinnersInfo,omitempty"`
	ReservesInfo               *ReservesInfo           `json:"ReservesInfo,omitempty"`
	CancelsInfo                *CancelsInfo            `json:"CancelsInfo,omitempty"`
	BidInfo                    *BidInfo                `json:"BidInfo,omitempty"`
	OfferInfo                  *OfferInfo              `json:"OfferInfo,omitempty"`
	EasyPaymentInfo            *EasyPaymentInfo        `json:"EasyPaymentInfo,omitempty"`
	StorePayment               *StorePayment           `json:"StorePayment,omitempty"`
}

// Response models
type AuctionItemResponse struct {
	ResultSet struct {
		TotalResultsAvailable int               `json:"totalResultsAvailable"`
		TotalResultsReturned  int               `json:"totalResultsReturned"`
		FirstResultPosition   int               `json:"firstResultPosition"`
		Result                AuctionItemDetail `json:"Result"`
	} `json:"ResultSet"`
}

// GetAuctionItem gets auction item information (public API)
func (c *Client) GetAuctionItem(ctx context.Context, req AuctionItemRequest) (*AuctionItemResponse, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/auctionItem", params, nil, AuthTypeNone)
	if err != nil {
		return nil, err
	}

	var auctionItemResponse AuctionItemResponse
	if err := c.parseResponse(resp, &auctionItemResponse); err != nil {
		return nil, err
	}

	return &auctionItemResponse, nil
}

func (c *Client) MockGetAuctionItem(ctx context.Context, req AuctionItemRequest) (*AuctionItemResponse, error) {
	auctionItemResponse := AuctionItemResponse{
		ResultSet: struct {
			TotalResultsAvailable int               `json:"totalResultsAvailable"`
			TotalResultsReturned  int               `json:"totalResultsReturned"`
			FirstResultPosition   int               `json:"firstResultPosition"`
			Result                AuctionItemDetail `json:"Result"`
		}{
			Result: AuctionItemDetail{
				AuctionID:      "x123456789",
				CategoryID:     22216,
				CategoryFarm:   2,
				CategoryIdPath: "0,2084005403,22216",
				CategoryPath:   "オークション > 音楽 > CD > R&B、ソウル",
				Title:          "【新品未開封】サンプルCD アルバム 限定版",
				SeoKeywords:    "CD,R&B,ソウル,新品,未開封",
				Seller: SellerInfo{
					AucUserId: "sample_seller_123",
					Rating: Rating{
						Point:                   150,
						TotalGoodRating:         145,
						TotalNormalRating:       3,
						TotalBadRating:          2,
						SellerTotalGoodRating:   120,
						SellerTotalNormalRating: 2,
						SellerTotalBadRating:    1,
						IsSuspended:             false,
						IsDeleted:               false,
					},
					AucUserIdItemListURL: "https://auctions.yahooapis.jp/AuctionWebService/V2/sellingList?appid=xxxxx&ItemListAucUserIdUrl=sample_seller_123",
					AucUserIdRatingURL:   "https://auctions.yahooapis.jp/AuctionWebService/V1/ShowRating?appid=xxxxx&RatingAucUserIdUrl=sample_seller_123",
					DisplayName:          "サンプルセラー",
					StoreName:            "サンプルストア",
					IconUrl128:           "https://s.yimg.jp/images/auct/profile/icon/128/sample_seller_123.jpg",
					IconUrl256:           "https://s.yimg.jp/images/auct/profile/icon/256/sample_seller_123.jpg",
					IconUrl512:           "https://s.yimg.jp/images/auct/profile/icon/512/sample_seller_123.jpg",
					IsStore:              true,
					ShoppingSellerId:     "store_12345",
					Performance:          map[string]interface{}{},
				},
				ShoppingItemCode: "shopping_item_abc123",
				AuctionItemUrl:   "https://page.auctions.yahoo.co.jp/jp/auction/x123456789",
				Img: Images{
					Image1: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890abc.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像1",
					},
					Image2: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890def.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像2",
					},
					Image3: &ImageInfo{
						URL:    "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-img600x450-1234567890ghi.jpg",
						Width:  600,
						Height: 450,
						Alt:    "商品画像3",
					},
				},
				ImgColor: "red",
				Thumbnails: Thumbnails{
					Thumbnail1: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890abc.jpg",
					Thumbnail2: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890def.jpg",
					Thumbnail3: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890ghi.jpg",
				},
				Initprice:         1000,
				LastInitprice:     1200,
				Price:             2480,
				TaxinStartPrice:   1080,
				TaxinPrice:        2678.4,
				TaxinBidorbuy:     5400,
				Bidorbuy:          5000,
				TaxRate:           8,
				Quantity:          2,
				AvailableQuantity: 1,
				WatchListNum:      42,
				Bids:              5,
				YPoint:            10,
				ItemStatus: ItemStatus{
					Condition: "new",
					Comment:   "新品・未開封品です",
				},
				ItemReturnable: ItemReturnable{
					Allowed: true,
					Comment: "未開封のため返品可能です",
				},
				StartTime:                 "2025-01-15T10:00:00+09:00",
				EndTime:                   "2025-02-15T23:59:59+09:00",
				IsBidCreditRestrictions:   true,
				IsBidderRestrictions:      true,
				IsBidderRatioRestrictions: false,
				IsEarlyClosing:            false,
				IsAutomaticExtension:      true,
				IsOffer:                   true,
				IsCharity:                 false,
				Option: Option{
					StoreIcon:            "https://image.auctions.yahoo.co.jp/images/store.gif",
					FeaturedIcon:         "https://image.auctions.yahoo.co.jp/images/featured.gif",
					FreeshippingIcon:     "https://image.auctions.yahoo.co.jp/images/freeshipping.gif",
					NewItemIcon:          "https://image.auctions.yahoo.co.jp/images/newitem.gif",
					EasyPaymentIcon:      "https://img.yahoo.co.jp/images/pay/icon_s16.gif",
					IsTradingNaviAuction: true,
				},
				Description:          "<![CDATA[新品未開封のCDアルバムです。限定版となります。<br>送料無料でお届けします。]]>",
				ItemDescriptionURL:   "https://pageX.auctions.yahoo.co.jp/jp/show/description?aID=x123456789&plainview=1",
				DescriptionInputType: "html",
				Payment: Payment{
					YBank: map[string]interface{}{},
					EasyPayment: &EasyPayment{
						SafeKeepingPayment: "1.00",
						IsCreditCard:       true,
						AllowInstallment:   true,
						IsPayPay:           true,
					},
					Bank: &BankPayment{
						TotalBankMethodAvailable: 3,
						Method: []BankMethod{
							{Name: "三菱UFJ銀行", BankID: "0005"},
							{Name: "みずほ銀行", BankID: "0001"},
							{Name: "ゆうちょ銀行", BankID: "9900"},
						},
					},
					CashRegistration: "可能",
					PostalTransfer:   "可能",
					PostalOrder:      "可能",
					CashOnDelivery:   "可能",
					Other: &OtherPayment{
						TotalOtherMethodAvailable: 1,
						Method:                    []string{"手渡し"},
					},
				},
				BlindBusiness:              "impossible",
				SevenElevenReceive:         "impossible",
				ChargeForShipping:          "seller",
				Location:                   "東京都",
				IsWorldwide:                true,
				ShipTime:                   "after",
				ShippingInput:              "now",
				IsYahunekoPack:             true,
				IsJPOfficialDelivery:       true,
				IsPrivacyDeliveryAvailable: true,
				ShipSchedule:               1,
				ManualStartTime:            "2025-01-15T10:00:00+09:00",
				Shipping: &Shipping{
					TotalShippingMethodAvailable: 4,
					LowestIndex:                  0,
				},
				BaggageInfo: BaggageInfo{
					Size:        "～70cm",
					SizeIndex:   1,
					Weight:      "～4kg",
					WeightIndex: 2,
				},
				IsAdult:            false,
				IsCreature:         false,
				IsSpecificCategory: false,
				IsCharityCategory:  false,
				CharityOption: &CharityOption{
					Proportion: 10,
				},
				AnsweredQAndANum:  3,
				Status:            "open",
				CpaRate:           5,
				BiddingViaCpa:     true,
				BrandLineIDPath:   "brand123|line456",
				BrandLineNamePath: "サンプルブランド|サンプルライン",
				ItemSpec: ItemSpec{
					Size:    "M",
					Segment: "メンズ",
				},
				CatalogId:   "catalog_12345",
				ProductName: "サンプルCD アルバム",
				Car: &Car{
					TotalCosts:              250000,
					TaxinTotalCosts:         270000,
					TotalPrice:              1250000,
					TaxinTotalPrice:         1350000,
					TotalBidorbuyPrice:      1500000,
					TaxinTotalBidorbuyPrice: 1620000,
					OverheadCosts:           150000,
					TaxinOverheadCosts:      162000,
					LegalCosts:              100000,
					ContactTelNumber:        "03-1234-5678",
					ContactReceptionTime:    "平日10:00-18:00",
					ContactUrl:              "https://example.com/contact",
					Regist: CarRegist{
						Model: "2020年式",
					},
					Options: CarOptions{
						Item: []string{"カーナビ", "ETC", "バックカメラ"},
					},
					TotalAmountComment: "諸費用込みの総額です",
				},
				OfferNum:              2,
				HasOfferAccept:        false,
				ArticleNumber:         "1234567890123",
				IsDsk:                 true,
				CategoryInsuranceType: 1,
				ExternalFleaMarketInfo: &ExternalFleaMarketInfo{
					IsWinner: false,
				},
				ShoppingSpecs: &ShoppingSpecs{
					TotalShoppingSpecs: 2,
					Spec: []ShoppingSpec{
						{ID: 100, ValueID: 1001},
						{ID: 200, ValueID: 2001},
					},
				},
				ItemTagList: &ItemTagList{
					TotalItemTagList: 2,
					Tag:              []string{"adidas", "Nike"},
				},
				ShoppingItem: ShoppingItem{
					PostageSetId:    12345,
					PostageId:       67890,
					LeadTimeId:      5000,
					ItemWeight:      500,
					IsOptionEnabled: true,
				},
				IsWatched:           true,
				NotifyID:            "notify_abc123",
				StoreSearchKeywords: "CD,音楽,限定版",
				SellingInfo: &SellingInfo{
					PageView:                        1250,
					WatchListNum:                    42,
					ReportedViolationNum:            0,
					AnsweredQAndANum:                3,
					UnansweredQAndANum:              1,
					OfferNum:                        2,
					UnansweredOfferNum:              0,
					AffiliateRatio:                  5,
					PageViewFromAff:                 180,
					WatchListNumFromAff:             8,
					BidsFromAff:                     2,
					IsWon:                           false,
					IsFirstSubmit:                   true,
					Duration:                        7,
					FirstAutoResubmitAvailableCount: 3,
					AutoResubmitAvailableCount:      3,
					FeaturedDpd:                     "500",
					ResubmitPriceDownRatio:          5,
					IsNoResubmit:                    false,
					BidQuantityLimit:                5,
				},
				AucUserIdContactUrl: "https://auctions.yahoo.co.jp/jp/show/contact?aID=x123456789",
				WinnersInfo: &WinnersInfo{
					WinnersNum: 1,
					Winner: []Winner{
						{
							AucUserId: "winner_user_1",
							Rating: WinnerRating{
								Point:       120,
								IsSuspended: false,
								IsDeleted:   false,
								IsNotRated:  false,
							},
							IsRemovable:        true,
							RemovableLimitTime: 1739836800,
							WonQuantity:        1,
							LastBidQuantity:    1,
							WonPrice:           2480,
							TaxinWonPrice:      2678.4,
							LastBidTime:        1707955199,
							BuyTime:            "2025-02-15T23:59:59+09:00",
							IsFnaviBundledDeal: false,
							ShoppingInfo: &WinnerShoppingInfo{
								OrderId: "order_12345",
							},
							DisplayName: "落札者1",
							IconUrl128:  "https://s.yimg.jp/images/auct/profile/icon/128/winner_user_1.jpg",
							IconUrl256:  "https://s.yimg.jp/images/auct/profile/icon/256/winner_user_1.jpg",
							IconUrl512:  "https://s.yimg.jp/images/auct/profile/icon/512/winner_user_1.jpg",
							IsStore:     false,
						},
					},
				},
				ReservesInfo: &ReservesInfo{
					ReservesNum: 1,
					Reserve: []Reserve{
						{
							AucUserId: "reserve_user_1",
							Rating: ReserveRating{
								Point:       85,
								IsSuspended: false,
								IsDeleted:   false,
							},
							LastBidQuantity:   1,
							LastBidPrice:      2450,
							TaxinLastBidPrice: 2646,
							LastBidTime:       1707955150,
							DisplayName:       "次点者1",
							IconUrl128:        "https://s.yimg.jp/images/auct/profile/icon/128/reserve_user_1.jpg",
							IconUrl256:        "https://s.yimg.jp/images/auct/profile/icon/256/reserve_user_1.jpg",
							IconUrl512:        "https://s.yimg.jp/images/auct/profile/icon/512/reserve_user_1.jpg",
							IsStore:           false,
						},
					},
				},
				CancelsInfo: &CancelsInfo{
					CancelsNum: 0,
					Cancel:     []Cancel{},
				},
				BidInfo: &BidInfo{
					IsHighestBidder: true,
					IsWinner:        false,
					IsDeletedWinner: false,
					IsNextWinner:    false,
					LastBid: LastBid{
						Price:              2480,
						TaxinPrice:         2678.4,
						Quantity:           1,
						Partial:            false,
						IsFnaviBundledDeal: false,
					},
					NextBid: NextBid{
						Price:         2580,
						LimitQuantity: 1,
						UnitPrice:     100,
					},
				},
				OfferInfo: &OfferInfo{
					OfferCondition:      1,
					SellerOfferredPrice: 4500,
					BidderOfferredPrice: 4000,
					RemainingOfferNum:   2,
				},
				EasyPaymentInfo: &EasyPaymentInfo{
					EasyPayment: EasyPaymentDetail{
						AucUserId:  "winner_user_1",
						Status:     "completed",
						LimitTime:  1710547199,
						UpdateTime: 1709337599,
					},
				},
				StorePayment: &StorePayment{
					TotalStorePaymentMethodAvailable: 2,
					Method:                           []string{"クレジットカード", "代金引換"},
					UpdateTime:                       1704614400,
				},
			},
			TotalResultsAvailable: 1,
			TotalResultsReturned:  1,
			FirstResultPosition:   1,
		},
	}
	return &auctionItemResponse, nil
}

// GetAuctionItemAuth gets authenticated auction item information
func (c *Client) GetAuctionItemAuth(ctx context.Context, req AuctionItemRequest, yahooAccountID string) (*AuctionItemResponse, error) {
	params := url.Values{}
	params.Set("auctionID", req.AuctionID)
	params.Set("yahoo_account_id", yahooAccountID)
	if req.AppID != "" {
		params.Set("appid", req.AppID)
	}

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/auctionItemAuth", params, nil, AuthTypeHMAC)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusBadRequest:
				{
					return nil, bizErr.BizError{
						Status:  http.StatusNotFound,
						ErrCode: http.StatusNotFound,
						ErrMsg:  "auction item not found",
					}
				}
			}
		}
		return nil, bizErr.BizError{
			Status:  http.StatusInternalServerError,
			ErrCode: http.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}

	var auctionItemAuthResponse AuctionItemResponse
	if err := c.parseResponse(resp, &auctionItemAuthResponse); err != nil {
		return nil, err
	}

	return &auctionItemAuthResponse, nil
}
