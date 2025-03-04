package mock

import (
	"github.com/buyandship/supply-svr/biz/common/config"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/mercari"
)

func MockMercariGetItemError(itemId string) error {
	if config.GlobalServerConfig.Env != "development" {
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

	return nil
}

func MockMercariPostOrderError(itemId string) error {
	if config.GlobalServerConfig.Env != "development" {
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

	return nil
}

func MockMercariPostMessageError(itemId string) error {
	if config.GlobalServerConfig.Env != "development" {
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

	return nil
}

func MockMercariItemResponse(resp *mercari.GetItemByIDResponse) error {
	if config.GlobalServerConfig.Env != "development" {
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
	}
	return nil
}

func MockMercariCategory(cid string) error {
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
