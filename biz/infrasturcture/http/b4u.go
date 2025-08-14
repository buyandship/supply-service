package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
)

type Notifier struct {
	Endpoint string
	Token    string
}

func GetNotifier() *Notifier {
	endpoint := "https://b4u-admin-test.buynship.com"
	if config.GlobalServerConfig.Env == "production" {
		endpoint = "https://b4u-admin.buynship.com"
	}
	token := config.GlobalServerConfig.B4uToken
	return &Notifier{
		Endpoint: endpoint,
		Token:    token,
	}
}

func (n *Notifier) Notify(ctx context.Context, body any) error {
	cli := http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("%s/api/admin/callback/mercari/switch-account/", n.Endpoint)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", n.Token))
	req.Header.Set("Content-Type", "application/json")

	resp, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to notify b4u: %s", resp.Status)
	}

	return nil
}
