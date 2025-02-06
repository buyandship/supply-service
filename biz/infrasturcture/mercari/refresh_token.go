package mercari

import (
	"bytes"
	"encoding/json"
	"fmt"
	bizErr "github.com/buyandship/supply-svr/biz/common/err"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"io"
	"net/http"
)

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func (m *Mercari) RefreshToken(buyerID string) error {
	// call refresh token
	acc, ok := m.Accounts[buyerID]
	if !ok {
		hlog.Errorf("buyer not exists, buyer_id: %s", buyerID)
		return bizErr.InvalidBuyerError
	}

	secret, err := m.GenerateSecret(buyerID)
	if err != nil {
		return err
	}

	headers := map[string][]string{
		"Content-Type":  []string{"application/x-www-form-urlencoded"},
		"Authorization": []string{fmt.Sprintf("Basic %s", secret)},
	}

	body := map[string]string{
		"grant_type":    "refresh_token",
		"scope":         "openapi:buy%20offline_access",
		"refresh_token": acc.RefreshToken,
	}

	reqBody, err := json.Marshal(body)
	if err != nil {
		hlog.Errorf("json marshal error, body: %s, err: %v", body, err)
		return bizErr.InternalError
	}

	data := bytes.NewBuffer(reqBody)
	httpReq, err := http.NewRequest("POST", fmt.Sprintf("%s/jp/v1/token", m.AuthServiceDomain), data)
	if err != nil {
		hlog.Errorf("http request error, err: %v", err)
		return bizErr.InternalError
	}
	httpReq.Header = headers
	client := &http.Client{}
	httpRes, err := client.Do(httpReq)
	defer func() {
		if err := httpRes.Body.Close(); err != nil {
			hlog.Errorf("http close error: %s", err)
		}
	}()
	if err != nil {
		hlog.Errorf("http error, err: %v", err)
		return bizErr.InternalError
	}
	if httpRes.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(httpRes.Body)
		hlog.Errorf("purchase mercari item error: %s", respBody)
		return bizErr.UnauthorisedError
	}

	resp := &RefreshTokenResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(resp); err != nil {
		hlog.Errorf("decode http response error, err: %v", err)
		return bizErr.InternalError
	}

	go func() {

	}()

	acc.AccessToken = resp.AccessToken
	return nil
}
