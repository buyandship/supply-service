package mock

import (
	"github.com/buyandship/bns-golib/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
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
