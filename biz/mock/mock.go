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

func MockMercariItemResponse(resp *mercari.GetItemByIDResponse) error {
	if config.GlobalServerConfig.Env != "development" {
		return nil
	}

	if resp != nil {
		if resp.Id == "m30254025612" {
			resp.Price = 300
			resp.Discounts.TotalReturnAbsolute = 100
		}

		if resp.Id == "m24744986573" {
			resp.Price = 300
			resp.ItemDiscount.ReturnAbsolute = 10
		}
	}
	return nil
}
