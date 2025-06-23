package err

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type BizError struct {
	Status  int    `json:"status"`
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func (e BizError) Error() string {
	return e.ErrMsg
}

func ConvertErr(err error) BizError {
	Err := BizError{}
	if errors.As(err, &Err) {
		return Err
	}
	s := BizError{
		Status:  consts.StatusInternalServerError,
		ErrCode: 1000, // unknown error
		ErrMsg:  err.Error(),
	}
	return s
}

var (
	UnauthorisedError     = BizError{Status: consts.StatusUnauthorized, ErrCode: consts.StatusUnauthorized, ErrMsg: "unauthorized"}
	InvalidParameterError = BizError{Status: consts.StatusBadRequest, ErrCode: consts.StatusBadRequest, ErrMsg: "invalid parameter"}
	ConflictError         = BizError{Status: consts.StatusConflict, ErrCode: consts.StatusConflict, ErrMsg: "conflict"}
	BadRequestError       = BizError{Status: consts.StatusBadRequest, ErrCode: consts.StatusBadRequest, ErrMsg: "bad request"}
	ForbiddenError        = BizError{Status: consts.StatusForbidden, ErrCode: consts.StatusForbidden, ErrMsg: "forbidden"}
	PaymentRequiredError  = BizError{Status: consts.StatusPaymentRequired, ErrCode: consts.StatusPaymentRequired, ErrMsg: "payment required"}
	NotFoundError         = BizError{Status: consts.StatusNotFound, ErrCode: consts.StatusNotFound, ErrMsg: "not found"}
	InvalidInputError     = BizError{Status: consts.StatusForbidden, ErrCode: consts.StatusForbidden, ErrMsg: "invalid input"}
	TooManyRequestError   = BizError{Status: consts.StatusTooManyRequests, ErrCode: consts.StatusTooManyRequests, ErrMsg: "too many request"}
	AccountBannedError    = BizError{Status: consts.StatusForbidden, ErrCode: consts.StatusForbidden, ErrMsg: "account banned"}

	MercariInternalError      = BizError{Status: consts.StatusInternalServerError, ErrCode: consts.StatusInternalServerError, ErrMsg: "mercari internal error"}
	InternalError             = BizError{Status: consts.StatusInternalServerError, ErrCode: 1000, ErrMsg: "internal error"}
	InvalidBuyerError         = BizError{Status: consts.StatusInternalServerError, ErrCode: 1001, ErrMsg: "invalid buyer id"}
	TooLowReferencePriceError = BizError{Status: consts.StatusInternalServerError, ErrCode: 1002, ErrMsg: "the reference price is too low"}
	ItemNotOnSaleError        = BizError{Status: consts.StatusInternalServerError, ErrCode: 1003, ErrMsg: "item not on sale"}
	InvalidCheckSumError      = BizError{Status: consts.StatusInternalServerError, ErrCode: 1004, ErrMsg: "invalid checksum"}
	RefIdDuplicatedError      = BizError{Status: consts.StatusInternalServerError, ErrCode: 1005, ErrMsg: "duplicated ref_id"}
	RateLimitError            = BizError{Status: consts.StatusInternalServerError, ErrCode: 1006, ErrMsg: "rate limit"}
	UnloginError              = BizError{Status: consts.StatusInternalServerError, ErrCode: 1007, ErrMsg: "mercari unlogin"}
	UndefinedError            = BizError{Status: consts.StatusInternalServerError, ErrCode: 5678, ErrMsg: "undefined error"}
	ACLBanError               = BizError{Status: consts.StatusForbidden, ErrCode: 17, ErrMsg: "acl ban"}
)
