package yahoo

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	bizErr "github.com/buyandship/supply-service/biz/common/err"
	"github.com/buyandship/supply-service/biz/model/bns/supply"
	"github.com/buyandship/supply-service/biz/model/yahoo"
)

type GetCategoryTreeResponse struct {
	ResultSet struct {
		Result                yahoo.Category `json:"Result"`
		TotalResultsAvailable int            `json:"totalResultsAvailable,omitempty"`
		TotalResultsReturned  int            `json:"totalResultsReturned,omitempty"`
		FirstResultPosition   int            `json:"firstResultPosition,omitempty"`
	}
}

func (c *Client) GetCategoryTree(ctx context.Context, req *supply.YahooGetCategoryTreeReq) (*GetCategoryTreeResponse, error) {
	params := url.Values{}
	params.Set("category", strconv.Itoa(int(req.Category)))
	params.Set("adf", strconv.Itoa(int(req.Adf)))
	params.Set("is_fnavi_only", strconv.Itoa(int(req.IsFnaviOnly)))

	resp, err := c.makeRequest(ctx, "GET", "/api/v1/categoryTree", params, nil, AuthTypeNone)
	if err != nil {
		if resp != nil {
			switch resp.StatusCode {
			case http.StatusUnprocessableEntity:
				return nil, bizErr.BizError{
					Status:  http.StatusBadRequest,
					ErrCode: http.StatusBadRequest,
					ErrMsg:  "validation error",
				}
			case http.StatusBadRequest:
				return nil, bizErr.BizError{
					Status:  http.StatusNotFound,
					ErrCode: http.StatusNotFound,
					ErrMsg:  "category not found",
				}
			case http.StatusInternalServerError:
				return nil, bizErr.BizError{
					Status:  http.StatusInternalServerError,
					ErrCode: 10001,
					ErrMsg:  "internal server error",
				}
			}
		}
		return nil, bizErr.BizError{
			Status:  http.StatusInternalServerError,
			ErrCode: http.StatusInternalServerError,
			ErrMsg:  "internal server error",
		}
	}
	var httpResp GetCategoryTreeResponse
	if err := c.parseResponse(resp, &httpResp); err != nil {
		return nil, err
	}

	return &httpResp, nil
}
