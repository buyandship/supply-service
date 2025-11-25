package mercari

import (
	"context"
	"fmt"

	"github.com/buyandship/bns-golib/cache"
	"github.com/buyandship/supply-service/biz/common/config"
	"github.com/buyandship/supply-service/biz/infrastructure/db"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	model "github.com/buyandship/supply-service/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func RegisterAccountService(ctx context.Context, req *supply.MercariRegisterAccountReq) (*supply.MercariRegisterAccountResp, error) {
	hlog.CtxInfof(ctx, "RegisterAccountService is called: %+v", req)

	acc := &model.Account{
		Email:          req.GetEmail(),
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
		Priority:       int(req.GetPriority()),
	}

	if err := db.GetHandler().UpsertAccount(ctx, acc); err != nil {
		return nil, err
	}
	// delete cache
	if err := cache.GetRedisClient().Del(ctx, fmt.Sprintf(config.MercariAccountPrefix, acc.ID)); err != nil {
		return nil, err
	}

	acc, err := db.GetHandler().GetAccount(ctx, int32(acc.ID))
	if err != nil {
		return nil, err
	}

	return &supply.MercariRegisterAccountResp{
		Account: acc.Thrift(),
	}, nil
}
