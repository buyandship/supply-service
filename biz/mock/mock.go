package mock

import (
	"context"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/bns-golib/config"
	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/infrastructure/mercari"
	"github.com/buyandship/supply-service/biz/infrastructure/yahoo"
	"github.com/google/uuid"
)

func MockMercariGetItemError(itemId string) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}
	if itemId == "m92155064693" {
		return bizErr.BadRequestError
	}

	if itemId == "m11873603604" {
		return bizErr.UnauthorisedError
	}

	if itemId == "m45423510521" {
		return bizErr.InvalidInputError
	}

	if itemId == "m88281581466" {
		return bizErr.NotFoundError
	}

	if itemId == "m81372405611" {
		return bizErr.ConflictError
	}

	if itemId == "m62402785917" {
		return bizErr.MercariInternalError
	}

	if itemId == "m85986863833" {
		return bizErr.InternalError
	}

	if itemId == "m31731301399" {
		return bizErr.UndefinedError
	}

	if itemId == "m64928494499" {
		return bizErr.TooManyRequestError
	}

	if itemId == "m65772869996" {
		return bizErr.NotFoundError
	}

	return nil
}

func MockMercariPostOrderError(itemId string) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if itemId == "m62418751854" {
		return bizErr.BadRequestError
	}

	if itemId == "m64057832200" {
		return bizErr.PaymentRequiredError
	}

	if itemId == "m10317110362" {
		return bizErr.ConflictError
	}

	if itemId == "m72838781411" {
		return bizErr.MercariInternalError
	}

	if itemId == "m79102263412" {
		return bizErr.InternalError
	}

	if itemId == "m41528787258" {
		return bizErr.ForbiddenError
	}

	if itemId == "m65378237792" {
		return bizErr.NotFoundError
	}

	if itemId == "m55618538005" {
		return bizErr.UndefinedError
	}

	if itemId == "m44557269639" {
		return bizErr.TooManyRequestError
	}

	if itemId == "m83954959553" {
		return bizErr.BadRequestError
	}

	if itemId == "m63823469042" {
		return bizErr.PaymentRequiredError
	}

	if itemId == "m24876695495" {
		return bizErr.MercariInternalError
	}

	if itemId == "m71491191679" {
		return bizErr.InternalError
	}

	/*
		if itemId == "m51533021958" {
			return bizErr.PaymentRequiredError
		}

		if itemId == "m81061748245" {
			return bizErr.InternalError
		}
	*/

	return nil
}

func MockMercariPostMessageError(itemId string) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if itemId == "m10398666665" {
		return bizErr.BadRequestError
	}

	if itemId == "m75543122082" {
		return bizErr.UnauthorisedError
	}

	if itemId == "m75718811593" {
		return bizErr.NotFoundError
	}

	if itemId == "m70687287564" {
		return bizErr.TooManyRequestError
	}

	if itemId == "m61551956239" {
		return bizErr.NotFoundError
	}

	if itemId == "m66542179466" {
		return bizErr.ConflictError
	}

	if itemId == "m31731301399" {
		return bizErr.MercariInternalError
	}

	if itemId == "m66542179466" {
		return bizErr.InternalError
	}

	return nil
}

func MockMercariSellerError(sellerId string) error {

	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if sellerId == "127569133" {
		return bizErr.BadRequestError
	}

	if sellerId == "582382837" {
		return bizErr.UnauthorisedError
	}

	if sellerId == "459538216" {
		return bizErr.ConflictError
	}

	if sellerId == "949242743" {
		return bizErr.MercariInternalError
	}

	if sellerId == "619741287" {
		return bizErr.InternalError
	}

	if sellerId == "107875930" {
		return bizErr.UndefinedError
	}

	if sellerId == "454117226" {
		return bizErr.InvalidInputError
	}

	if sellerId == "233277177" {
		return bizErr.NotFoundError
	}

	if sellerId == "110763907" {
		return bizErr.ConflictError
	}

	if sellerId == "568974622" {
		return bizErr.MercariInternalError
	}

	if sellerId == "215320064" {
		return bizErr.InternalError
	}

	if sellerId == "850025121" {
		return bizErr.TooManyRequestError
	}

	return nil
}

func MockMercariItemResponse(resp *mercari.GetItemByIDResponse) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if resp != nil {
		if resp.Id == "m31548572003" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m65596663947" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m85991927697" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m40565066644" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m16920742682" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m74571668950" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m31592028768" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m13506094471" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m30254025612" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 100
		}

		if resp.Id == "m80942035809" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m40193840651" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m15614884768" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m89066517919" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m77065327937" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m95350323871" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m89987642745" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m29930173669" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m38657767746" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m72592393057" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m88973643686" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m28450369044" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m28315619029" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m70978824732" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m79555391109" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}

		if resp.Id == "m97778471314" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 50
		}

		if resp.Id == "m94826714492" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 10
		}

		if resp.Id == "m98984482817" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 50
		}

	}
	return nil
}

func MockMercariCategory(cid string) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if cid == "208" {
		return bizErr.BadRequestError
	}
	if cid == "4364" {
		return bizErr.UnauthorisedError
	}
	if cid == "7708" {
		return bizErr.ConflictError
	}
	if cid == "4779" {
		return bizErr.MercariInternalError
	}
	if cid == "412" {
		return bizErr.InternalError
	}
	if cid == "179" {
		return bizErr.UndefinedError
	}
	if cid == "3838" {
		return bizErr.InvalidInputError
	}
	if cid == "4359" {
		return bizErr.NotFoundError
	}
	return nil
}

func MockMercariGetTransactionByItemId(resp *mercari.GetTransactionByItemIDResponse) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	// Case 1: Transactions with tracking numbers
	trackingNumberCases := map[string]string{
		"m92760180388": "520198765432",
		"m52779338879": "520287654321",
		"m57479709388": "520376543210",
		"m88774196481": "520465432109",
		"m89528493541": "520554321098",
		"m61022850419": "520643210987",
		"m86545992008": "520732109876",
		"m82944619012": "520821098765",
		"m80349753918": "520912345678",
		"m82088723367": "521023456789",
		"m78778189968": "521134567890",
		"m29456318542": "521245678901",
		"m94614181196": "521356789012",
		"m83031001810": "521467890123",
		"m96164664967": "521578901234",
		"m94455048583": "521689012345",
		"m65651422692": "521790123456",
		"m24472580811": "521801234567",
		"m80382710973": "521912345678",
		"m19205542174": "522023456789",
		"m15297147252": "522134567890",
		"m56496480575": "522245678901",
		"m87339221708": "522356789012",
		"m27781850905": "522467890123",
		"m67177021875": "522578901234",
		"m27928417925": "522689012345",
		"m67570147112": "522790123456",
		"m42339612386": "522801234567",
		"m24477724560": "522912345678",
		"m84912732956": "random_tracking_number",
		"m23552581043": "522912345670",
		"m21398498724": "",
	}

	// Case 2: Transactions without tracking numbers
	noTrackingCases := map[string]bool{
		"m16155925172": true,
		"m20792871193": true,
		"m45868093966": true,
		"m57275155732": true,
		"m80883954147": true,
		"m58081973812": true,
		"m59919792096": true,
		"m42717636168": true,
		"m49675936947": true,
		"m46075964356": true,
		"m70687287564": true,
		"m61551956239": true,
		"m31337561540": true,
		"m26922603870": true,
		"m41740765280": true,
		"m39411058099": true,
		"m49714846192": true,
		"m96503874959": true,
		"m91470723883": true,
	}

	waitShippingCases := map[string]string{
		"m87457435408": "random_tracking_number",
	}

	// Case 3: Error cases
	errorCases := map[string]error{
		"m54081840575": bizErr.UndefinedError,
		"m79752279771": bizErr.UnauthorisedError,
		"m92588011186": bizErr.ConflictError,
		"m81844429902": bizErr.RateLimitError,
		"m31143390731": bizErr.MercariInternalError,
		"m21457538128": bizErr.InternalError,
	}

	if trackingNumber, ok := trackingNumberCases[resp.ItemId]; ok {
		if trackingNumber == "random_tracking_number" {
			resp.ShippingInfo.TrackingNumber = uuid.NewString()
		} else {
			resp.ShippingInfo.TrackingNumber = trackingNumber
		}
		resp.Status = "wait_review"
		return nil
	}

	if _, ok := noTrackingCases[resp.ItemId]; ok {
		resp.ShippingInfo.TrackingNumber = ""
		resp.Status = "wait_review"
		return nil
	}

	if trackingNumber, ok := waitShippingCases[resp.ItemId]; ok {
		if trackingNumber == "random_tracking_number" {
			resp.ShippingInfo.TrackingNumber = uuid.NewString()
		} else {
			resp.ShippingInfo.TrackingNumber = trackingNumber
		}
		resp.Status = "wait_shipping"
		return nil
	}

	if err, ok := errorCases[resp.ItemId]; ok {
		return err
	}

	return nil
}

func MockYahooGetAuctionItemError(auctionId string) error {

	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if auctionId == "w1211863278" {
		return bizErr.BadRequestError
	}

	if auctionId == "n1211899153" {
		return bizErr.InternalError
	}

	if auctionId == "c1211905254" {
		return bizErr.TimeoutError
	}

	if auctionId == "w1211863278" {
		return bizErr.BadRequestError
	}

	return nil
}

func MockYahooPlaceBidError(auctionId string) error {

	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if auctionId == "l1211860869" {
		return bizErr.BadRequestError
	}

	if auctionId == "l1211858121" {
		return bizErr.InternalError
	}

	if auctionId == "s1211903305" {
		return bizErr.BadRequestError
	}

	if auctionId == "l1211896005" {
		return bizErr.NotFoundError
	}

	if auctionId == "n1211897053" {
		return bizErr.InternalError
	}
	return nil
}

func MockYahooGetAuctionItemDetail(resp *yahoo.AuctionItemResponse) error {
	if config.GlobalAppConfig.Env != "dev" {
		return nil
	}

	if resp.ResultSet.Result.AuctionID == "b1212185797" {
		resp.ResultSet.Result.Quantity = 5
		resp.ResultSet.Result.AvailableQuantity = 4
	}

	return nil
}

func UpdateNextBidPrice() {
	p := &yahoo.AuctionItemResponse{}
	if err := cache.GetRedisClient().Get(context.Background(), "bravo_test_item", p); err != nil {
		return
	}

	p.ResultSet.Result.Price = float64(p.ResultSet.Result.BidInfo.NextBid.Price)
	if p.ResultSet.Result.Price > 2000 {
		// Win bid
		p.ResultSet.Result.WinnersInfo.Winner = []yahoo.Winner{{
			AucUserId: "AnzTKsBM5HUpBc3CCQc3dHpETkds1",
			Rating: yahoo.WinnerRating{
				Point: 150,
			},
			WonPrice: int(p.ResultSet.Result.Price),
		}}
	}
	p.ResultSet.Result.TaxinPrice = float64(p.ResultSet.Result.BidInfo.NextBid.Price) * 1.1
	p.ResultSet.Result.BidInfo.NextBid.Price += 100
	if err := cache.GetRedisClient().Set(context.Background(), "bravo_test_item", p, 0); err != nil {
		return
	}

}

func TestAuction() *yahoo.AuctionItemResponse {
	p := &yahoo.AuctionItemResponse{}
	if err := cache.GetRedisClient().Get(context.Background(), "bravo_test_item", p); err == nil {
		return p
	}

	initProduct := &yahoo.AuctionItemResponse{
		ResultSet: struct {
			TotalResultsAvailable int                     `json:"@totalResultsAvailable"`
			TotalResultsReturned  int                     `json:"@totalResultsReturned"`
			FirstResultPosition   int                     `json:"@firstResultPosition"`
			Result                yahoo.AuctionItemDetail `json:"Result"`
		}{
			Result: yahoo.AuctionItemDetail{
				AuctionID:      "bravo_test_item",
				CategoryID:     22216,
				CategoryFarm:   2,
				CategoryIdPath: "0,2084005403,22216",
				CategoryPath:   "オークション > 音楽 > CD > R&B、ソウル",
				Title:          "【新品未開封】サンプルCD アルバム 限定版",
				SeoKeywords:    "CD,R&B,ソウル,新品,未開封",
				Seller: yahoo.SellerInfo{
					AucUserId: "AnzTKsBM5HUpBc3CCQc3dHpETkds1",
					Rating: yahoo.Rating{
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
				ImgColor:         "red",
				Thumbnails: yahoo.Thumbnails{
					Thumbnail1: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890abc.jpg",
					Thumbnail2: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890def.jpg",
					Thumbnail3: "https://auctions.c.yimg.jp/images.auctions.yahoo.co.jp/image/dr000/auc0101/users/1/2/3/4/sample_user-thumb-1234567890ghi.jpg",
				},
				Initprice:         1000,
				LastInitprice:     1200,
				Price:             1000,
				TaxinStartPrice:   1100,
				TaxinPrice:        1100,
				TaxinBidorbuy:     5500,
				Bidorbuy:          5000,
				TaxRate:           10,
				Quantity:          2,
				AvailableQuantity: 1,
				WatchListNum:      42,
				Bids:              5,
				YPoint:            10,
				ItemStatus: yahoo.ItemStatus{
					Condition: "new",
					Comment:   "新品・未開封品です",
				},
				ItemReturnable: yahoo.ItemReturnable{
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
				Option: yahoo.Option{
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
				Payment: yahoo.Payment{
					YBank: map[string]interface{}{},
					EasyPayment: &yahoo.EasyPayment{
						SafeKeepingPayment: "1.00",
						IsCreditCard:       true,
						AllowInstallment:   true,
						IsPayPay:           true,
					},
					Bank: &yahoo.BankPayment{
						TotalBankMethodAvailable: 3,
						Method: []yahoo.BankMethod{
							{Name: "三菱UFJ銀行", BankID: "0005"},
							{Name: "みずほ銀行", BankID: "0001"},
							{Name: "ゆうちょ銀行", BankID: "9900"},
						},
					},
					CashRegistration: "可能",
					PostalTransfer:   "可能",
					PostalOrder:      "可能",
					CashOnDelivery:   "可能",
					Other: &yahoo.OtherPayment{
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
				Shipping: &yahoo.Shipping{
					TotalShippingMethodAvailable: 4,
					LowestIndex:                  0,
				},
				BaggageInfo: yahoo.BaggageInfo{
					Size:        "～70cm",
					SizeIndex:   1,
					Weight:      "～4kg",
					WeightIndex: 2,
				},
				IsAdult:            false,
				IsCreature:         false,
				IsSpecificCategory: false,
				IsCharityCategory:  false,
				CharityOption: &yahoo.CharityOption{
					Proportion: 10,
				},
				AnsweredQAndANum:  3,
				Status:            "open",
				CpaRate:           5,
				BiddingViaCpa:     true,
				BrandLineIDPath:   "brand123|line456",
				BrandLineNamePath: "サンプルブランド|サンプルライン",
				ItemSpec: yahoo.ItemSpec{
					Size:    "M",
					Segment: "メンズ",
				},
				CatalogId:   "catalog_12345",
				ProductName: "サンプルCD アルバム",
				Car: &yahoo.Car{
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
					Regist: yahoo.CarRegist{
						Model: "2020年式",
					},
					Options: yahoo.CarOptions{
						Item: []string{"カーナビ", "ETC", "バックカメラ"},
					},
					TotalAmountComment: "諸費用込みの総額です",
				},
				OfferNum:              2,
				HasOfferAccept:        false,
				ArticleNumber:         "1234567890123",
				IsDsk:                 true,
				CategoryInsuranceType: 1,
				ExternalFleaMarketInfo: &yahoo.ExternalFleaMarketInfo{
					IsWinner: false,
				},
				ShoppingSpecs: &yahoo.ShoppingSpecs{
					TotalShoppingSpecs: 2,
					Spec: []yahoo.ShoppingSpec{
						{ID: 100, ValueID: 1001},
						{ID: 200, ValueID: 2001},
					},
				},
				ItemTagList: &yahoo.ItemTagList{
					TotalItemTagList: 2,
					Tag:              []string{"adidas", "Nike"},
				},
				ShoppingItem: &yahoo.ShoppingItem{
					PostageSetId:    12345,
					PostageId:       67890,
					LeadTimeId:      5000,
					ItemWeight:      500,
					IsOptionEnabled: true,
				},
				IsWatched:           true,
				NotifyID:            "notify_abc123",
				StoreSearchKeywords: "CD,音楽,限定版",
				SellingInfo: &yahoo.SellingInfo{
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
				WinnersInfo: &yahoo.WinnersInfo{
					WinnersNum: 1,
					Winner: []yahoo.Winner{
						{
							AucUserId: "winner_user_1",
							Rating: yahoo.WinnerRating{
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
							ShoppingInfo: &yahoo.WinnerShoppingInfo{
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
				ReservesInfo: &yahoo.ReservesInfo{
					ReservesNum: 1,
					Reserve: []yahoo.Reserve{
						{
							AucUserId: "reserve_user_1",
							Rating: yahoo.ReserveRating{
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
				CancelsInfo: &yahoo.CancelsInfo{
					CancelsNum: 0,
					Cancel:     []yahoo.Cancel{},
				},
				BidInfo: &yahoo.BidInfo{
					IsHighestBidder: true,
					IsWinner:        false,
					IsDeletedWinner: false,
					IsNextWinner:    false,
					LastBid: yahoo.LastBid{
						Price:              1000,
						TaxinPrice:         1100,
						Quantity:           1,
						Partial:            false,
						IsFnaviBundledDeal: false,
					},
					NextBid: yahoo.NextBid{
						Price:         1100,
						LimitQuantity: 1,
						UnitPrice:     100,
					},
				},
				OfferInfo: &yahoo.OfferInfo{
					OfferCondition:      1,
					SellerOfferredPrice: 4500,
					BidderOfferredPrice: 4000,
					RemainingOfferNum:   2,
				},
				EasyPaymentInfo: &yahoo.EasyPaymentInfo{
					EasyPayment: yahoo.EasyPaymentDetail{
						AucUserId:  "winner_user_1",
						Status:     "completed",
						LimitTime:  1710547199,
						UpdateTime: 1709337599,
					},
				},
				StorePayment: &yahoo.StorePayment{
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

	if err := cache.GetRedisClient().Set(context.Background(), "bravo_test_item", initProduct, 0); err != nil {
		return nil
	}

	return initProduct
}
