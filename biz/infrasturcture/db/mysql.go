package db

import (
	"context"
	"fmt"
	"github.com/buyandship/supply-svr/biz/common/config"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	once    sync.Once
	Handler *H
)

type H struct {
	cli *gorm.DB
}

func GetHandler() *H {
	once.Do(func() {
		c, err := gorm.Open(mysql.New(
			mysql.Config{
				DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
					config.GlobalServerConfig.Mysql.Username,
					config.GlobalServerConfig.Mysql.Password,
					config.GlobalServerConfig.Mysql.Address,
					config.GlobalServerConfig.Mysql.DBName,
				),
			},
		), &gorm.Config{})
		if err != nil {
			hlog.Fatal("mysql init err:", err)
		}
		Handler = &H{cli: c}
	})
	return Handler
}

func (h *H) HealthCheck() error {
	sqlDB, err := h.cli.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	return nil
}

func (h *H) UpsertAccount(ctx context.Context, account *model.Account) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Where(&model.Account{BuyerID: account.BuyerID}).
		FirstOrCreate(&account).Error; err != nil {
		hlog.Errorf("upsert account err, %s", err.Error())
		return err
	}
	return nil
}

func (h *H) InsertMessage(ctx context.Context, message *model.Message) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Create(&message).Error; err != nil {
		hlog.Errorf("insert message err, %s", err.Error())
		return err
	}
	return nil
}

func (h *H) InsertTransaction(ctx context.Context, transaction *model.Transaction) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Create(&transaction).Error; err != nil {
		hlog.Errorf("insert transaction err, %s", err.Error())
		return err
	}
	return nil
}

func (h *H) UpdateTransaction(ctx context.Context, cond *model.Transaction) error {
	if err := h.cli.
		WithContext(ctx).
		Debug().
		Where("trx_id = ?", cond.TrxID).
		Updates(cond).Error; err != nil {
		hlog.Errorf("update transaction err, %s", err.Error())
		return err
	}
	return nil
}

func (h *H) GetAccounts() (map[string]model.Account, error) {
	accounts := make(map[string]model.Account)
	var dbAccounts []*model.Account
	if err := h.cli.Debug().Find(&dbAccounts).Error; err != nil {
		hlog.Errorf("get account err %s", err.Error())
		return nil, err
	}
	for _, dbAccount := range dbAccounts {
		accounts[dbAccount.BuyerID] = *dbAccount
	}
	return accounts, nil
}
