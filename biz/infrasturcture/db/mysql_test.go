package db

import (
	"context"
	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/model/mercari"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	t.Run("GetAccounts", func(t *testing.T) {
		config.LoadTestConfig()
		h := GetHandler()
		mapAcc, err := h.GetAccounts()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(mapAcc)
	})
}

func TestUpsertAccount(t *testing.T) {
	t.Run("UpsertAccount", func(t *testing.T) {
		config.LoadTestConfig()
		h := GetHandler()
		acc := &mercari.Account{
			BuyerID: "1",
		}
		if err := h.UpsertAccount(context.Background(), acc); err != nil {
			t.Fatal(err)
		}
		t.Log(acc)
	})
}
