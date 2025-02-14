package mercari

import (
	"context"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/buyandship/supply-svr/biz/infrasturcture/db"
	"github.com/buyandship/supply-svr/biz/model/bns/supply"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func RegisterAccountService(ctx context.Context, req *supply.MercariRegisterAccountReq) (*supply.MercariRegisterAccountResp, error) {
	hlog.CtxInfof(ctx, "RegisterAccountService is called: %+v", req)

	if req.GetBuyerID() == 0 {
		hlog.CtxErrorf(ctx, "empty buyer_id")
		return nil, bizErr.InvalidParameterError
	}

	if err := db.GetHandler().UpsertAccount(ctx, &model.Account{
		Email:          req.GetEmail(),
		BuyerID:        req.GetBuyerID(),
		Prefecture:     req.GetPrefecture(),
		FamilyName:     req.GetFamilyName(),
		FirstName:      req.GetFirstName(),
		FamilyNameKana: req.GetFamilyNameKana(),
		FirstNameKana:  req.GetFirstNameKana(),
		Telephone:      req.GetTelephone(),
		ZipCode1:       req.GetZipCode1(),
		ZipCode2:       req.GetZipCode2(),
		City:           req.GetCity(),
		Address1:       req.GetAddress1(),
		Address2:       req.GetAddress2(),
	}); err != nil {
		return nil, err
	}

	acc, err := db.GetHandler().GetAccount(ctx, req.GetBuyerID())
	if err != nil {
		return nil, err
	}

	return &supply.MercariRegisterAccountResp{
		Email:          acc.Email,
		BuyerID:        acc.BuyerID,
		Prefecture:     acc.Prefecture,
		FamilyName:     acc.FamilyName,
		FirstName:      acc.FirstName,
		FirstNameKana:  acc.FirstNameKana,
		FamilyNameKana: acc.FamilyNameKana,
		Telephone:      acc.Telephone,
		ZipCode1:       acc.ZipCode1,
		ZipCode2:       acc.ZipCode2,
		City:           acc.City,
		Address1:       acc.Address1,
		Address2:       acc.Address2,
	}, nil
}
