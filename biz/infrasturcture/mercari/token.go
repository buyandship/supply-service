package mercari

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
	"net/url"
)

type GetTokenRequest struct {
	BuyerID     string `json:"buyer_id"`
	RedirectUrl string `json:"redirect_url"`
}

type GetTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (m *Mercari) GetToken(ctx context.Context, req *GetTokenRequest) (*GetTokenResponse, error) {
	parsedUrl, err := url.Parse(req.RedirectUrl)
	if err != nil {
		hlog.Errorf("url parse error: %s", err.Error())
		return nil, bizErr.InvalidParameterError
	}
	code := parsedUrl.Query().Get("code")
	scope := parsedUrl.Query().Get("scope")

	secret, err := m.GenerateSecret(req.BuyerID)
	if err != nil {
		return nil, err
	}

	headers := map[string][]string{
		"Content-Type":  []string{"application/x-www-form-urlencoded"},
		"Authorization": []string{fmt.Sprintf("Basic %s", secret)},
	}

	body := map[string]string{
		"grant_type":   "authorization_code",
		"scope":        scope,
		"redirect_uri": req.RedirectUrl,
		"code":         code,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		hlog.Errorf("json marshal error, body: %s, err: %v", body, err)
		return nil, bizErr.InternalError
	}
	data := bytes.NewBuffer(reqBody)
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/jp/v1/token", m.AuthServiceDomain), data)
	if err != nil {
		hlog.Errorf("http request error, err: %v", err)
		return nil, bizErr.InternalError
	}
	httpReq.Header = headers
	client := &http.Client{}
	httpRes, err := client.Do(httpReq)
	defer func() {
		if err := httpRes.Body.Close(); err != nil {
			hlog.CtxErrorf(ctx, "http close error: %s", err)
		}
	}()
	if err != nil {
		hlog.Errorf("http error, err: %v", err)
		return nil, bizErr.InternalError
	}

	if httpRes.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpRes.Body)
		hlog.Errorf("get mercari token error: %s", respBody)
		return nil, bizErr.BizError{
			Status:  httpRes.StatusCode,
			ErrCode: httpRes.StatusCode,
			ErrMsg:  "get mercari token failed",
		}
	}

	resp := &GetTokenResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
		hlog.Errorf("decode http response error, err: %v", err)
		return nil, bizErr.InternalError
	}
	return resp, nil
}
